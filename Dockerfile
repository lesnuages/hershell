FROM golang:alpine
RUN apk add --update make bash git openssl && go get github.com/lesnuages/hershell
WORKDIR /go/src/github.com/lesnuages/hershell/

# To retreive built binary, use docker run -v "/tmp:/go/bin" hershell:latest $1 $2

# todo Add docker entrypoint to grab params and bach compile, copy generated certs to bin folder (hershellCerts/)
# Make sure NOT TO REGEN certs
ENTRYPOINT [ "bash" ]