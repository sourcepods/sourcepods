FROM alpine

RUN apk add --update ca-certificates

ADD ./ui /usr/bin/

EXPOSE 3010

ENTRYPOINT ["/usr/bin/ui"]
