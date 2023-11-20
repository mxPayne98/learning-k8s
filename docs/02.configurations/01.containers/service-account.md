## Kubernetes Service Accounts

Service Accounts in Kubernetes are used to provide an identity for processes that run in a Pod, enabling fine-grained control over access to the Kubernetes API and resources.

### How Service Accounts Work

1. **Automatic Provisioning**: By default, Kubernetes automatically provisions a default service account in each namespace.
2. **API Access**: Service accounts are tied to a set of credentials stored as Secrets, which can be used by applications inside a Pod to interact with the Kubernetes API.
3. **Role-Based Access Control (RBAC)**: You can grant permissions to a service account using RBAC policies, allowing you to define what actions the service account can perform in the Kubernetes cluster.

### Most Useful Commands

- **Creating a Service Account**:
  ```bash
  kubectl create serviceaccount <serviceaccount-name>
  ```
  Creates a new service account in the current namespace.

- **Listing Service Accounts**:
  ```bash
  kubectl get serviceaccounts
  ```
  Lists all service accounts in the current namespace.

- **Describing a Service Account**:
  ```bash
  kubectl describe serviceaccount <serviceaccount-name>
  ```
  Shows detailed information about a specific service account.

### YAML Configuration Example

**Service Account YAML Example**:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-service-account  # Name of the service account
```

#### Explanation of Each Field

- `apiVersion`: Version of the Kubernetes API.
- `kind`: Type of Kubernetes resource, here it is a ServiceAccount.
- `metadata`:
  - `name`: The name of the service account.

### Gotchas

- **Default Service Account**: Pods use the default service account in their namespace if no other service account is specified. This default account may have broader permissions than required.
- **Secrets Auto Mount**: By default, service account credentials are automatically mounted into Pods. This can be disabled if not needed.

### Using Service Accounts in Pods

To use a specific service account in a Pod, specify it in the Pod manifest:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  serviceAccountName: my-service-account  # Specifying the service account
  containers:
  - name: my-container
    image: nginx
```

- `serviceAccountName`: The name of the service account to use.

## Service Account Token API

The Service Account Token API is used to get bearer tokens associated with a service account. These tokens can authenticate against the Kubernetes API server.

### Commands to Get the Token

1. **Get the Name of the Secret**:
   First, identify the secret associated with the service account. This secret contains the token.
   ```bash
   kubectl get serviceaccount <serviceaccount-name> -o jsonpath="{.secrets[*].name}"
   ```

2. **Retrieve the Token from the Secret**:
   Use the secret name obtained from the previous step to get the token.
   ```bash
   kubectl get secret <secret-name> -o jsonpath="{.data.token}" | base64 --decode
   ```

This sequence of commands fetches the token that can be used for API authentication.

## Service Account Changes in Kubernetes 1.22 and 1.24

### Kubernetes 1.22

- **Bound Service Account Token Volumes**: This feature, which had been in beta, became the default behavior in Kubernetes 1.22. It improves the security of service account tokens by:
  - Ensuring they are only mountable by Pods that explicitly reference the service account.
  - Auto-rotating the tokens issued to a service account.
  - Allowing the specification of an audience (intended recipient) for the token.

### Kubernetes 1.24

- **Removal of Legacy Service Account Tokens**: Kubernetes 1.24 plans to remove the legacy service account token mechanism. This means:
  - Tokens will no longer be stored in Secrets by default.
  - The auto-mounting of the API credentials for a service account inside a Pod will be disabled by default for new clusters.

### Impact and Best Practices

- **Migration to Bound Tokens**: Users should ensure they are using bound service account tokens for better security. Applications relying on the automatic mounting of service account credentials need to be updated accordingly.
- **Managing Token Lifetime**: With the auto-rotation of tokens, it's crucial to handle token expiration and renewal in your applications.
- **Audience Specification**: When creating a service account token, you can specify the intended audience (`aud`) for the token, enhancing security by preventing the token's misuse.
