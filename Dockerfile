
ARG IMAGE_BASE_REGISTRY

ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.3

FROM ${IMAGE_BASE_REGISTRY}golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build

LABEL maintainer=arash.idelchi

USER root

SHELL ["/bin/ash", "-o", "pipefail", "-c"]

RUN apk add --no-cache \
    tzdata

WORKDIR /work

ARG GOMODCACHE=/home/${USER}/.cache/.go-mod
ARG GOCACHE=/home/${USER}/.cache/.go

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    go mod download

COPY . .
ARG TCISD_VERSION="unofficial & built by unknown"
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    CGO_ENABLED=0 go install -ldflags="-s -w -X 'main.version=${TCISD_VERSION}'" ./...

ARG USER=user
RUN addgroup -S -g 1001 ${USER} && \
    adduser -S -u 1001 -G ${USER} -h /home/${USER} -s /bin/ash ${USER}

ENV TZ=Europe/Zurich

FROM scratch AS standalone

LABEL maintainer=arash.idelchi

ARG USER=user

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /go/bin/tcisd /tcisd

USER ${USER}
WORKDIR /home/${USER}

ENTRYPOINT ["/tcisd"]
CMD [""]

ENV TZ=Europe/Zurich

FROM ${IMAGE_BASE_REGISTRY}alpine:${ALPINE_VERSION}

LABEL maintainer=arash.idelchi

USER root

RUN apk add --no-cache \
    tzdata

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/bin/tcisd /usr/local/bin/tcisd

USER ${USER}
WORKDIR /home/${USER}

ENTRYPOINT [""]
CMD ["/bin/ash"]

ENV TZ=Europe/Zurich
