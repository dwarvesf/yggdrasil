FROM alpine:3.7
RUN apk add --no-cach ca-certificates
WORKDIR /
COPY server /server
EXPOSE 2700
ENTRYPOINT ["/server"]