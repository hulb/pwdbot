FROM golang:latest AS build
WORKDIR /app
COPY . /app
RUN go build -o pwdbot -i main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build pwdbot /app/

ENTRYPOINT ["/app/pwdbot"]