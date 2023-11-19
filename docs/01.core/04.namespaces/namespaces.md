## Kubernetes Namespaces

Namespaces in Kubernetes are a fundamental concept used to organize and isolate groups of resources within a single cluster. They're particularly important in environments where multiple teams or projects share a Kubernetes cluster.

### How Namespaces Work

1. **Isolation**: Namespaces provide a logical partitioning of cluster resources. Each namespace is a distinct and separate entity within the cluster, and resources in one namespace are hidden from other namespaces.
2. **Resource Management**: They help in managing resource allocation, quotas, and access controls. For example, you can allocate specific resources to different namespaces.
3. **Organization**: Namespaces are useful for organizing resources for different environments, projects, or teams.

### Common Use Cases

- **Multi-tenancy**: Different teams or projects can safely operate in separate namespaces within the same cluster.
- **Environment Segregation**: Separate namespaces for development, testing, staging, and production environments.
- **Resource Quotas**: Enforcing limits on resource usage per namespace.

### Key Namespace Commands

1. **Creating a Namespace**:
   ```bash
   kubectl create namespace <namespace-name>
   ```

2. **Listing Namespaces**:
   ```bash
   kubectl get namespaces
   ```

3. **Deleting a Namespace**:
   ```bash
   kubectl delete namespace <namespace-name>
   ```

4. **Working Within a Specific Namespace**:
   To run commands within a specific namespace, add `--namespace=<namespace-name>` to your `kubectl` commands.

### Example YAML Configuration for a Namespace

**my-namespace.yaml**:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
```

#### Explanation of Each Field

- `apiVersion`: Specifies the version of the Kubernetes API you're using to create the object.
- `kind`: Indicates that the object type is a Namespace.
- `metadata`: 
  - `name`: The name of the Namespace.

### Gotchas and Best Practices

- **Default Namespace**: If you don't specify a namespace, Kubernetes uses the `default` namespace.
- **System Namespaces**: Kubernetes system namespaces like `kube-system` and `kube-public` should not be used for your applications.
- **Network Policies**: Remember that by default, pods in different namespaces can communicate with each other. Use Network Policies for isolation if needed.
- **Resource Visibility**: Resources like Nodes and Persistent Volumes are not namespaced. They're available cluster-wide.


## Default Namespaces in Kubernetes

Kubernetes comes with a few default namespaces that serve specific purposes:

1. **`default`**:
   - **Purpose**: This is the default namespace for objects with no other namespace. It’s commonly used for user operations and resources.
   - **Usage**: If you don't specify a namespace in your resource definitions or kubectl commands, Kubernetes assumes the `default` namespace.

2. **`kube-system`**:
   - **Purpose**: It contains resources created by the Kubernetes system, mostly system processes run by Kubernetes like the Kube DNS server, and other system components.
   - **Usage**: It’s generally recommended not to modify resources in this namespace.

3. **`kube-public`**:
   - **Purpose**: This namespace is created automatically and is readable by all users (including those not authenticated). Commonly used for resources that should be publicly accessible to all users.
   - **Usage**: Resources here can be publicly accessible and visible.

4. **`kube-node-lease`**:
   - **Purpose**: It holds lease objects associated with each node. These lease objects contain information about node health and are used by the control plane to detect node failures.

## DNS Within and Across Namespaces

Kubernetes DNS assigns a DNS name to every service, which includes the service name and namespace. This allows for easy discovery and communication between services.

### DNS Naming Convention

- Format: `<service-name>.<namespace-name>.svc.cluster.local`
- Example: For a service named `my-service` in the `my-namespace` namespace, the DNS name would be `my-service.my-namespace.svc.cluster.local`.

### Within a Namespace

- Services can be reached by other Pods within the same namespace using just the service name (e.g., `my-service`).

### Across Namespaces

- To access a service from a different namespace, you need to use the full DNS name.
- Example: A Pod in `namespace-a` accessing a service in `namespace-b` would use `service.namespace-b.svc.cluster.local`.

## Working with Namespaces in kubectl

### kubectl Commands with Namespace Flags

- **List Resources in a Specific Namespace**:
  ```bash
  kubectl get pods --namespace=<namespace-name>
  ```

- **Create Resources in a Specific Namespace**:
  ```bash
  kubectl create -f <file.yaml> --namespace=<namespace-name>
  ```

### Changing the Context Namespace

To avoid always specifying `--namespace=<namespace-name>` with every command, you can change the default namespace for your current context:

1. **Set a Namespace for Current Context**:
   ```bash
   kubectl config set-context --current --namespace=<namespace-name>
   ```
   This sets the default namespace for the current context to `<namespace-name>`.

2. **Check the Current Context**:
   ```bash
   kubectl config current-context
   ```
   Verify which context is currently active.

3. **View Configured Contexts and Namespaces**:
   ```bash
   kubectl config view
   ```
   This shows all configured contexts, including their default namespaces.


### Advanced Usage

- **Resource Quotas**: You can apply ResourceQuotas to a namespace to limit resource consumption.
- **Role-Based Access Control (RBAC)**: Use RBAC to define permissions on a per-namespace basis.


## Setting Resource Quotas in Namespaces

Resource Quotas are used to limit the amount of resources a namespace can consume. This ensures fair usage of resources and prevents any single namespace from consuming more than its allocated share.

### Example Configuration for Resource Quotas

**resource-quota.yaml**:
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: example-quota
  namespace: my-namespace
spec:
  hard:
    pods: "10"               # Maximum number of pods
    requests.cpu: "1"        # Total CPU request limit
    requests.memory: "1Gi"   # Total memory request limit
    limits.cpu: "2"          # Total CPU limit
    limits.memory: "2Gi"     # Total memory limit
```

#### Explanation of Each Field

- `apiVersion`, `kind`, `metadata`: Standard Kubernetes object metadata.
- `spec`:
  - `hard`: Defines the hard resource usage limits for the namespace.
    - `pods`: Maximum number of pods that can be created.
    - `requests.cpu`, `requests.memory`: Limits for total CPU and memory requests.
    - `limits.cpu`, `limits.memory`: Limits for total CPU and memory usage.

Apply this Resource Quota to a namespace with `kubectl apply -f resource-quota.yaml`.

## Configuring RBAC in Namespaces

RBAC in Kubernetes enables fine-grained control over what actions users and processes can perform on various resources in a namespace.

### Example Configuration for RBAC

Let’s define a Role and RoleBinding within a namespace:

1. **Role** (role.yaml): Defines permissions within a namespace.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: my-namespace
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```

- This Role `pod-reader` allows a user to perform `get`, `watch`, and `list` actions on Pods within `my-namespace`.

2. **RoleBinding** (rolebinding.yaml): Binds the Role to a user.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: my-namespace
subjects:
- kind: User
  name: jane # Name of a user
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

- This RoleBinding `read-pods` grants the `pod-reader` role to the user `jane` in `my-namespace`.

### Applying RBAC Configuration

Apply these configurations with:

```bash
kubectl apply -f role.yaml
kubectl apply -f rolebinding.yaml
```

### Considerations and Best Practices

- **Specificity in Roles**: Define Roles as specifically as possible to adhere to the principle of least privilege.
- **RoleBinding vs ClusterRoleBinding**: RoleBindings grant permissions in a specific namespace, whereas ClusterRoleBindings grant permissions cluster-wide.
- **Audit and Review**: Regularly audit and review RBAC configurations for compliance and security purposes.

By implementing Resource Quotas and RBAC, you can effectively manage and control resource consumption and access to resources in your Kubernetes namespaces, ensuring a secure and efficient environment.
