FROM alpine

RUN apk add --update ca-certificates

ADD ./api /usr/bin/
RUN chmod +x /usr/bin/api

EXPOSE 3020

ENTRYPOINT ["/usr/bin/api"]
