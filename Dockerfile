FROM alpine

COPY  ./fankserver-cli /fankserver-cli

RUN adduser -D -u 997 cli && \
    apk add --no-cache curl
USER cli

# This container will be executable
ENTRYPOINT ["/fankserver-cli"]
