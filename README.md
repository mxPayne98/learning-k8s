# Learning k8s

Learn k8s and its core concepts step by step.

We'll be using the following minimal k8s cluster setup during the learning journey:

- A separate namespace for our microservices.
- Node.js and Go applications, each with its own Redis instance, isolated within the namespace.
- An NGINX Ingress controller to manage external access, routing specific paths to each application.

[Cluster setup](cluster/SETUP.md)

[App setup](apps/SETUP.md)

## Docs

### Core concepts

1. [Pods](docs/01.core/01.pods/pods.md)
2. [ReplicaSets](docs/01.core/02.02.replicasets/replicasets.md)
3. [Deployments](docs/01.core/03.deployments/deployments.md)
4. [Namespaces](docs/01.core/04.namespaces/namespaces.md)
