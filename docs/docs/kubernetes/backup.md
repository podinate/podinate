# Backing up a Kubernetes Cluster

Configuring backup for your Kubernetes cluster is an important step when setting up a cluster. Even if you aren't storing important data (like on a dev or testing cluster), backups will help to easily roll your Kubernetes configuration back to a previous configuration if an administrator makes a mistake. 

## Velero
[Velero](https://velero.io/) is the leading backup solution for Kubernetes. It supports backing up the entire state of a Kubernetes cluster, including the state of the objects in the cluster, and the contents of volumes, to any object storage provider. 

### Basic Velero Setup
First, the Velero command line interface needs to be installed on your development machine. To get a basic an effective backup for the Kubernetes cluster, there are a couple of configuration options that need to be set. 
#### Install Velero
1. Follow the steps in [Basic Install](https://velero.io/docs/v1.14/basic-install/) to install Velero on the machine used to control the Kubernetes cluster (your laptop).
1. Create a cloud storage bucket and credentials that can only access that bucket. Put these credentials into a file:
    ```toml title="credentials-velero"
    [default]
    aws_access_key_id = minio
    aws_secret_access_key = minio123
    ```
1. Install Velero
    ```bash
    velero install --features=EnableCSI --use-node-agent \
    --provider aws \
    --use-volume-snapshots=true \
    --bucket <your-bucket> \
    --plugins velero/velero-plugin-for-aws:v1.10.0 \
    --secret-file ./credentials-velero \
    --kubeconfig ~/.kube/config \
    --backup-location-config region=nl-ams,s3Url=https://s3.nl-ams.scw.cloud
    ```
    If you're using a provider other than AWS, add these options to the backup location config `--backup-location-config region=nl-ams,s3Url=https://s3.nl-ams.scw.cloud` 

## See Also
- [Velero](https://velero.io/)
- [Velero Documentation](https://velero.io/docs/v1.14)