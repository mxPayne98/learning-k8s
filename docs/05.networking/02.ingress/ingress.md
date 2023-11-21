## Kubernetes Ingress

Ingress in Kubernetes is an API object that manages external access to services in a cluster, typically HTTP. Ingress can provide load balancing, SSL termination, and name-based virtual hosting.

### How Ingress Works

- **Function**: An Ingress is a set of rules that allow inbound connections to reach the cluster services.
- **Ingress Controller**: To use Ingress, you need an Ingress controller, like NGINX Ingress Controller, Traefik, or HAProxy. The controller reads the Ingress Resource information and processes the data accordingly.
- **Routing Traffic**: Ingress can be configured to give services externally-reachable URLs, load balance traffic, terminate SSL, offer name-based virtual hosting, etc.

### Basic Ingress Features

1. **Path-based Routing**: Routes traffic to different services based on the URL path.
2. **Host-based Routing**: Routes traffic based on the requested hostname.
3. **TLS/SSL Termination**: Handles encrypted traffic (https termination), offloading the work from the services.

### Example YAML Configuration for Ingress

Here's an example of an Ingress that routes traffic based on the URL path:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
spec:
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /service1
        pathType: Prefix
        backend:
          service:
            name: service1
            port:
              number: 80
      - path: /service2
        pathType: Prefix
        backend:
          service:
            name: service2
            port:
              number: 80
```

- This Ingress routes traffic to `service1` when the path starts with `/service1` and to `service2` for paths starting with `/service2`.
- The `host` field allows routing based on the domain name (optional).

### Advanced Ingress Features

1. **TLS Configuration**:
   - Secure your Ingress by adding TLS configuration. You can specify a secret that contains a TLS private key and certificate.
   ```yaml
   spec:
     tls:
     - hosts:
       - myapp.example.com
       secretName: myapp-secret-tls
   ```

2. **Rewrite Targets**:
   - Some Ingress controllers support rewriting the URL path.

3. **Annotations**:
   - You can use annotations to customize the behavior of the Ingress controller.

4. **Authentication and Authorization**:
   - Certain Ingress controllers allow you to integrate authentication and authorization directly.

5. **Rate Limiting and CORS**:
   - Advanced features like rate limiting and Cross-Origin Resource Sharing (CORS) can be configured on some Ingress controllers.

### Considerations

- **Ingress Controller Installation**: An Ingress controller is not automatically started with a cluster; you must deploy it yourself.
- **DNS Configuration**: Ingress does not provide DNS services. You need to configure DNS records manually or via external DNS providers.
- **Resource Isolation**: Since Ingress operates at the edge of the network, it’s crucial to ensure its security and isolation from other cluster resources.

### Commands for Managing Ingress in Kubernetes

Managing Ingress resources in Kubernetes typically involves creating and applying YAML configuration files. However, you can use `kubectl` for basic operations like viewing and editing.

#### Common `kubectl` Commands for Ingress

1. **List All Ingress Resources in the Current Namespace**:
   ```bash
   kubectl get ingress
   ```

2. **List Ingress Resources Across All Namespaces**:
   ```bash
   kubectl get ingress --all-namespaces
   ```

3. **Describe a Specific Ingress Resource**:
   ```bash
   kubectl describe ingress <ingress-name>
   ```
   - Provides detailed information about the specified Ingress.

4. **Delete an Ingress Resource**:
   ```bash
   kubectl delete ingress <ingress-name>
   ```

5. **Edit an Ingress Resource**:
   ```bash
   kubectl edit ingress <ingress-name>
   ```
   - Opens the Ingress resource in an editor for modification.

### Multiple Ingress Resources and Nuances

Having multiple Ingress resources in a Kubernetes cluster is common, especially in environments with numerous services that require external access.

#### Considerations with Multiple Ingress Resources

1. **Overlap and Conflicts**: Ensure that Ingress resources do not have conflicting rules or overlap in a way that can cause routing conflicts.
2. **Controller Specificity**: Different Ingress controllers might handle overlapping rules differently. Be aware of how your chosen Ingress controller resolves conflicts.
3. **Resource Separation**: Using multiple Ingress resources can help in separating concerns, such as differentiating between production and development environments.
4. **Host and Path-Based Routing**: You can have multiple Ingress resources with different hostnames and path-based rules, directing traffic to various services.
5. **Performance and Scalability**: Keep an eye on the performance and scalability aspects, as each Ingress controller might handle load differently.

### Commonly Used Ingress Annotations

Annotations in Ingress allow you to specify additional configurations and behaviors that are specific to the Ingress controller you are using.

#### Examples of Common Annotations

1. **Rewrite Target** (common in NGINX Ingress Controller):
   ```yaml
   nginx.ingress.kubernetes.io/rewrite-target: /
   ```
   - Modifies the path of the request before it is sent to the service.

2. **SSL Redirect**:
   ```yaml
   nginx.ingress.kubernetes.io/ssl-redirect: "true"
   ```
   - Redirects HTTP traffic to HTTPS.

3. **Force SSL** (using Let's Encrypt with cert-manager):
   ```yaml
   cert-manager.io/cluster-issuer: "letsencrypt-prod"
   ```
   - Automatically configures SSL certificates.

4. **Rate Limiting**:
   ```yaml
   nginx.ingress.kubernetes.io/limit-rps: "5"
   ```
   - Restricts the number of requests per second to the service.

5. **Affinity and Session Stickiness**:
   ```yaml
   nginx.ingress.kubernetes.io/affinity: "cookie"
   ```
   - Enables session affinity using cookies.

6. **Custom Error Pages**:
   ```yaml
   nginx.ingress.kubernetes.io/custom-http-errors: "404,500,502"
   ```
   - Defines custom error pages for specific HTTP errors.

7. **CORS Configuration**:
   ```yaml
   nginx.ingress.kubernetes.io/enable-cors: "true"
   ```
   - Enables Cross-Origin Resource Sharing (CORS) policies.

#### Examples of Advanced Configurations

1. **Customizing Timeouts**:
   - Configure timeouts for connections.
   ```yaml
   nginx.ingress.kubernetes.io/proxy-connect-timeout: "10"
   nginx.ingress.kubernetes.io/proxy-read-timeout: "20"
   ```

2. **Configuring Maximum Body Size**:
   - Set the maximum size of the client request body.
   ```yaml
   nginx.ingress.kubernetes.io/proxy-body-size: "8m"
   ```

3. **Websocket Support**:
   - Enable support for websockets.
   ```yaml
   nginx.ingress.kubernetes.io/websocket-services: "websocket-service-name"
   ```

4. **Setting Custom Headers**:
   - Add custom headers to the response.
   ```yaml
   nginx.ingress.kubernetes.io/add-headers: "custom-header-name"
   ```

5. **External Authentication**:
   - Use an external service for authentication.
   ```yaml
   nginx.ingress.kubernetes.io/auth-url: "https://auth-service.example.com/auth"
   nginx.ingress.kubernetes.io/auth-signin: "https://auth-service.example.com/start"
   ```

6. **Rewriting the Request Path**:
   - Change the path of the request before forwarding it to the service.
   ```yaml
   nginx.ingress.kubernetes.io/rewrite-target: /$1
   ```

7. **Using Snippets for Custom Configuration**:
   - Inject custom NGINX configuration.
   ```yaml
   nginx.ingress.kubernetes.io/configuration-snippet: |
     more_set_headers "X-Custom-Header: value";
   ```

8. **Source IP Preservation**:
   - Preserve the real source IP of the client.
   ```yaml
   nginx.ingress.kubernetes.io/preserve-client-ip: "true"
   ```
   