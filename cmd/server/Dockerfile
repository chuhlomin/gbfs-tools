FROM alpine:3.7 as alpine

RUN apk add -U --no-cache ca-certificates && update-ca-certificates

FROM scratch

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD server /server

ENTRYPOINT ["/server"]
