FROM golang AS builder
WORKDIR /client-server-mtls
COPY go.mod .
# COPY go.sum .
# RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server server.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/client client.go

FROM alpine
RUN addgroup -S example-group && adduser -S -D example -G example-group
USER example
WORKDIR /home/example
COPY --from=builder /bin/server ./
COPY --from=builder /bin/client ./
EXPOSE 8443
ENTRYPOINT ["./server"]
