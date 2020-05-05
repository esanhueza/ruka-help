FROM node:lts-alpine AS frontend-builder
ARG VUE_APP_BASE_PATH
ENV VUE_APP_BASE_PATH $VUE_APP_BASE_PATH
WORKDIR /client
ADD . /client
RUN npm --prefix ui install
RUN npm run --prefix ui build

FROM golang:alpine AS builder
WORKDIR /websocket-echo-client
COPY --from=frontend-builder /websocket-echo-client/ /websocket-echo-client/
RUN apk update && apk add git && apk add ca-certificates
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o vue-websocket-echo