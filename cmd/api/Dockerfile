FROM alpine

RUN apk add --update ca-certificates

ADD ./cmd/api/api /usr/bin/
ADD ./schema/postgres /migrations

EXPOSE 3020

ENTRYPOINT ["/usr/bin/api"]
