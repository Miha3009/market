FROM golang:1.20 as builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /usr/src/app

COPY . ./
RUN go mod download

WORKDIR /usr/src/app/cmd
RUN go build -o /usr/local/bin/notifier -buildvcs=false

FROM alpine:latest
COPY --from=builder /usr/local/bin/notifier ./
CMD ["./notifier"]
