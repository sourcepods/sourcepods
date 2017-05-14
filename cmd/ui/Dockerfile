FROM alpine

RUN apk add --update ca-certificates

ADD ./ui /usr/bin/
RUN chmod +x /usr/bin/ui

EXPOSE 3010

ENTRYPOINT ["/usr/bin/ui"]
