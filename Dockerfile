FROM quay.io/giantswarm/alpine:3.9-giantswarm

ADD nic-exporter /
USER giantswarm

ENTRYPOINT ["/nic-exporter"]
