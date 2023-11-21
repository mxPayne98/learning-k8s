## Kubernetes Readiness Probes

Readiness probes in Kubernetes are used to determine if a container is ready to start accepting traffic. They are critical for managing the lifecycle of Pods and ensuring that services direct traffic only to healthy instances.

### How Readiness Probes Work

- **Purpose**: A readiness probe is used to signal to the Kubernetes control plane that your application is ready to receive traffic.
- **Checking**: Kubernetes regularly performs the specified check (HTTP GET, TCP socket, or exec command) to determine the readiness of the container.
- **Traffic Routing**: If a container fails its readiness probe, it is removed from service endpoints, meaning it won't receive traffic from services.

### Types of Readiness Probes

1. **HTTP GET**: Kubernetes sends an HTTP GET request to the container. A response code within 200-399 indicates success.
2. **TCP Socket**: Kubernetes tries to establish a TCP connection to a specified port of the container. A successful connection indicates readiness.
3. **Exec Command**: Executes a specified command inside the container. A return code of 0 indicates success.

### Example YAML Configuration

Here’s an example of how to define readiness probes in your Pod's YAML file:

**readiness-probe.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: readiness-demo
spec:
  containers:
  - name: demo-container
    image: demo-image
    readinessProbe:
      httpGet:
        path: /healthz
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5
```

- `readinessProbe`: Configures the probe.
- `httpGet`: Specifies an HTTP GET probe.
  - `path`: The path to access on the container.
  - `port`: The port on which the server is listening.
- `initialDelaySeconds`: Number of seconds after the container starts before the probe is initiated.
- `periodSeconds`: How often (in seconds) to perform the probe.

### Best Practices and Considerations

1. **Choose the Right Probe**: Select a probe type that best reflects your container's readiness for traffic.
2. **Avoid Heavy Load Operations**: The readiness probe endpoint should be lightweight and not put significant load on the application.
3. **Initial Delay Configuration**: Set an appropriate initial delay to give your application enough time to start up before the first probe.
4. **Handling Probe Failures**: Ensure that your application correctly handles and recovers from a failed readiness state.
5. **Probe Endpoint Security**: If using an HTTP probe, secure the endpoint if it exposes sensitive information.

### Considerations
- **Accurate Readiness Checks**: Ensure that your readiness probe accurately reflects the readiness state of your application.
- **Avoid False Positives**: The readiness logic should be robust enough to avoid false positives, where the service is marked as ready but isn't actually able to handle requests properly.
- **Resource Usage**: Keep the probe lightweight to avoid excessive resource usage.
- **Probe Endpoint Security**: Secure the probe endpoint if it exposes sensitive information.

### Probe configurations

### 1. HTTP GET Readiness Probe

An HTTP GET probe is used when your container has an HTTP server running. Kubernetes sends an HTTP GET request to the specified path and port.

**Example YAML Configuration**:

```yaml
readinessProbe:
  httpGet:
    path: /ready  # Endpoint in your Node.js application
    port: 8080    # Port on which your app is listening
  initialDelaySeconds: 10
  periodSeconds: 5
  failureThreshold: 3
  successThreshold: 1
```

- `httpGet`: Defines the HTTP GET request.
- `initialDelaySeconds`: Number of seconds after the container has started before probes are initiated.
- `periodSeconds`: How often to perform the probe.
- `failureThreshold`: Number of times to retry the probe before giving up.
- `successThreshold`: Minimum consecutive successes for the probe to be considered successful after having failed.

### 2. TCP Socket Readiness Probe

A TCP Socket probe checks if your container is listening on a specified TCP port.

**Example YAML Configuration**:

```yaml
readinessProbe:
  tcpSocket:
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
  failureThreshold: 3
```

- `tcpSocket`: Kubernetes tries to establish a TCP connection on the specified port.
- Other fields are similar to the HTTP GET probe.

### 3. Exec Command Readiness Probe

The exec command probe executes a specified command inside the container. The probe is considered successful if the command exits with a status code of 0.

**Example YAML Configuration**:

```yaml
readinessProbe:
  exec:
    command:
    - cat
    - /tmp/healthy
  initialDelaySeconds: 10
  periodSeconds: 5
  failureThreshold: 3
```

- `exec`: Specifies the command to execute.
- `command`: The command and its arguments.

### Integrating `/ready` Endpoint in Node.js App

In your Node.js application, you can create a `/ready` endpoint to respond to readiness probes.

**Example Node.js Code**:

```javascript
const express = require('express');
const app = express();

app.get('/ready', (req, res) => {
  // Implement your logic to determine if the app is ready
  const isReady = /* logic to check readiness */;
  if (isReady) {
    res.status(200).send('Ready');
  } else {
    res.status(500).send('Not Ready');
  }
});

app.listen(8080, () => {
  console.log('Server started on port 8080');
});
```

- This endpoint checks if the application is ready to handle traffic and responds with an appropriate HTTP status code.
