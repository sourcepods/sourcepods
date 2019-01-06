FROM alpine

RUN apk add --update ca-certificates git

COPY ./dev/storage /usr/bin/

EXPOSE 3030

ENTRYPOINT ["/usr/bin/storage"]
