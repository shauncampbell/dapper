FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/shauncampbell/dapper/
COPY . .

RUN go build -o /go/bin/dapper github.com/shauncampbell/dapper/cmd/dapper

FROM alpine:3.12

COPY --from=builder /go/bin/dapper /go/bin/dapper
LABEL maintainer="Shaun Campbell <docker@shaun.scot>"
COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

VOLUME /config.yaml
ENV LDAP_BASE "dc=home,dc=lab"
EXPOSE 389

ENTRYPOINT ["./docker-entrypoint.sh"]
