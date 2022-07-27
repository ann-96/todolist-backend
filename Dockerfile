FROM golang:1.18.1 as builder

RUN mkdir /build
COPY . /build/
WORKDIR /build

RUN go test ./...
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -o service-binary .


FROM alpine:latest
COPY --from=builder /build/service-binary .
RUN mkdir -p ./migrations
COPY ./migrations/. ./migrations/

ENTRYPOINT [ "./service-binary" ]