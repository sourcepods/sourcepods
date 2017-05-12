FROM alpine

RUN apk add --update ca-certificates

ADD ./ui /usr/bin/gitpods-ui

EXPOSE 3010

ENTRYPOINT ["/usr/bin/gitpods-ui"]
