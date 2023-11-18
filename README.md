# Learning k8s

Learn k8s and its core concepts step by step.

We'll be using the following minimal k8s cluster setup during the learning journey:

- A separate namespace for our microservices.
- Node.js and Go applications, each with its own Redis instance, isolated within the namespace.
- An NGINX Ingress controller to manage external access, routing specific paths to each application.

[Cluster setup](cluster/SETUP.md)

## Commands

All k8s components use following set of commands for various operations:  

`kubectl create`

`kubectl get`

`kubectl describe`

`kubect delete`

`kubectl edit`

To get info about all components at once use `kubectl get all`


- ### Nodes
    Get all nodes in a cluster
    ```
    kubectl get nodes
    ```
    Use `-o wide` to get additional info
    ```
    kubectl get nodes -o wide
    ```

- ### PODs

    Create/Run a POD Directly from command using an image

    ```
    kubectl run nginx --image=nginx
    ```

    From YAML file definition

    ```
    kubectl create -f ./pods/nginx-pod.yml
    ```
    _**NOTE**: The `kubectl create` and `kubectl apply` commands work the same when creating something new_

    Get all running PODS:

    ```
    kubectl get pods
    ```
    Use `-o wide` to get additional info
    ```
    kubectl get pods -o wide
    ```
    Get all details about a POD 
    ```
    kubectl describe pod nginx
    ```
    To delete a pod
    ```
    kubectl delete replicaset nginx-pod
    ```

- ### ReplicaSet
    Create a replica set
    ```
    kubectl create -f ./replicasets/nginx-replicaset.yml
    ```
    Get all running ReplicaSets:
    ```
    kubectl get replicaset
    ```
    To delete a replicaset
    ```
    kubectl delete replicaset nginx-replicaset
    ```
    _**NOTE**: While deleting replicaset it will also delete the underlying PODs_

    To scale a replica set, either modify the YAML file with desired number of replicas and run
    ```
    kubectl replace replicaset -f ./replicasets/nginx-replicaset.yml
    ```
    _**NOTE**: YAML file of any running configuration can also be edited by running the `edit` command_
    ```
    kubectl edit replicaset nginx-replicaset
    ```
    Alternatively, use the `scale` command
    ```
    kubectl scale --replicas=6 -f ./replicasets/nginx-replicaset.yml
    ```
    Get all details about a replica set 
    ```
    kubectl describe replicaset nginx-replicaset
    ```

- ### Deployments
    All commands and YAML config same as replicaset (use 'deployment' instead of 'replicaset')

    #### **Updates and Rollbacks for a deployment**

    To update a deployment use
    ```
    kubectl apply -f ./deployments/nginx-deployment
    ```
    Alternatively for changing image version use:
    ```
    kubectl set image deployments/nginx-deployment nginx=nginx:1.9.1
    ```

    To check rollout status
    ```
    kubectl rollout status deployment/nginx-deployment
    ```
    To check rollout history
    ```
    kubectl rollout history deployment/nginx-deployment
    ```
    _**NOTE**:  History shows CHANGE-CAUSE for an update. While updating (creating/editing etc.) a deployment use flag `--record` to record the cause for the deployment update._


    To rollback an update
    ```
    kubectl rollout undo deployment/nginx-deployment
    ```
    _**NOTE**: Updates can be done using either Recreate or Rolling Update Strtegies. The Recreate strategy brings down the replica set all at once and deploys new replicaset where as the rolling update strategy creates a new replicaset and brings down and creates PODs in old and new replicasets one by one respectively._