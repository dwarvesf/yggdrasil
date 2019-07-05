FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY worker /worker
ENTRYPOINT ["/worker"]
