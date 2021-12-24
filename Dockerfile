FROM scratch
LABEL org.opencontainers.image.source https://github.com/fhofherr/rsched

ARG RSCHED_VERSION
ARG RESTIC_VERSION
ARG TARGETARCH

ENV RSCHED_RESTIC_BINARY /bin/restic

COPY ./bin/rsched_${RSCHED_VERSION}_linux_${TARGETARCH} /bin/rsched
COPY ./bin/restic_${RESTIC_VERSION}_linux_${TARGETARCH} ${RSCHED_RESTIC_BINARY}

ENTRYPOINT ["/bin/rsched"]
