## Kubernetes Resource Requirements

Resource requirements in Kubernetes allow you to specify how much CPU and memory (RAM) each container in a Pod needs. Understanding and correctly setting these requirements is crucial for the efficient and stable operation of your applications in a Kubernetes environment.

### How Resource Requirements Work

- **`requests`**: The amount of resources Kubernetes will guarantee for a container. If a container exceeds its request for a resource, it can use more up to the limit.
- **`limits`**: The maximum amount of a resource that a container can use. If a container exceeds its limit, it might be terminated or throttled.

### YAML Configuration Example

**Example Pod with Resource Requirements**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: resource-demo-pod
spec:
  containers:
  - name: demo-container
    image: nginx
    resources:
      requests:
        memory: "64Mi"       # Memory request
        cpu: "250m"          # CPU request (250 millicores)
      limits:
        memory: "128Mi"      # Memory limit
        cpu: "500m"          # CPU limit (500 millicores)
```

#### Explanation of Each Field

- `resources`:
  - `requests`:
    - `memory`: The amount of RAM the container is guaranteed. `64Mi` means 64 Mebibytes.
    - `cpu`: Guaranteed CPU. `250m` means 250 millicores, or 0.25 of a CPU.
  - `limits`:
    - `memory`: The maximum amount of RAM the container can use.
    - `cpu`: The maximum CPU the container can use.

### Most Useful Commands

- **Create a Pod with Resource Requirements**:
  ```bash
  kubectl apply -f <pod-spec.yaml>
  ```

- **Describe Pod to See Resources**:
  ```bash
  kubectl describe pod <pod-name>
  ```

### Gotchas and Best Practices

1. **Resource Exhaustion**:
   - Pods requesting more resources than available can lead to resource exhaustion and scheduling failures.
   - Always ensure that cluster resources match or exceed the sum of all resource requests.

2. **CPU and Memory Units**:
   - Understand the units for CPU (cores in millicores) and memory (bytes in Mi/Gi).
   - `1` CPU is equivalent to 1 AWS vCPU or 1 GCP Core.

3. **No Limits Specified**:
   - Containers without resource limits can potentially consume all available resources on a node, impacting other containers.

4. **Requests vs. Limits**:
   - Setting resource `requests` too low can lead to poor performance.
   - Setting `limits` too high can lead to inefficient resource utilization.

5. **Overcommitting Resources**:
   - Overcommitting resources (requests < limits) is acceptable but monitor performance and stability closely.
   - Overcommitting too much can lead to resource contention and degraded performance.

6. **Quality of Service (QoS) Classes**:
   - Kubernetes uses resource requests and limits to determine the QoS class for a Pod:
     - `Guaranteed`: If every container in the Pod has the same resource request and limit.
     - `Burstable`: If the requests and limits are not the same for all containers.
     - `BestEffort`: If no requests or limits are specified in the Pod spec.

7. **Node Pressure Eviction**:
   - Pods might be evicted under node pressure if they exceed their requests and the node is running out of resources.

8. **Horizontal Pod Autoscaler (HPA)**:
   - Consider using HPA to automatically scale applications based on observed CPU utilization or other select metrics.

### Example Scenarios

- **High Traffic Web Application**:
  - Set higher CPU and memory limits to handle traffic spikes, but monitor resource usage to right-size requests and limits.
- **Batch Processing Jobs**:
  - Allocate higher CPU for faster processing but be wary of resource contention with other critical applications.

### Best Practices for Setting Limits

1. **Setting Both CPU and Memory Limits**:
   - It's recommended to set both CPU and memory limits. This helps to prevent a container from using more than its share of resources and impacting other containers or the node itself.
   - CPU Limits: Prevents a container from using excessive CPU, which could starve other containers.
   - Memory Limits: Prevents a container from consuming too much memory, which could lead to system instability and killing of containers.

2. **Considerations**:
   - Understand the workload: For CPU-bound applications, CPU limits are crucial. For applications that consume more memory, focus on memory limits.
   - Avoid too restrictive limits, which can lead to poor performance and application crashes.

### Actions When a Pod Reaches its Limits

1. **CPU Limit Reached**:
   - The container is throttled, and its CPU usage is restricted. It won’t be killed for excessive CPU usage but will run slower.
   - Action: Monitor the application's performance. If it's degraded, consider increasing the CPU limit.

2. **Memory Limit Reached**:
   - The container might be terminated or restarted if it exceeds its memory limit.
   - Action: Investigate memory usage and leaks. Increase memory limits if necessary and feasible.

### Commands for Monitoring and Adjusting Resources

- **Monitoring Resource Usage**:
  ```bash
  kubectl top pod <pod-name>
  ```
  Displays the resource usage for Pods.

- **Updating Resource Limits**:
  ```bash
  kubectl edit pod <pod-name>
  ```
  Allows you to edit and update resource limits directly.

### Scaling Up Nodes in AWS EKS

When a node hits its CPU and memory limits, and the cluster gets affected, you might need to scale up your nodes:

1. **Using EKS Managed Node Groups**:
   - If using managed node groups, update the desired count:
     ```bash
     eksctl scale nodegroup --cluster=<cluster-name> --name=<nodegroup-name> --nodes=<desired-count>
     ```
   - This command increases the number of nodes in the specified node group.

2. **Using Auto Scaling Groups**:
   - If nodes are part of an Auto Scaling Group (ASG), adjust the desired capacity of the ASG in the AWS Management Console or using AWS CLI:
     ```bash
     aws autoscaling set-desired-capacity --auto-scaling-group-name <asg-name> --desired-capacity <new-desired-capacity>
     ```

3. **Cluster Autoscaler**:
   - For automatic scaling, consider using the Cluster Autoscaler. It automatically adjusts the size of the Kubernetes cluster when there are insufficient resources or nodes are underutilized.


### Setting Up Horizontal Pod Autoscaler (HPA)

HPA automatically scales the number of Pods in a deployment based on observed CPU utilization or other selected metrics.

#### Creating an HPA

1. **Example HPA YAML**:
   ```yaml
   apiVersion: autoscaling/v2beta2
   kind: HorizontalPodAutoscaler
   metadata:
     name: my-hpa
   spec:
     scaleTargetRef:
       apiVersion: apps/v1
       kind: Deployment
       name: my-deployment
     minReplicas: 1
     maxReplicas: 10
     metrics:
     - type: Resource
       resource:
         name: cpu
         target:
           type: Utilization
           averageUtilization: 50
   ```

   - This HPA targets the `my-deployment` Deployment and scales the number of Pods between 1 and 10, aiming for an average CPU utilization of 50%.

2. **Applying HPA**:
   ```bash
   kubectl apply -f hpa.yaml
   ```

### Best Practices for Using HPA

- **Understand Workload Patterns**: Configure HPA thresholds based on your application's normal workload patterns.
- **Right Metrics**: Besides CPU, consider other metrics like memory usage or custom metrics that better represent your application's load.
- **Testing**: Test HPA configurations in a staging environment to understand how it responds to different load conditions.
- **Alerting**: Set up monitoring and alerts for HPA actions to keep track of scaling events.
- **Setting Appropriate Requests**: It's crucial to set realistic CPU requests for your Pods. If the request is set too low, the Pods might be throttled unnecessarily. If set too high, it might lead to underutilization and less efficient scaling.
- **Observing Behavior**: Monitoring the behavior of your HPA and the performance of your applications will help you fine-tune the CPU requests and HPA settings.

In the context of a Horizontal Pod Autoscaler (HPA) in Kubernetes, `averageUtilization` within the `target` field refers to the average CPU utilization as a percentage of the **requested** CPU resources for the Pods.

### Understanding `averageUtilization`:

- **Relative to CPU Requests**: The `averageUtilization` is calculated relative to the CPU requests of the Pods, not their limits or the node's total CPU capacity.
- **How It Works**: If the HPA is set with an `averageUtilization` of 50%, it means the HPA will attempt to maintain the average CPU utilization of the Pods at 50% of what is requested by each Pod. If a Pod has requested 200 millicores (0.2 CPU), the HPA aims to keep the utilization at around 100 millicores (0.1 CPU).

### Example Scenario:

Suppose you have a Deployment where each Pod requests 200 millicores of CPU (`cpu: "200m"`), and you set up an HPA targeting this Deployment with `averageUtilization: 50`. The HPA will scale up or down to try to maintain the CPU usage of each Pod at about 100 millicores on average.

### Why This Approach:

- **Predictable Scaling**: By tying the scaling behavior to the requested resources, Kubernetes provides a predictable way to scale applications based on the resources they are expected to use.
- **Avoids Overloading Nodes**: This approach helps ensure that Pods do not use more resources than the node can handle, as it's based on requests rather than limits or total capacity.
