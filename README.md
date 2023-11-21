# Learning k8s

Learn k8s and its core concepts step by step.

We'll be using the following minimal k8s cluster setup during the learning journey:

- A separate namespace for our microservices.
- Node.js and Go applications, each with its own Redis instance, isolated within the namespace.
- An NGINX Ingress controller to manage external access, routing specific paths to each application.

[Cluster setup](cluster/SETUP.md)

[App setup](apps/SETUP.md)

## Table of Contents

| Category             | Topic            | Subtopic             | Link                                                                                  |
| -------------------- | ---------------- | -------------------- | ------------------------------------------------------------------------------------- |
| Core Concepts        | Pods             | -                    | [Pods](docs/01.core/01.pods/pods.md)                                                  |
|                      | ReplicaSets      | -                    | [ReplicaSets](docs/01.core/02.replicasets/replicasets.md)                             |
|                      | Deployments      | -                    | [Deployments](docs/01.core/03.deployments/deployments.md)                             |
|                      | Namespaces       | -                    | [Namespaces](docs/01.core/04.namespaces/namespaces.md)                                |
|                      | Jobs             | -                    | [Jobs](docs/01.core/05.jobs/jobs.md)                                                  |
| Configurations       | Containers       | Commands and Args    | [Commands and Args](docs/02.configurations/01.containers/commands-and-args.md)        |
|                      |                  | Environment          | [Environment](docs/02.configurations/01.containers/environment.md)                    |
|                      |                  | Security Context     | [Security Context](docs/02.configurations/01.containers/security-context.md)          |
|                      |                  | Service Account      | [Service Account](docs/02.configurations/01.containers/service-account.md)            |
|                      | Resources        | -                    | [Resource Requirements](docs/02.configurations/02.resources/resource-requirements.md) |
|                      | Scheduling       | Taints & Tolerations | [Taints & Tolerations](docs/02.configurations/03.scheduling/taints-tolerations.md)    |
|                      |                  | Node Selectors       | [Node Selectors](docs/02.configurations/03.scheduling/node-selectors.md)              |
|                      |                  | Node Affinity        | [Node Affinity](docs/02.configurations/03.scheduling/node-affinity.md)                |
| Multi-Container Pods | Intro            | Multi-Container Pods | [Multi-Container Pods](docs/03.multi-container-pods/multi-container-pods.md)          |
|                      | -                | Ambassador Example   | [Ambassador Example](docs/03.multi-container-pods/ambassador-example.md)              |
| Observability        |                  | Readiness Probe      | [Readiness Probe](docs/04.observability/readiness-probe.md)                           |
|                      | -                | Liveness Probe       | [Liveness Probe](docs/04.observability/liveness-probe.md)                             |
| Networking           | Services         | -                    | [Services](docs/05.networking/01.services/services.md)                                |
|                      | Ingress          | -                    | [Ingress](docs/05.networking/02.ingress/ingress.md)                                   |
|                      |                  | Controller Setup     | [Ingress Setup](docs/05.networking/02.ingress/ingress-setup.md)                       |
|                      | Network Policies | -                    | [Network Policies](docs/05.networking/03.network-policies/network-policies.md)        |
| Storage              |                  | Volumes              | [Volumes](docs/06.storage/volumes.md)                                                 |
|                      |                  | Persistent Volumes   | [Persistent Volumes](docs/06.storage/persistent-volumes.md)                           |
|                      |                  | PVC                  | [PVC](docs/06.storage/pvc.md)                                                         |
|                      |                  | Storage Class        | [Storage Class](docs/06.storage/storage-class.md)                                     |
|                      |                  | Stateful Sets        | [Stateful Sets](docs/06.storage/stateful-sets.md)                                     |
|                      | -                | MySQL Master-Slave   | [MySQL Master-Slave](docs/06.storage/mysql-master-slave.md)                           |
