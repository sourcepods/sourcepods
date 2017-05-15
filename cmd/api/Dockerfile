FROM alpine

RUN apk add --update ca-certificates

ADD ./api /usr/bin/

EXPOSE 3020

ENTRYPOINT ["/usr/bin/api"]
