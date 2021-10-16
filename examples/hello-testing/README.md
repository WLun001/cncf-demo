Example of integration testing that rely on database with docker compose

## Run app

```angular2html
MONGO_URI=mongodb://admin:apple123@localhost:27017 go run .
```

## Run test locally

```shell
docker-compose up
```

```
# open another terminal
go test -v .

2021/10/16 16:15:48 db client connected
2021/10/16 16:15:48 db client ping
=== RUN   TestAddItem
--- PASS: TestAddItem (0.00s)
=== RUN   TestGetItems
--- PASS: TestGetItems (0.00s)
PASS
ok      hello-testing   0.022s
```

```
# to shutdown and delete volume
docker-compose down -v
```
