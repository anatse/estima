# FROM alpine:latest
FROM golang:latest

MAINTAINER Anatoly Sementsov<anatolse@yandex.ru>

#RUN curl -O https://www.arangodb.com/repositories/arangodb31/Debian_8.0/Release.key
#RUN apt-key add - < Release.key

#RUN ls -l /etc/apt/sources.list.d
#RUN apt-get update && apt-get install apt-transport-https
#RUN mv /etc/apt/sources.list.d/google-cloud-sdk.list1 /etc/apt/sources.list.d/google-cloud-sdk.list

#RUN echo 'deb https://www.arangodb.com/repositories/arangodb31/Debian_8.0/ /' | tee /etc/apt/sources.list.d/arangodb.list
#RUN apt-get install apt-transport-https
#RUN apt-get install -y arangodb3
#RUN apt-get update && apt-get install arangodb3=3.1.22

ADD . /go/src/estima
ENV GOPATH /go/src/estima
WORKDIR /go/src/estima

RUN go build
ENTRYPOINT /go/src/estima/estima

EXPOSE 9080
