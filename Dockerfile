FROM golang:latest AS build
WORKDIR /app
COPY . /app
RUN mkdir -p /app/data
RUN env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -tags netgo pwdbot && rm -rf .git && rm go.mod && rm main.go && rm go.sum && rm README.md && rm .gitignore

FROM alpine:3.5
RUN apk update && apk add ca-certificates
COPY --from=build /app/pwdbot /app/pwdbot

ENTRYPOINT ["/app/pwdbot"]