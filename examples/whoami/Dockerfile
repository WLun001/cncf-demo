FROM golang AS builder
WORKDIR /whoami
COPY go.mod .
# COPY go.sum .
# RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/whoami .

FROM alpine
RUN addgroup -S example-group && adduser -S -D example -G example-group
USER example
WORKDIR /home/example
COPY --from=builder /bin/whoami ./
EXPOSE 3000
ENTRYPOINT ["./whoami"]
