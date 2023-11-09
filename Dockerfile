FROM node:alpine AS build-node
ARG version=0.0.0
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build-css

FROM golang:alpine AS build-go
ARG version=0.0.0
WORKDIR /app
COPY . .
RUN go build -o /go/bin ./...

FROM alpine:latest
ARG version
ENV APP_VERSION=$version
WORKDIR /app
COPY --from=build-go /go/bin/blog .
COPY --from=build-go /app/cmd/blog/dist ./dist
COPY --from=build-node /app/cmd/blog/dist/assets/output.css ./dist/assets/output.css
ENTRYPOINT ["./blog"]