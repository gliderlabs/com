FROM golang:1.9
WORKDIR /go/src/github.com/gliderlabs/com
RUN go get github.com/gliderlabs/ssh
RUN go get github.com/spf13/viper
RUN go get github.com/thejerf/suture
COPY . .
COPY ./example/cmd/example/example.toml ./example.toml
RUN go build -o /example ./example/cmd/example/...
RUN go build -buildmode=plugin -o /ssh.so ./example/plugins/ssh/...
ENTRYPOINT ["/example"]
