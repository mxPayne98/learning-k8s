## Setting up a MySQL cluster with 1 master and 2 slave nodes using StatefulSets

### Step 1: Define Storage Classes

Assuming dynamic provisioning with AWS storage, define storage classes. Here, we’ll use one storage class, but in a production environment, you might want different classes for master and slaves.

**Storage Class YAML (e.g., `aws-storage-class.yaml`):**

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: aws-gp2
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
reclaimPolicy: Retain
allowVolumeExpansion: true
```

Apply the storage class:

```bash
kubectl apply -f aws-storage-class.yaml
```

### Step 2: Create a Headless Service for StatefulSet

A headless service is needed for the StatefulSet to control the network domain.

**Headless Service YAML (e.g., `mysql-headless-service.yaml`):**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
spec:
  clusterIP: None
  selector:
    app: mysql
  ports:
    - port: 3306
      name: mysql
```

Apply the service:

```bash
kubectl apply -f mysql-headless-service.yaml
```

### Step 3: Create a StatefulSet for MySQL Cluster

**StatefulSet YAML (e.g., `mysql-statefulset.yaml`):**

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: "mysql"
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:5.7
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "yourpassword"
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: mysql-data
          mountPath: /var/lib/mysql
  volumeClaimTemplates:
  - metadata:
      name: mysql-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: "aws-gp2"
      resources:
        requests:
          storage: 10Gi
```

Apply the StatefulSet:

```bash
kubectl apply -f mysql-statefulset.yaml
```

### Step 4: Setting Up MySQL Replication

The setup of MySQL replication is generally done within the MySQL environment and requires executing SQL commands in the MySQL shell. Here's a high-level overview:

1. **On the Master**: Configure MySQL to generate binary logs.
   - `SET GLOBAL server_id = 1;`
   - `FLUSH TABLES WITH READ LOCK;`
   - `SHOW MASTER STATUS;` (Note the log file and position)
   - Unlock tables after backing up data.

2. **On Slave 1**: Clone data from the master.
   - `CHANGE MASTER TO MASTER_HOST='mysql-0.mysql', MASTER_USER='replica', MASTER_PASSWORD='password', MASTER_LOG_FILE='noted-log-file', MASTER_LOG_POS=noted-log-position;`
   - `START SLAVE;`

3. **On Slave 2**: Clone data from Slave 1.
   - Repeat the process, pointing to Slave 1 as the master.

### Additional Steps

- **Monitoring and Managing Replication**: Use MySQL's tools and commands to monitor the replication status and troubleshoot issues.
- **Database Initialization and User Setup**: These are typically part of the database setup process, outside the scope of the Kubernetes setup.

### Best Practices and Gotchas

- **Persistent Storage**: Ensure that the storage provisioned is sufficient and performs well with MySQL, especially under load.
- **StatefulSet Scaling**: Be cautious when scaling StatefulSets, as this can impact the replication setup.
- **Security**: Secure your MySQL root and replication user credentials. Using Kubernetes Secrets is advisable for managing sensitive data.
- **Backup and Recovery**: Plan for backup and recovery of your MySQL data.


### Backup Strategy for MySQL on Kubernetes

1. **Database-Level Backups**: Use MySQL tools like `mysqldump` for logical backups or tools like Percona XtraBackup for physical backups. These backups can be scheduled as CronJobs in Kubernetes.

2. **Volume Snapshots**: If using persistent volumes, you can leverage volume snapshot features provided by your cloud provider or storage solution.

3. **Backup to External Storage**: Backups should ideally be stored off-cluster, for example, in an S3 bucket or other cloud storage services, for added safety.

### Example: Kubernetes CronJob for MySQL Backup

You can create a Kubernetes CronJob to take regular backups using `mysqldump`.

#### CronJob YAML Configuration

```yaml
# mysql-backup-cronjob.yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: mysql-backup
spec:
  schedule: "0 2 * * *"  # Schedule backup at 2:00 AM every day
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: mysql-backup
            image: mysql:5.7
            command:
            - /bin/bash
            - -c
            - |
              mysqldump -h mysql-0.mysql -u root -p yourpassword --all-databases > /backup/mysql-backup-$(date +%Y%m%d).sql
            volumeMounts:
            - name: backup-volume
              mountPath: /backup
          volumes:
          - name: backup-volume
            persistentVolumeClaim:
              claimName: mysql-backup-pvc  # PVC for backup data
          restartPolicy: OnFailure
```

- `schedule`: Defines when the backup job should run.
- `mysqldump`: Command to take a backup of all databases. Adjust the command to suit your needs (e.g., specific databases).
- `volumeMounts`: Mount the backup volume where the backup file will be stored.

#### PersistentVolumeClaim for Backup Data

Create a PVC to store backup data. This PVC should ideally be backed by a storage class that replicates data across zones or offsite.

```yaml
# backup-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-backup-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: "aws-gp2"  # Example: AWS GP2 storage class
  resources:
    requests:
      storage: 10Gi
```

Apply the configurations:

```bash
kubectl apply -f backup-pvc.yaml
kubectl apply -f mysql-backup-cronjob.yaml
```

### Best Practices and Gotchas

- **Backup Frequency**: Determine the backup frequency based on your data criticality and change rate.
- **Backup Testing**: Regularly test backups to ensure they can be restored successfully.
- **Security**: Secure your backups, especially if they contain sensitive data. Consider encryption and secure access policies.
- **Monitoring and Alerts**: Set up monitoring for the backup process and alerts for failures or issues.
