FROM alpine

RUN apk add --update ca-certificates

ADD ./ui /usr/bin/gitpods-api

EXPOSE 3020

ENTRYPOINT ["/usr/bin/gitpods-api"]
