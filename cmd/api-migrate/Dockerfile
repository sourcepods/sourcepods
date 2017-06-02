FROM alpine

RUN apk add --update ca-certificates

ADD ./api-migrate /usr/bin/

ENTRYPOINT ["/usr/bin/api-migrate"]
