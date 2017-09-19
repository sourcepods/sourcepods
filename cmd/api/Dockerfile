FROM alpine

RUN apk add --update ca-certificates

COPY ./cmd/api/api /usr/bin/
COPY ./schema/postgres /migrations

EXPOSE 3020

ENTRYPOINT ["/usr/bin/api"]
