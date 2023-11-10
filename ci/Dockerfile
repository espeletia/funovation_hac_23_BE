FROM golang:1.21-alpine3.17

COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffmpeg /usr/bin/
COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffprobe /usr/bin/

ADD ./build/ /app/
ADD ./configurations /app/configurations/
ADD ./migrations /app/migrations


WORKDIR /app

ENTRYPOINT ./start_server
