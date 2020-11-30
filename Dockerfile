# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-Achecker"
LABEL REPO="https://github.com/itbiz-youtube/Achecker"

ENV PROJPATH=/go/src/github.com/itbiz-youtube/Achecker

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/itbiz-youtube/Achecker
WORKDIR /go/src/github.com/itbiz-youtube/Achecker

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/itbiz-youtube/Achecker"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/Achecker/bin

WORKDIR /opt/Achecker/bin

COPY --from=build-stage /go/src/github.com/itbiz-youtube/Achecker/bin/Achecker /opt/Achecker/bin/
RUN chmod +x /opt/Achecker/bin/Achecker

# Create appuser
RUN adduser -D -g '' Achecker
USER Achecker

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/Achecker/bin/Achecker"]
