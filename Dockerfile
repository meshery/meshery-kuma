FROM golang:1.13.7 as bd

WORKDIR /github.com/layer5io/meshery-kuma
<<<<<<< HEAD
ARG CONFIG_PROVIDER="viper"
=======
ARG CONFIG_PROVIDER="local"
>>>>>>> 3ced31b8d82ed2a98f924ad60f7b4c5040e61f1f
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o meshery-kuma main.go

FROM alpine
RUN apk --update add ca-certificates curl
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
<<<<<<< HEAD
	mkdir -p /home/scripts && \
	mkdir -p /root/.kube/ && \
	mkdir -p /root/.kuma/
=======
	mkdir -p /home/scripts/kuma && \
	mkdir -p /root/.kube/
>>>>>>> 3ced31b8d82ed2a98f924ad60f7b4c5040e61f1f

# Install kubectl
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl" && \
	chmod +x ./kubectl && \
	mv ./kubectl /usr/local/bin/kubectl

COPY --from=bd /github.com/layer5io/meshery-kuma/meshery-kuma /home/
<<<<<<< HEAD
COPY --from=bd /github.com/layer5io/meshery-kuma/scripts/** /home/scripts
=======
COPY --from=bd /github.com/layer5io/meshery-kuma/scripts/** /home/scripts/kuma
>>>>>>> 3ced31b8d82ed2a98f924ad60f7b4c5040e61f1f
WORKDIR /home
CMD ./meshery-kuma
