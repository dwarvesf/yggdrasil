FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY worker /worker
EXPOSE 3000
ENTRYPOINT ["/worker"]
