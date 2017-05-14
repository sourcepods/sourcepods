FROM alpine

RUN apk add --update ca-certificates

ADD ./api /usr/bin/gitpods-api
RUN chmod +x /usr/bin/gitpods-api

EXPOSE 3020

ENTRYPOINT ["/usr/bin/gitpods-api"]
