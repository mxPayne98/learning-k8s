To set up a complete Kubernetes cluster with all necessary components, including the Nginx ingress controller, namespaces, and ensuring only the applications (not Redis instances) are exposed to the internet, we'll need to follow these steps:

### 1. **Creating Namespaces**

Namespaces in Kubernetes provide a mechanism for isolating groups of resources within a single cluster. Let's create a namespace for our microservices environment.

#### Namespace YAML File

**microservices-namespace.yaml**:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: microservices
```

Apply this file using `kubectl apply -f microservices-namespace.yaml` to create the `microservices` namespace.

### 2. **Deploying Applications and Redis in the Namespace**

When deploying the Node.js, Go, and Redis applications, specify the namespace in their YAML files. This is done by adding the `namespace` field under `metadata`.

For example, in the Node.js deployment file, you would add:

```yaml
metadata:
  name: node-app
  namespace: microservices
```

Repeat this for all deployment and service YAML files for both the applications and their Redis instances.

### 3. **Setting Up NGINX Ingress Controller**

The NGINX Ingress Controller is a Kubernetes controller that manages external access to HTTP services in a Kubernetes cluster using NGINX as a reverse proxy and load balancer.

#### Install NGINX Ingress Controller

You can install the NGINX Ingress Controller using Helm (a Kubernetes package manager) or directly with YAML files. Here's the general command if using Helm:

```bash
helm install my-nginx-ingress ingress-nginx/ingress-nginx --namespace microservices
```

#### Configure Ingress Resource

Once the Ingress Controller is installed, you need to define Ingress Resources to route traffic to your Node.js and Go applications.

**ingress-resource.yaml**:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: microservices-ingress
  namespace: microservices
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - http:
      paths:
      - path: /node
        pathType: Prefix
        backend:
          service:
            name: node-service
            port:
              number: 80
      - path: /go
        pathType: Prefix
        backend:
          service:
            name: go-service
            port:
              number: 80
```

This Ingress Resource routes `/node` to the Node.js service and `/go` to the Go service.

### 4. **Applying the Configurations**

- Apply all updated YAML files for the Node.js, Go, and Redis applications within the `microservices` namespace.
- Apply the Ingress resource configuration.

### 5. **Accessing the Applications**

After applying the Ingress configuration, your applications should be accessible from outside the cluster. The exact URL will depend on your cluster setup and the domain you have configured.
