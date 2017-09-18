FROM alpine

RUN apk add --update ca-certificates git

ADD ./cmd/storage/storage /usr/bin/

EXPOSE 3030

ENTRYPOINT ["/usr/bin/storage"]
