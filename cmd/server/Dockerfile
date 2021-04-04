FROM alpine:3.7 as alpine

RUN apk add -U --no-cache ca-certificates && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \
    -D \
    -g "" \
    -h "/nonexistent" \
    -H \
    -s "/sbin/nologin" \
    -u "${UID}" \
    "${USER}"

ADD https://raw.githubusercontent.com/NABSA/gbfs/master/systems.csv /systems.csv
RUN chmod +r /systems.csv

FROM scratch

COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /systems.csv /systems.csv

# Use an unprivileged user
USER appuser:appuser

EXPOSE 8080

ADD server /server

ENTRYPOINT ["/server"]
