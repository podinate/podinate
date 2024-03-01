# Podinate Concepts
This page provides an overview of terms and components that make up the Podinate Cloud ecosystem.

## Account
A podinate account is the root level collection of cloud resources. All resources must be part of an account. This keeps ownership and access separate. 

## Cluster
A cluster is a set of servers which run your Podinate workloads. Multiple Podinate clusters can be managed by a central 

## Controller
A Podinate controller is a central control point to multiple clusters on which your workloads can run. You can either use the Podinate hosted controller, which enables you to get up and running quickly, or install your own hosted controller. 

## IAM 
Identitiy and Access Management. Podinate uses policy-based Identity and Access Manangement to decide which requestors can perform which actions on Podinate. 

## Pod
A pod is a running container or group of idential containers running on Podinate. 

## Project 
Resources in Podinate are grouped into projects to make it clear what their purpose is. For example you might have an [Account](#Account) called 'department-of-awesome' and then have projects for 'Beer Ordering System' and 'Backflip Counter'. If the department of awesome is done with backflips, they can delete the entire 'Backflip Counter' project.

<!-- ## Requestor
In [IAM](#IAM) a requestor is an Actor that is requesting the given Action.  -->

## User
A user is a person that can perform actions in the Podinate API. User accounts in Podinate are provided by a Git provider. 

## Volume
A volume is persistent storage attached to a Pod. 