FROM golang:1.20-alpine as build

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY publisher.go ./

RUN go build publisher.go

FROM alpine as runtime

COPY --from=build /src/publisher /app/publisher

CMD [ "/app/publisher" ]