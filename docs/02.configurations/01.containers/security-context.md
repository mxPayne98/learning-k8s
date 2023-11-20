## Kubernetes Security Contexts

Security Contexts in Kubernetes are used to define privilege and access control settings for Pods and containers. They allow you to control security-sensitive aspects of the container's environment.

### How Security Contexts Work

- **Pod-Level Security Context**: Applies to all containers in the Pod.
- **Container-Level Security Context**: Applies to a specific container. Overrides Pod-level settings if both are defined.

### Key Fields in Security Context

- `runAsUser`: The UID to run the entrypoint of the container process.
- `runAsGroup`: The GID to run the entrypoint of the container process.
- `fsGroup`: The GID associated with the container's filesystem.
- `privileged`: Determines if any container in a Pod can enable privileged mode.
- `readOnlyRootFilesystem`: Whether the container's root filesystem is read-only.
- `allowPrivilegeEscalation`: Controls whether a process can gain more privileges than its parent process.

### Most Useful Commands

- **Apply Pod/Deployment with Security Context**:
  ```bash
  kubectl apply -f <filename>.yaml
  ```
  This applies the configuration, including the security context defined in the YAML file.

### YAML Configuration Examples

1. **Pod-Level Security Context**:

    **pod-security-context.yaml**:
    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: example-pod
    spec:
      securityContext:
        runAsUser: 1000       # UID
        fsGroup: 2000         # GID for file system
      containers:
      - name: example-container
        image: nginx
    ```

    - Sets the user ID and file system group ID for all containers in the Pod.

2. **Container-Level Security Context**:

    **container-security-context.yaml**:
    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: example-pod
    spec:
      containers:
      - name: example-container
        image: nginx
        securityContext:
          runAsUser: 1000       # UID
          readOnlyRootFilesystem: true  # Root filesystem is read-only
    ```

    - Sets specific security settings for the `example-container` container.

### Gotchas

- **Compatibility with Volumes**: Ensure that the `fsGroup` and `runAsUser` settings are compatible with the permissions and ownership requirements of the volumes used by the Pod.
- **Non-Root Containers**: Running containers as non-root is a best practice. However, some legacy applications might not function correctly without root privileges.
- **Privileged Mode**: Use caution with `privileged` containers as they have elevated access and can potentially access the host's resources.

