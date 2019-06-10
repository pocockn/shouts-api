FROM vidsyhq/go-base:latest
LABEL maintainer="Nick Pocock"

ARG VERSION
LABEL version=$VERSION

ADD shouts-api /
ADD config /config

ENTRYPOINT ["/shouts-api"]