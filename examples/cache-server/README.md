API server with cache ability. Embeddable, in-memory db using [BadgerDB](https://dgraph.io/docs/badger/)
Run the API

```
go run .
```

Available APIs

- `/` no cache
- `/cache` has header Cache-Control: public, max-age=5
- `/server-cache` embeddable memory db cache, TTL 5 seconds
- `/server-and-cdn-cache`embeddable memory db cache, TTL to 30 seconds and has header Cache-Control: public, max-age=5
- `/db/all` query db records
