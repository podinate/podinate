# Grafana Monitoring

This documentation covers the Kubernetes Prometheus and Grafana based monitoring setup. This setup is what much of the Kubernetes community has standaridised on for an open-source monitoring solution for Kubernetes. 

## Installation
The easiest way to install is using the Helm charts at https://prometheus-community.github.io/helm-charts.
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update 
```

You can then install the helm chart with the default settings (good for getting started) with the below command:
```bash
helm install -n monitoring --create-namespace prometheus-stack prometheus-community/kube-prometheus-stack
```

## Access Grafana
You can get access to the Grafana dashboard by forwarding a port on localhost into the cluster:
```bash
kubectl port-forward svc/prometheus-stack-grafana 3000:80 -n monitoring
```
The initial password is `admin` and password is `prom-operator`.

Go into Dashboards in Grafana and you will see that a bunch of useful dashboards have been created for you. 

## Add Logs 
