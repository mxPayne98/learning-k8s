## Kubernetes Deployments

Deployments in Kubernetes are a higher-level concept that manages Pods and ReplicaSets. They provide declarative updates for Pods and ReplicaSets, which makes them very powerful and convenient for managing application deployments.

### How Deployments Work

- **Manages Pods and ReplicaSets**: A Deployment provides declarative updates to Pods and ReplicaSets. You describe a desired state in a Deployment, and the Deployment Controller changes the actual state to the desired state at a controlled rate.
- **Rolling Updates & Rollbacks**: They allow for easy updating of applications. When you update a Deployment, it creates a new ReplicaSet and scales it up while scaling down the old ReplicaSet. If something goes wrong, you can easily roll back to a previous version of the Deployment.
- **Self-Healing**: Similar to ReplicaSets, if a Pod in a Deployment is deleted or terminated, the Deployment will create a new one to replace it.
- **Scalability**: You can easily scale up or down the number of Pods via the Deployment.

### Differences from ReplicaSets and Replication Controllers

- **Update Mechanism**: Deployments provide a more sophisticated update mechanism than ReplicaSets or Replication Controllers. They manage the rollout of updates to the Pods and can handle rollback scenarios.
- **Use of ReplicaSets**: Deployments use ReplicaSets under the hood for managing Pod replicas, which provides more advanced features like rolling updates and rollbacks compared to Replication Controllers.

### Gotchas

- **Immutable Pod Template**: Once a Deployment is created, its Pod template (`.spec.template`) is immutable. To update Pods, you need to create a new Deployment revision.
- **Overlapping Deployments**: Be careful with label selectors that might overlap between different Deployments, as it can lead to unexpected behavior.

### Example YAML Configuration

Here's an example of a Deployment YAML configuration for the Node.js application:

**node-deployment.yaml**:
```yaml
apiVersion: apps/v1              # API version
kind: Deployment                 # Kubernetes resource type
metadata:
  name: node-deployment          # Name of the Deployment
spec:
  replicas: 3                    # Number of desired Pods
  selector:
    matchLabels:
      app: node-app              # Selector to match Pods
  template:                      # Pod template
    metadata:
      labels:
        app: node-app            # Labels applied to Pods
    spec:
      containers:
      - name: node-container     # Name of the container in the Pod
        image: [your-docker-username]/node-redis-app:latest # Docker image
        ports:
        - containerPort: 8080    # Port the container exposes
```

### Most Useful Commands

1. **Creating a Deployment**:
   ```bash
   kubectl apply -f node-deployment.yaml
   ```

2. **Listing Deployments**:
   ```bash
   kubectl get deployments
   ```

3. **Viewing Deployment Details**:
   ```bash
   kubectl describe deployment <deployment-name>
   ```

4. **Updating a Deployment**:
   - Update the YAML file and then run `kubectl apply -f node-deployment.yaml`.
   - Alternatively, use `kubectl set image` for image updates.

5. **Scaling a Deployment**:
   ```bash
   kubectl scale deployment <deployment-name> --replicas=<number>
   ```

6. **Rolling Back a Deployment**:
   ```bash
   kubectl rollout undo deployment <deployment-name>
   ```

7. **Viewing Rollout History**:
   ```bash
   kubectl rollout history deployment <deployment-name>
   ```

   > **NOTE**:  History shows CHANGE-CAUSE for an update. While updating (creating/editing etc.) a deployment use flag `--record` to record the cause for the deployment update.

8. **Deleting a Deployment**:
   ```bash
   kubectl delete deployment <deployment-name>
   ```

Deployments in Kubernetes are incredibly powerful, especially when it comes to managing the rollout of updates and handling rollbacks. Let's explore these aspects in more detail.

### Update Mechanisms in Deployments

1. **Rolling Update (Default Strategy)**:
   - **How it Works**: When you update the Pod template (e.g., update the container image), the Deployment creates a new ReplicaSet and gradually scales up the new ReplicaSet while scaling down the old one.
   - **Zero Downtime**: This ensures zero downtime and guarantees that a certain number of Pods are always available.

2. **Recreate Strategy**:
   - **Usage**: You can also set the Deployment strategy to `Recreate`, which terminates all existing Pods before new ones are created.
   - **Downtime Involved**: This approach will cause downtime but is useful when you cannot have two versions of an application running simultaneously.

### Rollouts in Deployments

- **Managed Rollouts**: When you update a Deployment, a rollout is triggered. This rollout is a managed deployment of the new version.
- **Rollout Status**: You can check the status of a rollout with the command:
  ```bash
  kubectl rollout status deployment/<deployment-name>
  ```

### Rollbacks in Deployments

- **Automatic Rollbacks**: If Kubernetes detects something wrong during the rollout (like failing health checks), it can automatically rollback to the previous version.
- **Manual Rollbacks**: You can manually trigger a rollback to a previous revision of the Deployment using:
  ```bash
  kubectl rollout undo deployment/<deployment-name>
  ```

### Commands for Managing Rollouts and Rollbacks

1. **Viewing Rollout History**:
   ```bash
   kubectl rollout history deployment/<deployment-name>
   ```
   - This command shows you the history of rollouts, including revisions.

2. **Rolling Back to a Specific Revision**:
   ```bash
   kubectl rollout undo deployment/<deployment-name> --to-revision=<revision-number>
   ```
   - This rolls back to a specific revision number.

3. **Pausing and Resuming a Rollout**:
   - **Pausing**: 
     ```bash
     kubectl rollout pause deployment/<deployment-name>
     ```
     This pauses the rollout of a Deployment, allowing you to apply multiple fixes or changes without triggering multiple rollouts.
   - **Resuming**:
     ```bash
     kubectl rollout resume deployment/<deployment-name>
     ```
     This resumes a paused rollout.

### Complex Concepts in Deployments

1. **Proactive Rollback**:
   - Deployments can be configured with readiness probes and health checks. If a new version starts failing these checks after a rollout, Kubernetes can automatically rollback.

2. **Canary Deployments** (Advanced Use Case):
   - This is not natively supported by basic Deployments, but you can simulate it by controlling the rollout process, where you update a Deployment in a controlled manner (like updating a small percentage of Pods at a time) to minimize the risk.

3. **Blue/Green Deployment** (Advanced Use Case):
   - This is a deployment strategy where you have two identical environments: Blue (current) and Green (new). Once the Green environment is tested and ready, the traffic is switched from Blue to Green. This can be managed by manipulating services and labels in Kubernetes but might require additional tooling for smoother operations.

4. **Resource Limits and Requests**:
   - During rollouts, ensure that resource requests and limits are appropriately set to prevent resource starvation and ensure smooth scaling.

Setting up Canary and Blue/Green deployment strategies in Kubernetes requires a bit of planning and careful configuration. These strategies are not directly built into Kubernetes but can be implemented using its features. Let's go through each of these deployment strategies.

## 1. Canary Deployments

Canary deployments are used to roll out updates to a small subset of users before making them available to everybody. This approach helps in minimizing the impact of any potential issues.

### Implementing a Canary Deployment

1. **Deploy the Initial Version**:
   First, you have your stable application version deployed.

2. **Deploy the Canary Version**:
   Next, you deploy the new version of your application as a separate deployment with a small number of replicas.

3. **Route Traffic**:
   Use Kubernetes services or an Ingress controller that can route a percentage of traffic to the canary deployment. This could be based on HTTP headers, cookies, or other criteria.

4. **Monitor and Scale**:
   Monitor the performance and behavior of the canary. If it performs well, you can gradually increase the traffic to it while decreasing traffic to the stable version.

### Example Configuration

Assuming you have two deployments: `app-stable` for the stable version and `app-canary` for the canary version.

**app-stable-deployment.yaml**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-stable
spec:
  replicas: 5
  ...
```

**app-canary-deployment.yaml**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-canary
spec:
  replicas: 1 # Start with a smaller number of replicas
  ...
```

You would then configure your service or Ingress to route a portion of the traffic to `app-canary`.

## 2. Blue/Green Deployments

Blue/Green deployment is a technique to reduce downtime and risk by running two identical environments: one hosting the current version of the application (Blue) and one hosting the new version (Green).

### Implementing a Blue/Green Deployment

1. **Deploy the Blue Environment**:
   This is your current production environment.

2. **Prepare the Green Environment**:
   Deploy the new version of your application in parallel to the Blue environment. This should be identical to the Blue environment except for the new changes.

3. **Switch Traffic**:
   Once you're ready to release, switch the service endpoint from the Blue to the Green environment. This is often done by updating the service selector in Kubernetes to point to the Green deployment.

4. **Monitor and Finalize**:
   After the switch, monitor the Green environment. If something goes wrong, you can quickly switch back to the Blue environment.

### Example Configuration

Assume you have two identical deployments: `app-blue` (current version) and `app-green` (new version). The service initially points to `app-blue`.

**app-blue-green-service.yaml**:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: app-service
spec:
  selector:
    app: app-blue # Initially points to the Blue environment
  ...
```

To switch to the Green environment, update the service selector:

```yaml
spec:
  selector:
    app: app-green # Now points to the Green environment
```

### Considerations

- **Canary Deployments** are great for user-based testing and gradual rollouts.
- **Blue/Green Deployments** are excellent for instant rollbacks and zero-downtime deployments.

Setting up routing rules and implementing monitoring strategies are crucial components of effectively managing Canary and Blue/Green deployments in Kubernetes. Let's delve into these aspects.

## Setting Up Routing Rules

Routing rules in Kubernetes are primarily managed through Services and Ingress controllers. The routing for Canary and Blue/Green deployments depends on how you manage and direct traffic to different versions of your application.

### 1. **Using Kubernetes Services**

For both Canary and Blue/Green deployments, Kubernetes Services can route traffic to different sets of Pods based on labels.

- **Label Selectors**: Update the labels and selectors in your Service definition to route traffic to different deployments.
- **Example**: In a Blue/Green deployment, change the selector in the Service configuration to switch from blue to green Pods.

### 2. **Using Ingress Controllers**

For more sophisticated routing (like splitting traffic between versions), an Ingress controller like NGINX or Traefik can be used.

- **Weighted Routing**: Implement weighted routing where a certain percentage of traffic is directed to the Canary version while the rest goes to the stable version.
- **Header-Based Routing**: Route traffic based on HTTP headers or cookies, which is useful for directing traffic to different versions based on user or session information.

### Configuration Example for Canary Deployment

Suppose you have an Ingress controller like NGINX:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: myapp.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-stable
            port:
              number: 80
      - path: /canary
        pathType: Prefix
        backend:
          service:
            name: app-canary
            port:
              number: 80
```

In this example, traffic going to `myapp.com/canary` will be routed to the `app-canary` service.

### Splitting Traffic by Percentage

To split traffic by percentage between different versions of your application (e.g., stable and canary), you can use the `nginx.ingress.kubernetes.io/canary` and `nginx.ingress.kubernetes.io/canary-weight` annotations.

Here's an example configuration:

1. **Stable Version Ingress**:
   No changes are needed to the Ingress resource for your stable version.

2. **Canary Version Ingress**:
   You create a separate Ingress resource for the canary version and use annotations to define the weight (percentage of traffic to send to the canary).

   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: canary-app-ingress
     annotations:
       nginx.ingress.kubernetes.io/canary: "true"
       nginx.ingress.kubernetes.io/canary-weight: "20"
   spec:
     ...
   ```

   In this example, `20%` of the traffic will be routed to the canary version.

### Routing Based on HTTP Header

You can route traffic based on HTTP headers by using the `nginx.ingress.kubernetes.io/canary-by-header` annotation.

Here's how you can set it up:

1. **Canary Version Ingress**:
   Define a specific header and value for routing traffic to the canary version.

   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: canary-app-ingress
     annotations:
       nginx.ingress.kubernetes.io/canary: "true"
       nginx.ingress.kubernetes.io/canary-by-header: "Canary-Version"
       nginx.ingress.kubernetes.io/canary-by-header-value: "v2"
   spec:
     ...
   ```

   In this configuration, only requests with a header `Canary-Version: v2` will be routed to the canary version.

### Considerations and Tips

- **Consistent Hashing**: If you are load balancing requests that need session persistence, ensure that the NGINX Ingress controller is configured for consistent hashing.
- **Versioning Headers**: It's a good practice to use version-specific headers or values for routing to different deployment versions.
- **Monitoring and Logging**: Make sure to monitor and log the traffic distribution to quickly identify any issues or unexpected behavior.

## Monitoring Strategies

Effective monitoring is essential for Canary and Blue/Green deployments to quickly identify and respond to issues.

### 1. **Metrics and Performance Monitoring**

- **Tools**: Use monitoring tools like Prometheus and Grafana to collect and visualize metrics from your Pods and services.
- **Key Metrics**: Monitor CPU, memory usage, response times, error rates, and other relevant metrics.
- **Alerts**: Set up alerts for anomalous behavior or thresholds that indicate problems with the new release.

### 2. **Log Analysis**

- **Centralized Logging**: Implement a centralized logging solution like ELK (Elasticsearch, Logstash, Kibana) stack or Fluentd to aggregate and analyze logs from all Pods.
- **Correlation**: Correlate logs between old and new versions to identify specific issues related to the deployment.

### 3. **Health Checks and Readiness Probes**

- **Liveness Probes**: Ensure that the application is running and restart it if it fails.
- **Readiness Probes**: Ensure that the application is ready to serve traffic. This is crucial for smooth rolling updates and avoiding routing traffic to Pods that are not ready.

### 4. **User Feedback**

- In Canary deployments, consider incorporating user feedback mechanisms to gauge the impact of new features or changes.

### 5. **Real-time Monitoring**

- Use real-time monitoring tools to observe the system behavior as soon as the traffic starts hitting the new version.
