FROM golang:latest AS build
WORKDIR /app
COPY . /app
RUN go build -o /app/pwdbot -i main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/pwdbot /app/pwdbot/pwdbot

ENTRYPOINT ["/app/pwdbot"]