#Dockerfile
FROM golang:1.16.4-alpine3.13 as builder
COPY go.mod go.sum /go/src/github.com/dgnabasik/fibonacci/
WORKDIR /go/src/github.com/dgnabasik/fibonacci
RUN go mod download
COPY . /go/src/github.com/dgnabasik/fibonacci
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/fibonacci github.com/dgnabasik/fibonacci

FROM alpine
#RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/dgnabasik/fibonacci/build/fibonacci /usr/bin/fibonacci
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/fibonacci"]