ARG PROJECT_VERSION=0.1.0

ARG GO_VERSION=1.12.5
ARG ALPINE_VERSION=3.9
ARG DOCKER_COMPOSE_VERSION=1.25.0-rc1

ARG BINARY_NAME=docker-cron
ARG BASE_DIR=docker-cron

ARG CONFIG_PATH=/etc/docker-config.env

# Compilation stage
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG DOCKER_REPO
ARG BINARY_NAME

# Install needed deps for building the binary
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git gcc libc-dev ca-certificates tzdata && update-ca-certificates

# Import project into the container
WORKDIR ${GOPATH}/src/github.com/${DOCKER_REPO}/
COPY . .

# Fetch dependencies; using go get.
RUN go get -d -v

# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/${BINARY_NAME}

# Destination stage
FROM docker/compose:${DOCKER_COMPOSE_VERSION}-alpine

ARG DOCKER_REPO
ARG DOCKER_TAG
ARG IMAGE_NAME
ARG COMMIT_MSG
ARG SOURCE_BRANCH
ARG SOURCE_COMMIT

ARG PROJECT_VERSION
ARG GO_VERSION
ARG ALPINE_VERSION
ARG DOCKER_COMPOSE_VERSION
ARG BINARY_NAME
ARG BASE_DIR
ARG CONFIG_PATH

LABEL maintainer="Corebreaker"
LABEL description="Have a cron scheduler for docker which run command in a docker container"
LABEL version="${PROJECT_VERSION}"

RUN echo "DOCKER_REPO='${DOCKER_REPO}'" >${CONFIG_PATH} \
    && echo "DOCKERFILE_PATH='${DOCKERFILE_PATH}'" >>${CONFIG_PATH} \
    && echo "DOCKER_TAG='${DOCKER_TAG}'" >>${CONFIG_PATH} \
    && echo "IMAGE_NAME='${IMAGE_NAME}'" >>${CONFIG_PATH} \
    && echo "COMMIT_MSG='${COMMIT_MSG}'" >>${CONFIG_PATH} \
    && echo "SOURCE_BRANCH='${SOURCE_BRANCH}'" >>${CONFIG_PATH} \
    && echo "SOURCE_COMMIT='${SOURCE_COMMIT}'" >>${CONFIG_PATH} \
    && echo "PROJECT_VERSION='${PROJECT_VERSION}'" >>${CONFIG_PATH} \
    && echo "ALPINE_VERSION='${ALPINE_VERSION}'" >>${CONFIG_PATH} \
    && echo "DOCKER_COMPOSE_VERSION='${DOCKER_COMPOSE_VERSION}'" >>${CONFIG_PATH} \
    && echo "BINARY_NAME='${BINARY_NAME}'" >>${CONFIG_PATH} \
    && echo "CRONBIN='/${BASE_DIR}/${BINARY_NAME}'" >>${CONFIG_PATH}

RUN mkdir /projects /${BASE_DIR}
COPY --from=builder /go/bin/${BINARY_NAME} /${BASE_DIR}/
COPY entry-point.sh /${BASE_DIR}/

WORKDIR /${BASE_DIR}/
VOLUME /projects

ENTRYPOINT ["./entry-point.sh"]
CMD ["start"]
