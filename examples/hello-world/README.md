HTTP hello world

```
$ curl localhost:8080

Hello, world!
Version: 1.0.0
```

### tekton

Create registry secret

```shell
# edit the config.json
# XXXX get from echo -n USER:ACCESS_TOKEN | base64
kubectl create secret generic kaniko-secret --from-file=config.json
```

Install task

```
tkn hub install task git-clone
tkn hub install task kaniko
```

Create Pipeline and PipelineRun

```
kubectl apply -f pipeline.yaml
kubectl create -f pipeline-run.yaml
```
