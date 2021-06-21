API server with cache ability. Embeddable, in-memory db using [BadgerDB](https://dgraph.io/docs/badger/)
Run the API

```
# initial run might take time
# to start BadgerDB
go run .
```

Available APIs

- `/` no cache
- `/cache` has header Cache-Control: public, max-age=5
- `/server-cache` embeddable memory db cache, TTL 5 seconds
- `/server-and-cdn-cache`embeddable memory db cache, TTL to 30 seconds and has header Cache-Control: public, max-age=5
- `error/:statusCode` example error page based on status code
- `server-1/example` example route
- `server-2/example`example route
- `/db/all` query db records
