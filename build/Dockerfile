FROM golang:1.22 as builder

WORKDIR /usr/app
COPY .. .

RUN go get -d -v ./...
RUN go install -v ./...
RUN  CGO_ENABLED=0 go build -o application ./cmd/.

FROM debian:stretch-slim
COPY --from=builder /usr/app/application /usr/app/application
EXPOSE 8080

CMD ["/usr/app/application"]
