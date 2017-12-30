FROM alpine
MAINTAINER Ian Auld <im-auld@github>

RUN apk update && apk add ca-certificates
COPY ./bin/slackify /