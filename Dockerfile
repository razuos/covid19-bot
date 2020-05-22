FROM golang:1.14 as build

WORKDIR /go/src/github.com/covid19-bot/covid19-bot
COPY . .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

FROM alpine
RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
  ca-certificates \
  && update-ca-certificates 2>/dev/null || true
WORKDIR /root/
COPY --from=build /go/src/github.com/covid19-bot/covid19-bot/bot .
CMD [ "./bot" ]