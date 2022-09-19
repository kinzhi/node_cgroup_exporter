ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

ARG ARCH="amd64"
ARG OS="linux"
COPY main /bin/node_cgroup_exporter

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_cgroup_exporter" ]