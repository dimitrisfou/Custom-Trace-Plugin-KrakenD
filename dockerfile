FROM golang:1.19.5-alpine3.17 as builder

COPY ./krakend.json /app/krakend/krakend.json
COPY ./go.mod /app/krakend/plugins/go.mod
COPY ./*.go /app/krakend/plugins/

WORKDIR /app/krakend/plugins

# Initialize the Go module and download dependencies
RUN go mod tidy
# RUN go mod download

# for alpine image
# https://stackoverflow.com/questions/43580131/exec-gcc-executable-file-not-found-in-path-when-trying-go-build
RUN apk add build-base

# ChatGPT for linux
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o pluginServer.so pluginServer.go

FROM devopsfaith/krakend:2.2.0

COPY --from=builder /app/krakend/krakend.json /etc/krakend/krakend.json
COPY --from=builder /app/krakend/plugins/*.so /etc/krakend/plugins/