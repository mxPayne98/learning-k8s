## Environment, ConfigMaps, and Secrets in Kubernetes

Kubernetes offers several ways to manage configuration data and sensitive information through Environment Variables, ConfigMaps, and Secrets.

### Environment Variables

Environment Variables are a simple way to pass configuration to your containers.

#### Ways to Create and Inject

- **Directly in Pod Spec**:
  ```yaml
  env:
    - name: ENV_VAR_NAME
      value: "value"
  ```

- **From ConfigMap or Secret**:
  ```yaml
  env:
    - name: ENV_VAR_NAME
      valueFrom:
        configMapKeyRef:
          name: configmap-name
          key: key-name
  ```

### ConfigMaps

ConfigMaps are used to store non-confidential data in key-value pairs. They can be used to store configuration files, command-line arguments, environment variables, port numbers, etc.

#### Ways to Create

- **Imperative Command**:
  ```bash
  kubectl create configmap <configmap-name> --from-literal=key1=value1 --from-file=path/to/file
  ```

- **Using a YAML File**:
  ```yaml
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: example-configmap  # Name of the ConfigMap
  data:
    key1: value1             # Key-value pairs
    key2: value2
  ```

  - `apiVersion`: Version of the Kubernetes API.
  - `kind`: Specifies the resource type, here it's a ConfigMap.
  - `metadata`: Metadata about the resource, including the name.
  - `data`: The data stored in the ConfigMap, provided as key-value pairs.

#### Injecting Values

- **As Environment Variables**:
  ```yaml
  envFrom:
    - configMapRef:
        name: <configmap-name>
  ```

- **As Volumes**:
  Mount ConfigMap as a volume in the Pod:
  ```yaml
  volumes:
    - name: config-volume
      configMap:
        name: <configmap-name>
  ```

### Secrets

Secrets are used to store and manage sensitive information, such as passwords, OAuth tokens, and SSH keys.

#### Ways to Create

- **Imperative Command**:
  ```bash
  kubectl create secret generic <secret-name> --from-literal=key1=value1
  ```

- **Using a YAML File**:
  Data in Secrets must be base64 encoded:
  ```yaml
  apiVersion: v1
  kind: Secret
  metadata:
    name: example-secret     # Name of the Secret
  type: Opaque              # Type of Secret; Opaque is the default
  data:
    key1: YmFzZTY0RW5jb2RlZFZhbHVl  # Base64 encoded values
  ```

  - `type`: The type of the Secret; `Opaque` is used for arbitrary user-defined data.
  - `data`: The data stored in the Secret, where each value must be a base64-encoded string.

#### Injecting Values

- **As Environment Variables**:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
spec:
  containers:
  - name: example-container
    image: nginx
    envFrom:
    - secretRef:
        name: example-secret   # Referencing the entire Secret
```

- `envFrom`: Defines environment variables based on values from resources.
- `secretRef`: Specifies the Secret resource to expose as environment variables.
- `name`: Name of the Secret to reference.

- **As Single Environment Variables**:
  ```yaml
  env:
    - name: SECRET_KEY
      valueFrom:
        secretKeyRef:
          name: <secret-name>
          key: key1
  ```

- **As Volumes**:
  Mount Secrets as a volume in the Pod:
  ```yaml
  volumes:
    - name: secret-volume
      secret:
        secretName: <secret-name>
  ```

### Note about Secrets

- A secret is only sent to a node if a pod on that node requires it.

- Kubelet stores the secret into a tmpfs so that the secret is not written to disk storage.

- Once the Pod that depends on the secret is deleted, kubelet will delete its local copy of the secret data as well.


### Encryption of Secrets

- **By Default**: Kubernetes Secrets are stored as base64-encoded strings. By default, they are not encrypted at rest.
- **Encryption at Rest**: Kubernetes supports [Encryption at Rest for Secrets](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/). This can be enabled by configuring the API server with encryption providers.
- **Using Cloud Providers**: Cloud providers like AWS offer their own solutions for secret management (e.g., AWS KMS) that can be integrated with Kubernetes for enhanced security.

#### Example: Integration with AWS KMS

- **AWS KMS**: Kubernetes can be configured to use AWS Key Management Service (KMS) for encrypting data at rest. This requires setting up the KMS plugin and modifying the API server's configuration to use this plugin for encryption.

## Encrypting Secrets at Rest

### Native Encryption in Kubernetes

Kubernetes supports encryption of Secrets at rest natively using an EncryptionConfiguration object.

**Steps**:

1. **Create an Encryption Configuration File**:
   This file specifies the encryption provider to be used. For example:

   ```yaml
   apiVersion: apiserver.config.k8s.io/v1
   kind: EncryptionConfiguration
   resources:
     - resources:
         - secrets
       providers:
         - aescbc:
             keys:
               - name: key1
                 secret: <base64-encoded-secret>
         - identity: {}
   ```

   Replace `<base64-encoded-secret>` with a base64-encoded AES key.

2. **Start the API Server with Encryption Configuration**:
   Pass the encryption configuration file to the Kubernetes API server using the `--encryption-provider-config` flag.

### Using AWS KMS

To use AWS Key Management Service (KMS) for encrypting Secrets:

1. **Set Up an AWS KMS Key**:
   - Create a customer master key (CMK) in AWS KMS.

2. **Configure the Kubernetes API Server**:
   - Use a KMS plugin for Kubernetes, like [aws-encryption-provider](https://github.com/kubernetes-sigs/aws-encryption-provider), which integrates with AWS KMS.
   - The API server must be configured to use this plugin as an encryption provider.

3. **Encryption Configuration with KMS**:
   - The EncryptionConfiguration file should specify the KMS provider:

     ```yaml
     apiVersion: apiserver.config.k8s.io/v1
     kind: EncryptionConfiguration
     resources:
       - resources:
           - secrets
         providers:
           - kms:
               name: aws-kms
               endpoint: unix:///var/run/kmsplugin/socket.sock
               cachesize: 1000
               timeout: 3s
           - identity: {}
     ```

   - The `endpoint` is where the KMS plugin is running.
