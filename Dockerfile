FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY ./bin/video_upload .

CMD ["./video_upload"]