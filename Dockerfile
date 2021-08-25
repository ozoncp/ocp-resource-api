FROM ubuntu AS builder

RUN apt update -y
RUN apt upgrade -y

RUN apt install -y locales
RUN apt install -y sudo

RUN echo "LC_ALL=en_US.UTF-8" >> /etc/environment && \
    echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen && \
    echo "LANG=en_US.UTF-8" > /etc/locale.conf && \
    locale-gen en_US.UTF-8

RUN DEBIAN_FRONTEND="noninteractive" apt install -y golang
RUN apt install -y ca-certificates && sudo update-ca-certificates
RUN apt install -y make git vim protobuf-compiler

RUN useradd -m -G sudo developer
RUN echo 'developer:developer' | chpasswd
ENV GOPATH /home/developer/go
RUN mkdir $GOPATH
ENV PATH $PATH:/home/developer/go/bin
COPY . /home/developer/go/src/github.com/ozoncp/ocp-resource-api/
RUN chown -R "developer:developer" /home/developer

USER developer
WORKDIR /home/developer/go/src/github.com/ozoncp/ocp-resource-api

RUN make deps && make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /home/developer/go/src/github.com/ozoncp/ocp-resource-api/.bin/ocp-resource-api .
RUN chown root:root ocp-resource-api
EXPOSE 7070
EXPOSE 7072
CMD ["./ocp-resource-api"]