### tekton example without auth (git & registry)

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
