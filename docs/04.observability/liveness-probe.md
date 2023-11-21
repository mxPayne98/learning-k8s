## Kubernetes Liveness Probes

Liveness probes in Kubernetes are used to determine if a container is running and healthy. If a liveness probe fails, the Kubernetes kubelet will kill the container, and the container will be subjected to its restart policy.

### How Liveness Probes Work

- **Purpose**: Liveness probes are used to check the health of a container. If the check fails, Kubernetes will restart the container based on the Pod's restart policy.
- **Checking Mechanisms**: Similar to readiness probes, liveness probes can be configured to perform checks using HTTP GET, TCP Socket, or an Exec command.

### Types of Liveness Probes

1. **HTTP GET**: Kubernetes sends an HTTP GET request to the container. If the call fails, the container is restarted.
2. **TCP Socket**: Kubernetes attempts to establish a TCP connection. If it fails, the container is restarted.
3. **Exec Command**: Executes a command inside the container. If the command exits with a non-zero status, the container is restarted.

### Example YAML Configurations

#### 1. HTTP GET Liveness Probe

```yaml
livenessProbe:
  httpGet:
    path: /health  # Health-check endpoint
    port: 8080     # Port of your application
  initialDelaySeconds: 15
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

#### 2. TCP Socket Liveness Probe

```yaml
livenessProbe:
  tcpSocket:
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 10
  failureThreshold: 3
```

#### 3. Exec Command Liveness Probe

```yaml
livenessProbe:
  exec:
    command:
    - cat
    - /tmp/alive
  initialDelaySeconds: 15
  periodSeconds: 10
  failureThreshold: 3
```

### Parameters Explained

- `initialDelaySeconds`: Time to wait before the first probe is initiated.
- `periodSeconds`: How often the probe is performed.
- `timeoutSeconds`: Number of seconds after which the probe times out (applicable for HTTP and Exec probes).
- `failureThreshold`: Number of times the probe is retried before giving up.

### Integrating Liveness Probe in Node.js Application

Add a `/health` endpoint in your Node.js application that returns a successful response if the application is healthy.

**Example Node.js Code**:

```javascript
const express = require('express');
const app = express();

app.get('/health', (req, res) => {
  // Add logic to check the health of your application
  const isHealthy = /* logic to check health */;
  if (isHealthy) {
    res.status(200).send('OK');
  } else {
    res.status(500).send('Error');
  }
});

app.listen(8080, () => {
  console.log('Server started on port 8080');
});
```

### Best Practices

- **Correct Endpoint**: Ensure the liveness probe points to an endpoint that accurately reflects the health of the application.
- **Avoid False Positives**: Avoid situations where the liveness probe is passing, but the application is unable to handle requests.
- **Resource Usage**: The liveness probe endpoint should be lightweight and should not significantly impact application performance.
- **Avoid Aggressive Probes**: Setting probes with too short a period can lead to unnecessary restarts, which could destabilize your application.
