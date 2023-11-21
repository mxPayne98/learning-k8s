## Kubernetes StatefulSets

StatefulSets are Kubernetes objects used to manage stateful applications, providing stable, unique network identifiers, stable, persistent storage, and ordered, graceful deployment and scaling.

### How StatefulSets Work

- **Stable Pod Identity**: StatefulSet Pods have a unique, stable identity that is maintained across any rescheduling.
- **Ordered Operations**: StatefulSets ensure that Pods are created, updated, and deleted in a predictable order, typically useful for stateful applications like databases.
- **Persistent Storage**: Each Pod in a StatefulSet can be associated with its Persistent Volume, which remains attached to the Pod regardless of rescheduling.

### Example: Creating a StatefulSet for a Database

#### YAML Configuration for a StatefulSet

```yaml
# statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
          name: web
      volumeMounts:
      - name: www
        mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

- `serviceName`: The name of the Service that the StatefulSet is governing.
- `replicas`: The number of desired replicas.
- `selector`: Selector for finding the Pods.
- `template`: Template for the Pods.
- `volumeClaimTemplates`: A template for Persistent Volume Claims. Each Pod will get its own PVC according to this template.

### Commands Associated with StatefulSets

- **Create a StatefulSet**:
  ```bash
  kubectl apply -f statefulset.yaml
  ```
- **Get StatefulSet Information**:
  ```bash
  kubectl get statefulsets
  ```
- **Delete a StatefulSet**:
  ```bash
  kubectl delete statefulset web
  ```

### Best Practices and Gotchas

- **Stable Network Identities**: StatefulSets provide each Pod with a stable hostname based on the Pod’s ordinal index, which is important for applications that rely on stable identifiers.
- **Ordered Deployment and Scaling**: Be aware that StatefulSets do not provide any guarantees on the termination order of Pods. However, the creation order is guaranteed to be from `0` to `N-1`.
- **Update Strategies**: `RollingUpdate` (default) updates Pods in reverse ordinal order, while `OnDelete` requires manual deletion of Pods.
- **Persistence**: Stateful applications often require persistent storage, which should be carefully planned and managed.
- **Backup and Disaster Recovery**: Implementing robust backup and disaster recovery plans for stateful applications is crucial.


### Storage Techniques in Relation to StatefulSets

StatefulSets in Kubernetes can be configured to use storage in different ways, depending on the requirements of the application. One common scenario is where each Pod in the StatefulSet has its own storage (Persistent Volume), and another is where all Pods share the same storage volume.

#### Example 1: Unique Storage for Each Pod in StatefulSet

In many stateful applications like databases, each Pod in a StatefulSet requires its own storage. This is typically achieved using `volumeClaimTemplates`, which ensures that each Pod gets its own Persistent Volume Claim (PVC) and, consequently, its own Persistent Volume (PV).

**YAML Configuration for StatefulSet with Unique Storage per Pod**:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: db-statefulset
spec:
  serviceName: "db"
  replicas: 3
  selector:
    matchLabels:
      app: db
  template:
    metadata:
      labels:
        app: db
    spec:
      containers:
      - name: mysql
        image: mysql:5.7
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: mysql-data
          mountPath: /var/lib/mysql
  volumeClaimTemplates:
  - metadata:
      name: mysql-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

- `volumeClaimTemplates`: This creates a PVC for each Pod in the StatefulSet.
- Each Pod will have its own storage, independent of the other Pods.

#### Example 2: Shared Storage for All Pods in StatefulSet

There are scenarios where you might want all Pods in a StatefulSet to share the same storage. This can be achieved using a pre-provisioned PV or by manually provisioning a shared PV and binding it to a single PVC that is shared across all Pods.

**YAML Configuration for StatefulSet with Shared Storage**:

First, create a shared PersistentVolume (assuming a NFS server):

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: shared-nfs
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteMany
  nfs:
    server: nfs-server.example.com
    path: "/shared"
```

Then, create a PVC that binds to this PV:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: shared-nfs-claim
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 10Gi
  volumeName: shared-nfs
```

Finally, reference this PVC in your StatefulSet:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web-statefulset
spec:
  serviceName: "web"
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: nginx
        image: nginx
        volumeMounts:
        - name: shared-storage
          mountPath: /usr/share/nginx/html
  volumes:
  - name: shared-storage
    persistentVolumeClaim:
      claimName: shared-nfs-claim
```

- `volumes.persistentVolumeClaim.claimName`: This references the PVC that is bound to the shared PV.

### Commands Associated with StatefulSets and PVCs

- **Create Resources**:
  ```bash
  kubectl apply -f pv.yaml
  kubectl apply -f pvc.yaml
  kubectl apply -f statefulset.yaml
  ```
- **Get StatefulSet Information**:
  ```bash
  kubectl get statefulset
  ```
- **Delete a StatefulSet**:
  ```bash
  kubectl delete statefulset web-statefulset
  ```

### Best Practices and Gotchas

1. **Data Isolation**: Ensure that applications are compatible with shared storage models, especially regarding data isolation and concurrent access.
2. **Storage Access Modes**: The access mode of the PV must match the needs of your application. `ReadWriteMany` is necessary for a shared storage scenario.
3. **Capacity Planning**: Properly plan capacity for shared storage to prevent space issues, especially in a scenario where all Pods write to the same volume.
