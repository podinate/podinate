#! /bin/bash

##==============================================================================
## Official Podinate installer. 
# This script is meant to be run on a fresh Ubuntu 24.04 LTS install.
# It will install all the base components needed to run the Podinate controller
##==============================================================================

## Steps:
# 1. Check if Kubernetes is installed or available 
# 2. Ask for user email
# 3. Install K3s if not already installed
# 4. Install certbot
# 5. Create Podinate namespace
# 6. Install Postgres 
# 7. Run goerd to create / migrate the database
# 8. Install Podinate controller
# 9. Check connection to the Podinate controller
# 10. Install the Podinate CLI
# 11. Check connection to the Podinate CLI

# To be confirmed:
# Create the initial user and account, setting the user email. 
# Display connection details to the user

## Install prerequisites



echo "1. Updating System..."
sudo apt-get update
sudo apt-get upgrade -y

echo "2. Installing Prerequisites..."
sudo apt-get install -y curl wget git nano dialog
sudo snap install helm 


echo "3. Installing Kubernetes..."
if kubectl get nodes ; then
    if dialog --stdout --title "Existing Cluster Detected" --clear --yesno "An existing Kubernetes cluster connection was detected. Do you want to install Podinate on the existing cluster?\n\nWe only recommend installing Podinate on a dedicated Kubernetes cluster!" 20 80 ; then
        clear
        echo "Installing Podinate on existing cluster."
    else
        clear
        echo "Installation cancelled."
        exit 1
    fi
else
    echo "Kubernetes is not installed. Installing K3s."
    curl -sfL https://get.k3s.io | sh -
    #sleep 10
fi

if details=$(dialog --stdout \
--title "About You" \
--clear --form "Please enter your email address. This will be used for any Let's Encrypt reminders and to create your account on your Podinate cluster." \
 20 80 \
 0 "Email:" 1 1 "$EMAIL" \
  1 10 50 0 \
 ); then
    echo "Email entered."
else
    echo "No email entered. Exiting."
    exit 1
fi

# Clear to prevent blue screen from prompt continuing to show
clear

EMAIL=$(echo $details | awk '{print $1}')
echo "Email: $EMAIL"

echo "4. Setting up certbot..."
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.0/cert-manager.yaml

# Wait for cert-manager to be ready
echo "5. Waiting for Kubernetes install..."
until kubectl -n cert-manager wait pod --for condition=Ready -l app.kubernetes.io/component=webhook --timeout 180s
do 
    echo "5. Waiting for Kubernetes cluster to be ready..."
    sleep 5
done

cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
    name: letsencrypt-prod
spec:
    acme:
        server: https://acme-v02.api.letsencrypt.org/directory
        email: $EMAIL
        privateKeySecretRef:
            name: letsencrypt-account-prod
        solvers:
            - http01:
                ingress:
                    class: traefik
EOF

# Creating the Podinate namespace...
echo "Creating Podinate namespace..."
if ! (kubectl create namespace podinate --dry-run=client -o yaml | kubectl apply -f -) ; then
    #echo "Podinate namespace created."
    echo "Error creating Podinate namespace."
    exit 1
fi

echo "6. Installing Postgres..."
kubectl apply -f https://raw.githubusercontent.com/podinate/podinate/main/kubernetes/masterdb-postgres.yaml

# Make the passwords random
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: masterdb-secret
  namespace: podinate
type: Opaque
stringData:
  superUserPassword: $(openssl rand -base64 32 | tr -cd '[:alpha:]')
  replicationUserPassword: $(openssl rand -base64 32 | tr -cd '[:alpha:]')
EOF

echo "6. Waiting for database to be ready..."
kubectl -n podinate wait pod --for=condition=Ready -l app=postgres --timeout 180s

# Run the Postgres migrations
echo "6. Setting up database..."
cat << EOF | kubectl apply -f - 
apiVersion: v1
kind: ConfigMap
metadata:
  name: migration-script
  namespace: podinate
data:
  migration.sh: |
    #!/bin/sh
    wget https://raw.githubusercontent.com/podinate/podinate/main/database/atlas.hcl
    /atlas schema apply --url "postgres://\$POSTGRES_USER:\$POSTGRES_PASSWORD@\$POSTGRES_HOST:5432/\$POSTGRES_DB?sslmode=disable" --to file://atlas.hcl --auto-approve

---
apiVersion: batch/v1
kind: Job
metadata:
  name: atlas
  namespace: podinate
  labels:
    app: atlas
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: atlas
          image: arigaio/atlas:latest-alpine
          #command: [ "/migrations/migration.sh" ]
          command: [
            "sh",
            "/migrations/migration.sh"
          ]
          env:
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: masterdb-secret
                  key: superUserPassword
            - name: POSTGRES_DB
              value: podinate
            - name: POSTGRES_HOST
              value: postgres
            - name: POSTGRES_PORT
              value: "5432"
          volumeMounts:
            - name: migration-script
              mountPath: /migrations

      volumes:
        - name: migration-script
          configMap:
            name: migration-script
            items:
              - key: migration.sh
                path: migration.sh
            defaultMode: 0777
EOF

echo "6. Waiting for migrator setup..."
kubectl -n podinate wait pod --for=condition=Ready -l job-name=atlas --timeout 120s
kubectl -n podinate logs -f -l job-name=atlas

until kubectl -n podinate wait --for condition=complete job/atlas --timeout 10s
do
    echo "6. Migrations not complete. Let's try again..."
    sleep 1
    kubectl -n podinate logs -f -l job-name=atlas
done


# Install the Podinate controller
echo "7. Installing Podinate controller..."
kubectl -n podinate delete deployment podinate-controller
kubectl -n podinate apply -f https://raw.githubusercontent.com/podinate/podinate/main/kubernetes/controller.yaml

echo "7. Waiting for controller to be ready..."
kubectl -n podinate wait pod --for=condition=Ready -l app=podinate-controller --timeout 180s

echo "8. Initializing Podinate (nearly done!)..."
# Gets the IP address of the first "enp*" interface
IP=$(ip -4 -o addr show scope global | grep enp | awk '{gsub(/\/.*/,"",$4); print $4}')

# Runs Podinate init to create the initial user and copies the profile out
kubectl -n podinate exec -it $(kubectl -n podinate get pod -l app=podinate-controller -o jsonpath='{.items[0].metadata.name}') -- podinate init --email $EMAIL --ip $IP
kubectl -n podinate cp $(kubectl -n podinate get pod -l app=podinate-controller -o jsonpath='{.items[0].metadata.name}'):/profile.yaml config.yaml

echo "Yippee! Podinate controller installed."