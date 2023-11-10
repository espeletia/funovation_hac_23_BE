FROM golang:1.21-alpine3.17

ADD ./build/ /app/
ADD ./configurations /app/configurations/
ADD ./migrations /app/migrations


WORKDIR /app

ENTRYPOINT ./start_server
