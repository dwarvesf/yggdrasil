FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY server /server
EXPOSE 1492
ENTRYPOINT ["/server"]