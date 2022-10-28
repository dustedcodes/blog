FROM golang:alpine AS build-env
RUN apk add git
ARG version=0.0.0
WORKDIR /app
COPY . .
RUN go build -o /go/bin ./...

FROM alpine:latest
ARG version
ENV APP_VERSION=$version
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
WORKDIR /app
COPY --from=build-env /go/bin/blog .
COPY --from=build-env /app/cmd/blog/dist ./dist
ENTRYPOINT ["./blog"]