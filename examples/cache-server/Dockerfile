FROM golang AS builder
WORKDIR /cache-server
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/cache-server .

FROM alpine
RUN addgroup -S example-group && adduser -S -D example -G example-group
USER example
WORKDIR /home/example
COPY --from=builder /bin/cache-server ./
COPY --from=builder /cache-server/views ./views
EXPOSE 3000
ENTRYPOINT ["./cache-server"]
