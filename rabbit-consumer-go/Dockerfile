FROM golang:1.20-alpine as build

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY consumer.go ./

RUN go build consumer.go

FROM alpine as runtime

COPY --from=build /src/consumer /app/consumer

CMD [ "/app/consumer" ]