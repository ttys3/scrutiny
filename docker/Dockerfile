########
FROM golang:1.15.3-buster as backendbuild

WORKDIR /go/src/github.com/analogj/scrutiny

COPY . /go/src/github.com/analogj/scrutiny

RUN go mod vendor && \
    go build -tags timetzdata -ldflags '-w -extldflags "-static"' -o scrutiny webapp/backend/cmd/scrutiny/scrutiny.go && \
    go build -tags timetzdata -ldflags '-w -extldflags "-static"' -o scrutiny-collector-selftest collector/cmd/collector-selftest/collector-selftest.go && \
    go build -tags timetzdata -ldflags '-w -extldflags "-static"' -o scrutiny-collector-metrics collector/cmd/collector-metrics/collector-metrics.go

########
FROM node:lts-slim as frontendbuild

#reduce logging, disable angular-cli analytics for ci environment
ENV NPM_CONFIG_LOGLEVEL=warn NG_CLI_ANALYTICS=false

WORKDIR /scrutiny/src
COPY webapp/frontend /scrutiny/src

RUN npm install -g @angular/cli@9.1.4 && \
    mkdir -p /scrutiny/dist && \
    npm install && \
    ng build --output-path=/scrutiny/dist --deploy-url="/web/" --base-href="/web/" --prod


########
FROM ubuntu:focal as runtime

EXPOSE 8080
WORKDIR /scrutiny
ENV PATH="/scrutiny/bin:${PATH}"

ADD https://github.com/dshearer/jobber/releases/download/v1.4.4/jobber_1.4.4-1_amd64.deb /tmp/

RUN apt install /tmp/jobber_1.4.4-1_amd64.deb

RUN apt-get update && apt-get install -y smartmontools ca-certificates curl tzdata && update-ca-certificates

ADD https://github.com/just-containers/s6-overlay/releases/download/v2.1.0.0/s6-overlay-amd64.tar.gz /tmp/
# prevent tar override /bin/sh

# mkdir -p /tmp/s6 && tar -xzf /tmp/s6-overlay-amd64.tar.gz -C /tmp/s6 && cp -ra /tmp/s6/* / && /bin/sh --version
# cp: cannot overwrite non-directory '/bin' with directory '/tmp/s6/bin'

# because ubuntu focal /bin is a symlink
# ❯ docker run -ti --rm ubuntu:focal ls -lhp /
# total 48K
# lrwxrwxrwx   1 root root    7 Sep 25 01:20 bin -> usr/bin
# drwxr-xr-x   2 root root 4.0K Apr 15  2020 boot/
# drwxr-xr-x   5 root root  360 Oct 17 07:05 dev/
# drwxr-xr-x   1 root root 4.0K Oct 17 07:05 etc/
# drwxr-xr-x   2 root root 4.0K Apr 15  2020 home/
# lrwxrwxrwx   1 root root    7 Sep 25 01:20 lib -> usr/lib
# lrwxrwxrwx   1 root root    9 Sep 25 01:20 lib32 -> usr/lib32
# lrwxrwxrwx   1 root root    9 Sep 25 01:20 lib64 -> usr/lib64
# lrwxrwxrwx   1 root root   10 Sep 25 01:20 libx32 -> usr/libx32
# drwxr-xr-x   2 root root 4.0K Sep 25 01:20 media/
# drwxr-xr-x   2 root root 4.0K Sep 25 01:20 mnt/
# drwxr-xr-x   2 root root 4.0K Sep 25 01:20 opt/
# dr-xr-xr-x 772 root root    0 Oct 17 07:05 proc/
# drwx------   2 root root 4.0K Sep 25 01:23 root/
# drwxr-xr-x   1 root root 4.0K Sep 25 22:34 run/
# lrwxrwxrwx   1 root root    8 Sep 25 01:20 sbin -> usr/sbin
# drwxr-xr-x   2 root root 4.0K Sep 25 01:20 srv/
# dr-xr-xr-x  13 root root    0 Oct 17 07:05 sys/
# drwxrwxrwt   2 root root 4.0K Sep 25 01:23 tmp/
# drwxr-xr-x   1 root root 4.0K Sep 25 01:20 usr/
# drwxr-xr-x   1 root root 4.0K Sep 25 01:23 var/

# fucking dash has no --version param

RUN mkdir -p /tmp/s6 && tar -xzf /tmp/s6-overlay-amd64.tar.gz -C /tmp/s6 && mv /tmp/s6/bin/* /tmp/s6/usr/bin/ && rmdir /tmp/s6/bin && cp -ra /tmp/s6/* / && ls -lhp /bin/sh

COPY /rootfs /

COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny /scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-selftest /scrutiny/bin/
COPY --from=backendbuild /go/src/github.com/analogj/scrutiny/scrutiny-collector-metrics /scrutiny/bin/
COPY --from=frontendbuild /scrutiny/dist /scrutiny/web

RUN chmod +x /scrutiny/bin/scrutiny && \
    chmod +x /scrutiny/bin/scrutiny-collector-selftest && \
    chmod +x /scrutiny/bin/scrutiny-collector-metrics && \
    mkdir -p /scrutiny/web && \
    mkdir -p /scrutiny/config && \
    mkdir -p /scrutiny/jobber

CMD ["/init"]
