FROM alpine:latest

RUN apk --no-cache add curl

COPY ./marge-backend /usr/local/bin

WORKDIR /usr/local/bin

CMD ["marge-backend"]
