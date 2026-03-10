# HTTP/3 Configuration in K3s Traefik

## Status: CONFIGURED ✅

HTTP/3 has been successfully enabled on the cluster.

### Fix Applied (2026-03-08)

The `ERR_QUIC_PROTOCOL_ERROR` was caused by HTTP/3 not being properly configured in Traefik. The following changes were made:

1. Created `traefik-config.yaml` with HTTP/3 enabled
2. Applied HelmChart configuration with:
   - `http3.enabled: true` on websecure port
   - `service.single: false` to create separate TCP and UDP services
   - Both services share the same MetalLB IP (192.168.0.124)

### Fix Applied (2026-03-09) - Missing websecure Port

**Issue**: After recent changes, HTTPS services (including pihole.borkowskij.com/admin) became inaccessible.

**Root Cause**: The k3s Traefik addon HelmChart was missing the `websecure` port configuration in `ports`. This caused the Traefik pod to not expose port 8443, resulting in no service endpoint for port 443.

**Symptoms**:
- `kubectl describe svc traefik -n kube-system` showed empty Endpoints for port 443
- External connections to HTTPS services timed out
- Pod was listening on 8443 internally but port wasn't declared in container spec

**Solution**:
```bash
# Patch the HelmChart to include websecure port and image configuration
kubectl patch helmchart -n kube-system traefik --type=merge -p '{"spec":{"valuesContent":"image:\n  repository: rancher/mirrored-library-traefik\n  tag: \"3.6.9\"\napi:\n  dashboard: true\nports:\n  web:\n    port: 8000\n    expose:\n      default: true\n    exposedPort: 80\n    protocol: TCP\n    http:\n      redirections:\n        entryPoint:\n          to: websecure\n          scheme: https\n          permanent: true\n  websecure:\n    port: 8443\n    expose:\n      default: true\n    exposedPort: 443\n    protocol: TCP\n    http:\n      tls:\n        enabled: true\n    http3:\n      enabled: true\n      advertisedPort: 443\nproviders:\n  kubernetesIngress:\n    enabled: true\n    publishedService:\n      enabled: true\n  kubernetesCRD:\n    enabled: true\nservice:\n  enabled: true\n  single: false\n  annotations:\n    metallb.universe.tf/loadBalancerIPs: 192.168.0.124\n    metallb.universe.tf/allow-shared-ip: traefik-shared"}}'
```

**Note**: The k3s addon controller manages the Traefik HelmChart. When patching, ensure the `image` section is included with the correct tag format (`3.6.9` not `v3.6.9`).

### Result
- ✅ TCP Service: `traefik` on port 443 (HTTPS)
- ✅ UDP Service: `traefik-udp` on port 443 (HTTP/3/QUIC)
- ✅ Both services use the same LoadBalancer IP
- ✅ UDP port 443 is reachable and responding

### Working Applications (HTTP/3)

| Application | Protocol | Status |
|-------------|----------|--------|
| pihole.DOMAIN | h3 | ✅ Working |
| shoppinglist.DOMAIN | h3 | ✅ Working |
| filebrowser.DOMAIN | h3 | ✅ Working |

### Non-Working Applications (HTTP/2)

| Application | Protocol | Status |
|-------------|----------|--------|
| home.DOMAIN | h2 | ❌ Not upgrading |
| grafana.DOMAIN | h2 | ❌ Not upgrading |

## Configuration Applied

### Traefik HelmChart (`trafik-config.yaml`)

```yaml
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: traefik
  namespace: kube-system
spec:
  chart: traefik
  repo: https://helm.traefik.io/traefik
  targetNamespace: kube-system
  set:
    additionalArguments[0]: "--log.level=DEBUG"
  valuesContent: |-
    api:
      dashboard: true
    
    ports:
      web:
        port: 8000
        expose:
          default: true
        exposedPort: 80
        protocol: TCP
        http:
          redirections:
            entryPoint:
              to: websecure
              scheme: https
              permanent: true
      websecure:
        port: 8443
        expose:
          default: true
        exposedPort: 443
        protocol: TCP
        http:
          tls:
            enabled: true
        http3:
          enabled: true
          advertisedPort: 443
    
    providers:
      kubernetesIngress:
        enabled: true
        publishedService:
          enabled: true
      kubernetesCRD:
          enabled: true
    
    service:
      enabled: true
      single: false
      annotations:
        metallb.universe.tf/loadBalancerIPs: 192.168.0.124
        metallb.universe.tf/allow-shared-ip: "traefik-shared"
```

### Key Configuration Changes

1. **Use `ports` instead of `entryPoints`**: The Helm chart uses `ports` configuration, not `entryPoints`
2. **Enable HTTP/3 on websecure**: `http3.enabled: true` with `advertisedPort: 443`
3. **Dual Service Setup**: Chart automatically creates `traefik` (TCP) and `traefik-udp` (UDP) services
4. **Shared IP Annotation**: Both services must share the same IP using MetalLB annotations

### Service Status

```bash
$ kubectl get svc -n kube-system | grep traefik
traefik      LoadBalancer  192.168.0.124  80:30230/TCP,443:30284/TCP
traefik-udp  LoadBalancer  192.168.0.124  443:31623/UDP
```

### Traefik Arguments

```
--entryPoints.websecure.address=:8443/tcp
--entryPoints.websecure.http.tls=true
--entryPoints.websecure.http3
--entryPoints.websecure.http3.advertisedPort=443
```

### Response Headers

All applications correctly advertise HTTP/3:
```
alt-svc: h3=":443"; ma=2592000
```

## Verification

### Check Protocol in Browser
```javascript
// In browser console
const entries = performance.getEntriesByType('navigation');
console.log(entries[0].nextHopProtocol); // "h3" or "h2"
```

### Test UDP Connectivity
```bash
nc -zvu home.DOMAIN 443
# Connection succeeded
```

### Check Traefik HTTP/3 Status
```bash
kubectl logs -n kube-system deployment/traefik | grep -i "http3\|quic"
```

## Known Issues

### UDP Buffer Size Warning
Traefik logs show:
```
failed to sufficiently increase receive buffer size (was: 208 kiB, wanted: 7168 kiB, got: 416 kiB)
```
This is a warning but doesn't prevent HTTP/3 from working (as evidenced by working applications).

### Why Some Apps Don't Use HTTP/3

Possible reasons:
1. **Browser caching**: Chrome may have cached that HTTP/3 doesn't work for specific domains
2. **Certificate/QUIC handshake issues**: Some handshakes may be failing silently
3. **Application-specific issues**: Some backends may not properly support the HTTP/3 upgrade
4. **Connection pooling**: Existing HTTP/2 connections may persist

## Troubleshooting Steps Tried

- ✅ Fixed HelmChart configuration (changed `entryPoints` to `ports`)
- ✅ Enabled HTTP/3 on websecure entrypoint
- ✅ Configured dual services (TCP + UDP) with shared IP
- ✅ Verified UDP port 443 is accessible
- ✅ Confirmed `alt-svc` header is present
- ✅ Verified Traefik is listening on UDP 8443

## Notes

- HTTP/3 requires UDP port 443 to be open and accessible
- The `alt-svc` header advertises HTTP/3 availability: `h3=":443"; ma=2592000`
- Browsers handle the HTTP/3 upgrade automatically when they see the alt-svc header
- **HTTP/3 upgrade process**: Browser connects via HTTP/2 → sees alt-svc header → upgrades to HTTP/3 on next connection
- Some applications may take several page loads (2-3) before switching to HTTP/3
- Chrome/Edge may require clearing browser cache or using Incognito mode to test HTTP/3
- Once upgraded to HTTP/3, the browser will continue using it for subsequent requests

## Next Steps

To investigate why some applications remain on HTTP/2:

1. Check Chrome's net internals: `chrome://net-export/` (requires enabling logging)
2. Test with Firefox which has different HTTP/3 implementation
3. Verify TLS 1.3 is properly configured on all certificates
4. Check if any middleware or backend services are interfering with HTTP/3
5. Try accessing from a fresh browser profile to rule out caching issues
