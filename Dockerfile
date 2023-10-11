FROM golang:1.20-alpine as builder
RUN apk add build-base

WORKDIR /app
ADD . /app
#COPY .env /app

RUN go build -o /iate

FROM alpine:3.18.4

COPY --from=builder /iate /
COPY --from=builder /app/.env /

EXPOSE 80

CMD [ "/iate" ]