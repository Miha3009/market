FROM golang:1.20 as builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /usr/src/app

COPY . ./
RUN go mod download

WORKDIR /usr/src/app/cmd
RUN go build -o /usr/local/bin/inventory -buildvcs=false

FROM alpine:latest
COPY --from=builder /usr/local/bin/inventory ./
CMD ["./inventory"]
EXPOSE 8080
