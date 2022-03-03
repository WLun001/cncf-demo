### tekton example with auth (git & registry)

Create registry secret

```shell
# edit the config.json
# XXXX get from echo -n USER:ACCESS_TOKEN | base64
kubectl create secret generic kaniko-secret --from-file=config.json
```

Create secret,
and [add new ssh to github](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account)

```shell
# create git ssh
# check command in secret.yaml and replace it
# add to github ssh

kubectl apply -f secret.yaml
```

Install task

```shell
tkn hub install task git-clone
tkn hub install task kaniko
```

Create Pipeline and PipelineRun

```shell
kubectl apply -f pipeline.yaml
kubectl create -f pipeline-run.yaml
```
