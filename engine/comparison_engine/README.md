# Comparison Engine

While building the Pod logic, I realised a lot of logic was being duplicated in checking and comparing the resources. So I decided to build this package which takes advantage of the Kubernetes RestHelper concept to apply this logic to any resource we could want. 

ALL LOGIC IN THIS PACKAGE MUST BE STATELESS, AND NOT MUTATE THE KUBERNETES STATE. 

## Responsibilities
This is the logic I found myself writing over and over for different resource types.

1. Check if the resource exists, if it doesn't, we create it.
1. The resource exists, run a dry run update on it. This will either tell us the resource is invalid, or it will complete a bunch of fields with their default values for us and make comparison easy. We then apply this dry-ran resource. 
1. 