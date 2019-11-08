FROM golang:1.13.3 as bd
RUN adduser --disabled-login appuser
WORKDIR /github.com/layer5io/meshery-kuma
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -a -o /meshery-kuma .
RUN find . -name "*.go" -type f -delete; mv kuma /

FROM alpine
RUN apk --update add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=bd /meshery-kuma /app/
COPY --from=bd /kuma /app/kuma
COPY --from=bd /etc/passwd /etc/passwd
USER appuser
WORKDIR /app
CMD ./meshery-kuma
