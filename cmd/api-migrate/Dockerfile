FROM alpine

RUN apk add --update ca-certificates

ADD ./cmd/api-migrate/api-migrate /usr/bin/
ADD ./schema/postgres /migrations

ENV GITPODS_MIGRATIONS_PATH /migrations

ENTRYPOINT ["/usr/bin/api-migrate"]
