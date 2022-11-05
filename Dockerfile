FROM golang:1.18.4-alpine AS build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod vendor

RUN go build -o main .

FROM alpine

WORKDIR /

COPY --from=build /app/main /usr/local/bin/main

COPY --from=build /app/docs .

ENTRYPOINT [ "main" ]