## Step-by-step guide to setup Ingress Controller from scratch

This guide will walk you through creating a namespace, deployment, service, ConfigMap, Secret, and roles necessary for setting up a basic Ingress controller, using NGINX as an example.

### Step 1: Create a Namespace

It's good practice to isolate the Ingress controller in its own namespace.

```yaml
# nginx-ingress-namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: nginx-ingress
```

Apply it with:
```bash
kubectl apply -f nginx-ingress-namespace.yaml
```

### Step 2: Create a ConfigMap

ConfigMaps store configuration data that can be consumed by pods.

```yaml
# nginx-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-config
  namespace: nginx-ingress
data:
  # Configuration values for NGINX
  keepalive_timeout: "65"
  max_worker_connections: "1024"
```

Apply it with:
```bash
kubectl apply -f nginx-configmap.yaml
```

### Step 3: Create a Secret

If you need to use SSL/TLS certificates, store them in a Secret.

```yaml
# nginx-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: nginx-secret
  namespace: nginx-ingress
type: Opaque
data:
  tls.crt: <base64-encoded-cert>
  tls.key: <base64-encoded-key>
```

Apply it with:
```bash
kubectl apply -f nginx-secret.yaml
```

### Step 4: Create a Service Account and Roles

Create a Service Account and the necessary RBAC roles for NGINX.

```yaml
# nginx-rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx-ingress-serviceaccount
  namespace: nginx-ingress
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: nginx-ingress-role
  namespace: nginx-ingress
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - secrets
  - services
  - endpoints
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: nginx-ingress-role-nisa-binding
  namespace: nginx-ingress
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: nginx-ingress
```

Apply it with:
```bash
kubectl apply -f nginx-rbac.yaml
```

### Step 5: Create a Deployment for the Ingress Controller

Deploy the NGINX Ingress controller.

```yaml
# nginx-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
  namespace: nginx-ingress
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-ingress
  template:
    metadata:
      labels:
        app: nginx-ingress
    spec:
      serviceAccountName: nginx-ingress-serviceaccount
      containers:
        - name: nginx-ingress-controller
          image: nginx/nginx-ingress:latest
          args:
            - /nginx-ingress-controller
            - --configmap=$(POD_NAMESPACE)/nginx-config
          env:
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          ports:
            - name: http
              containerPort: 80
            - name: https
              containerPort: 443
```

Apply it with:
```bash
kubectl apply -f nginx-deployment.yaml
```

### Step 6: Create a Service for the Ingress Controller

Expose the Ingress controller using a Service.

```yaml
# nginx-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress
  namespace: nginx-ingress
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: http
    - port: 443
      targetPort: https
  selector:
    app: nginx-ingress
```

Apply it with:
```bash
kubectl apply -f nginx-service.yaml
```

---

## Using Helm

The most common and straightforward way to set up an Ingress controller in a Kubernetes environment is by using Helm charts or applying YAML files from trusted community repositories. These methods ensure a robust and tested configuration, often with sensible defaults and the flexibility to customize as needed. 

For illustration, I'll describe setting up the NGINX Ingress controller using Helm, which is a widely adopted package manager for Kubernetes. Helm simplifies the deployment and management of complex Kubernetes applications.

### Step 1: Install Helm

First, you need to have Helm installed in your system. You can download and install Helm from its [official website](https://helm.sh/docs/intro/install/).

### Step 2: Add the NGINX Ingress Helm Repository

Add the repository that contains the NGINX Ingress Helm chart:

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
```

### Step 3: Install the NGINX Ingress Controller Using Helm

Deploy the NGINX Ingress controller in your cluster:

```bash
helm install nginx-ingress ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace
```

This command installs the NGINX Ingress controller in the `ingress-nginx` namespace, which will be created if it doesn't exist.

### Step 4: Verify the Installation

Check if the Ingress controller pods are running:

```bash
kubectl get pods -n ingress-nginx
```

### Customization and Configuration

You can customize the installation by overriding default values in the Helm chart. Create a YAML file with your configurations and use it during the installation:

```bash
helm install nginx-ingress ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace -f my-values.yaml
```

In `my-values.yaml`, you can specify various configurations like resource limits, annotations, node selectors, etc.

### Advantages of Using Helm

- **Simplicity**: Helm charts abstract the complexity of the application's structure and dependencies.
- **Manageability**: Easy to upgrade, rollback, and customize.
- **Community Support**: Helm charts are often maintained by the community or the tool's maintainers, ensuring best practices and updates.

---

## Using Community YAML Files

As an alternative to Helm, you can use YAML files provided by the community or the maintainers of the Ingress controller. For instance, the NGINX Ingress controller’s GitHub repository typically includes deployment YAMLs.

### Example: Setting Up NGINX Ingress Controller Using Official YAML Files

The NGINX Ingress controller, maintained by the Kubernetes community, provides YAML files that you can use to deploy the controller. Here’s how you can do it:

#### Step 1: Download the Official YAML Files

First, you need to download the YAML files from the official NGINX Ingress controller repository. These files are usually available in the repository's GitHub page under the `deploy` directory.

For the NGINX Ingress controller, you can find them here: [NGINX Ingress Controller Deployment YAMLs](https://github.com/kubernetes/ingress-nginx/tree/main/deploy/static/provider)

#### Step 2: Apply the YAML Files

You can directly apply these YAML files using `kubectl`. For a standard cloud environment, the process might look like this:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml
```

This command will set up the NGINX Ingress controller in your cluster, including all necessary roles, service accounts, and configurations.

#### Step 3: Verify the Installation

Check if the Ingress controller pods are running:

```bash
kubectl get pods -n ingress-nginx
```

Make sure that the pods are in a `Running` state.

### Customization

One of the key advantages of using YAML files is the ability to customize the deployment to your specific needs. You can download the YAML file, edit it to change configurations, resource limits, annotations, etc., and then apply it with `kubectl`.

For example, you might want to edit the service type, update resource requests and limits, or add specific annotations for your cloud provider.

#### Customizing Resource Requirements

Here’s an example snippet from the deployment YAML where you could customize the resource requests and limits:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ingress-nginx-controller
  namespace: ingress-nginx
  ...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - name: controller
        image: k8s.gcr.io/ingress-nginx/controller:...
        ...
        resources:
          requests:
            cpu: 100m
            memory: 90Mi
          limits:
            cpu: 200m
            memory: 250Mi
```

### Considerations

- **Compatibility**: Ensure that the version of the YAML files is compatible with your Kubernetes cluster version.
- **Security**: Review the configurations, especially if you're deploying in a production environment, to ensure they comply with your security policies.
- **Updates**: Keep an eye on the repository for updates or changes to the YAML files, especially for security patches or new features.
