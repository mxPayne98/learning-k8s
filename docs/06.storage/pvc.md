## Kubernetes Persistent Volume Claims (PVCs)

Persistent Volume Claims (PVCs) in Kubernetes are a way for users to request storage resources from the cluster. PVCs abstract the details of how storage is provided and manage the consumption of Persistent Volumes (PVs).

### How PVCs Work

- **Purpose**: PVCs allow users to request specific sizes and access modes for storage, without needing to know the underlying storage system.
- **Binding**: When a PVC is created, it is automatically bound to a suitable Persistent Volume (PV) that meets its requirements. This binding is exclusive; a PV bound to one PVC cannot be bound to another PVC.

### PVC Lifecycle

1. **Provisioning**: A PVC is created with specific storage requirements. Depending on the cluster setup, a corresponding PV can be dynamically provisioned if a matching one does not exist.
2. **Binding**: The PVC is matched with an available PV and becomes bound to it. The matching is based on size, access modes, and user-specified labels or selectors.
3. **Using**: Pods can then use the PVC as a volume. The PVC's lifecycle is tied to the Pod's lifecycle, but the data can persist beyond it.
4. **Releasing**: When the Pod no longer needs the volume, the PVC can be deleted. The PV then goes into a `Released` state.
5. **Reclaiming**: Depending on the PV’s `persistentVolumeReclaimPolicy`, the underlying storage can be retained, recycled, or deleted.

### Example: Creating a PVC

#### YAML Configuration for a PVC

```yaml
# persistent-volume-claim.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: example-pvc
spec:
  accessModes:
    - ReadWriteOnce  # The PVC requires the PV to be mounted as read-write by a single node
  resources:
    requests:
      storage: 1Gi  # Requesting a PV with at least 1 GiB of storage
  storageClassName: standard  # Specifies the StorageClass
```

- `accessModes`: How the volume can be mounted (e.g., ReadWriteOnce, ReadOnlyMany, ReadWriteMany).
- `resources.requests.storage`: The minimum amount of storage requested.
- `storageClassName`: Specifies which StorageClass the PVC should use.

### Commands Associated with PVCs

- **Create a PVC**:
  ```bash
  kubectl apply -f persistent-volume-claim.yaml
  ```
- **Get PVC Information**:
  ```bash
  kubectl get pvc
  ```
- **Delete a PVC**:
  ```bash
  kubectl delete pvc example-pvc
  ```

### Gotchas

- **PVC-PV Binding**: Once bound, a PVC exclusively binds to a PV. This means the PV cannot be bound to another PVC and vice versa.
- **Storage Class Match**: Ensure that the storage class you specify in the PVC is available in your cluster and matches your requirements.
- **PVC Resizing**: While resizing PVCs is supported, not all storage classes or PV types support dynamic resizing.
- **Capacity Planning**: Be aware of your cluster’s storage capacity to avoid issues with PVCs failing to bind due to insufficient resources.

### More Gotchas in Persistent Volume Claims

#### 1. Capacity Requested Not Available

When a PVC requests a capacity that is not available, it affects the status and behavior of the PVC:

- **Status**: The PVC will remain in the `Pending` state until a matching PV becomes available.
- **Effect on Binding**: If no existing PV meets the capacity requirement and dynamic provisioning is not set up or unable to fulfill the request, the PVC will not be bound.
- **Manual Intervention**: An administrator might need to provision additional storage or adjust existing PVs to meet the capacity requirements.

#### 2. Access Modes Mismatch

Each PV has certain access modes, and a PVC must request an access mode that is supported by the PV:

- **Matching Requirement**: The PVC and PV must have matching access modes for a successful binding. Common access modes include `ReadWriteOnce`, `ReadOnlyMany`, and `ReadWriteMany`.
- **PV Availability**: If no PV with the requested access mode is available, the PVC will remain in a `Pending` state.
- **Dynamic Provisioning**: With dynamic provisioning, the storage class should support the requested access mode.

### Advanced Use Cases

#### Example: PVC for Stateful Applications

Stateful applications, like databases, often require persistent storage with specific access patterns.

**YAML Configuration for a Stateful Application PVC**:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: db-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: high-performance
```

- `accessModes: ReadWriteOnce`: Suitable for databases where a single instance requires read-write access.
- `storage: 10Gi`: Requests more storage for database needs.
- `storageClassName: high-performance`: Assumes a StorageClass `high-performance` that is tailored for database workloads.

#### Example: Configuring a MySQL Database with PVC

Setting up a MySQL database in Kubernetes involves creating a deployment that uses a PVC for persistent storage.

1. **Create a PVC for MySQL**:
   ```yaml
   apiVersion: v1
   kind: PersistentVolumeClaim
   metadata:
     name: mysql-pvc
   spec:
     accessModes:
       - ReadWriteOnce
     resources:
       requests:
         storage: 5Gi
   ```

2. **MySQL Deployment**:
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: mysql
   spec:
     selector:
       matchLabels:
         app: mysql
     strategy:
       type: Recreate
     template:
       metadata:
         labels:
           app: mysql
       spec:
         containers:
         - image: mysql:5.7
           name: mysql
           env:
           - name: MYSQL_ROOT_PASSWORD
             value: mypassword
           ports:
           - containerPort: 3306
             name: mysql
           volumeMounts:
           - name: mysql-persistent-storage
             mountPath: /var/lib/mysql
         volumes:
         - name: mysql-persistent-storage
           persistentVolumeClaim:
             claimName: mysql-pvc
   ```

   - The MySQL container uses the `mysql-pvc` for its data storage.
   - The deployment strategy `Recreate` ensures that the old Pod is removed before a new one is created, which is important for single-instance stateful applications like MySQL.

### Useful Commands

- **Create PVC and Deployment**:
  ```bash
  kubectl apply -f mysql-pvc.yaml
  kubectl apply -f mysql-deployment.yaml
  ```
- **Check PVC Status**:
  ```bash
  kubectl get pvc mysql-pvc
  ```
- **Check Pod and Deployment Status**:
  ```bash
  kubectl get pods
  kubectl get deployment mysql
  ```
