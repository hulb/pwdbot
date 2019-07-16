FROM golang:latest
WORKDIR /app
COPY . /app
RUN mkdir -p /app/data
RUN go build pwdbot && rm -rf .git && rm go.mod && rm main.go && rm go.sum && rm README.md && rm .gitignore

ENTRYPOINT ["/app/pwdbot"]