Example of integration testing that rely on database with docker compose

## Run app

```bash
MONGO_URI=mongodb://admin:apple123@localhost:27017 go run .
```

## Run test locally

```bash
docker-compose up
```

```bash
# open another terminal
MONGO_URI=mongodb://admin:apple123@localhost:27017 go test -v .

2021/10/16 16:15:48 db client connected
2021/10/16 16:15:48 db client ping
=== RUN   TestAddItem
--- PASS: TestAddItem (0.00s)
=== RUN   TestGetItems
--- PASS: TestGetItems (0.00s)
PASS
ok      hello-testing   0.022s
```

```bash
# to shutdown and delete volume
docker-compose down -v
```

## Run test on Cloud Build

```bash
gcloud builds submit .
```

## Run test on tekton

Install task

```
tkn hub install task git-clone
```

Create Pipeline and PipelineRun

```
kubectl apply -f pipeline.yaml
kubectl create -f pipeline-run.yaml
```
