#Dockerfile for fibonacci project.
FROM golang:1.16.4-alpine3.13 as builder
RUN echo "Don't forget to first disable any local instances of Postgres with: sudo systemctl stop postgresql@12-main.service"
ENV GO111MODULE=on
COPY go.mod go.sum migration.sh /go/src/github.com/dgnabasik/fibonacci/
WORKDIR /go/src/github.com/dgnabasik/fibonacci
RUN go mod download
COPY . /go/src/github.com/dgnabasik/fibonacci
RUN go get github.com/dgnabasik/fibonacci
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/fibonacci github.com/dgnabasik/fibonacci

FROM alpine
#RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/dgnabasik/fibonacci/build/fibonacci /usr/bin/fibonacci
EXPOSE 8080 8080
#RUN ["/usr/bin/fibonacci/migration.sh"]
ENTRYPOINT ["/usr/bin/fibonacci"]