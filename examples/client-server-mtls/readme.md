# example of mtls (client and server)

## Generate certs
```bash
.certs/gen-certs.sh
```
## Env
envs and default values
### Server

```bash
SERVER_ADDR=:8443
CERT_PATH=certs/server.crt
KEY_PATH=certs/server.key
CA_CERTS_PATH=certs/ca.crt
```

### Client
```bash
CERT_PATH=certs/client.crt
KEY_PATH=certs/client.key
CA_CERTS_PATH=certs/ca.crt
LOOP=false
SERVER_URL=https://localhost:8443
```
## Example deployment

### Secrets
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: vault-pki-client
  namespace: default
type: Opaque
data:
  ca.crt: xxx
  client.crt: xxx
  client.key: xxx
```
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: vault-pki-server
  namespace: default
type: Opaque
data:
  ca.crt: xxx
  server.crt: xxx
  server.key: xxx
```

### Server
```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: mtls-server
  labels:
    name: mtls-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mtls-server
  template:
    metadata:
      labels:
        app: mtls-server
    spec:
      containers:
        - name: mtls
          image: ghcr.io/wlun001/mtls:latest
          ports:
            - containerPort: 8443
          volumeMounts:
            - mountPath: /home/example/certs
              name: vault-pki
      volumes:
        - name: vault-pki
          secret:
            secretName: vault-pki-server
            defaultMode: 420
---
apiVersion: v1
kind: Service
metadata:
  name: mtls-server
spec:
  ports:
    - name: https
      port: 8443
  selector:
    app: mtls-server
```

### Client
```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: mtls-client
  labels:
    name: mtls-client
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mtls-client
  template:
    metadata:
      labels:
        app: mtls-client
    spec:
      containers:
        - name: mtls-client
          image: ghcr.io/wlun001/mtls:latest
          command:
            - ./client
          env:
            - name: LOOP
              value: "true"
            - name: SERVER_URL
              value: "https://mtls.default.svc.cluster.local:8443"
          volumeMounts:
            - mountPath: /home/example/certs
              name: vault-pki
      volumes:
        - name: vault-pki
          secret:
            secretName: vault-pki-client
            defaultMode: 420
```