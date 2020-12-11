FROM golang as builder
WORKDIR /cloud-run-example
COPY go.mod .
# COPY go.sum .
# RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/hello-world .

FROM alpine
RUN addgroup -S example-group && adduser -S -D example -G example-group
USER example
WORKDIR /home/example
COPY --from=builder /bin/hello-world ./
EXPOSE 8080
ENTRYPOINT ["./hello-world"]
