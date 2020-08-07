FROM golang:1.14-alpine3.11 as build-img
ARG CONFIG_PROVIDER="viper"
RUN apk update && apk add --no-cache git libc-dev gcc pkgconf && mkdir /home/adaptor
COPY ${PWD} /go/src/github.com/layer5io/meshery-kuma/
WORKDIR /go/src/github.com/layer5io/meshery-kuma/
RUN go mod vendor && go build -ldflags="-w -s -X main.configProvider=$CONFIG_PROVIDER" -a -o /home/adaptor/kuma

FROM alpine
RUN apk --update add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
	mkdir ${HOME}/.kuma/ && \
	mkdir /home/adaptor/scripts/
COPY --from=bd /home/adaptor /home/
COPY --from=bd /go/src/github.com/layer5io/meshery-kuma/scripts/** /home/adaptor/scripts/
WORKDIR /home/adaptor
CMD ./kuma
