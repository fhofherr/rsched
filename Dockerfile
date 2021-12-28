FROM alpine:3.15.0
LABEL org.opencontainers.image.source https://github.com/fhofherr/rsched

ARG RSCHED_VERSION
ARG RESTIC_VERSION
ARG TARGETARCH

ENV RSCHED_RESTIC_BINARY /usr/local/bin/restic

COPY ./bin/rsched_${RSCHED_VERSION}_linux_${TARGETARCH} /usr/local/bin/rsched
COPY ./bin/restic_${RESTIC_VERSION}_linux_${TARGETARCH} ${RSCHED_RESTIC_BINARY}

ENTRYPOINT ["/usr/local/bin/rsched"]
