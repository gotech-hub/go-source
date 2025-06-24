# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# [START cloudrun_helloworld_dockerfile]
# [START run_helloworld_dockerfile]

# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.21.5-alpine3.19  as builder

RUN apk add --no-cache git ca-certificates build-base

ARG GITLAB_JOB_USER
ARG GITLAB_JOB_TOKEN

RUN go env -w GOPRIVATE=gitlab.com/*

RUN git config --global url.https://${GITLAB_JOB_USER}:${GITLAB_JOB_TOKEN}@gitlab.com.insteadOf https://gitlab.com
RUN echo "machine gitlab.com login $GITLAB_JOB_USER password $GITLAB_JOB_TOKEN" > ~/.netrc

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./

RUN go mod download

# Copy local code to the container image.
COPY . .

RUN go mod tidy

# Build the binary.
RUN CGO_ENABLED=1 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -tags musl -v -o runner main.go

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3.19.0

USER root

RUN set -x && apk update && apk add --no-cache \
    ca-certificates curl  && \
    rm -rf /var/lib/apt/lists/*

# Create user
ARG USER=zdeploy_docker
ARG GROUP=zdeploy_docker
ARG GID=2000
ARG UID=2000
ENV HOME_DIR=/home/$USER
RUN addgroup -g ${GID} ${GROUP} && adduser -u ${UID} -G ${GROUP} -s /bin/sh -D ${USER}

# Tell docker that all future commands should run as the user
USER $USER
WORKDIR $HOME_DIR

# Copy the binary to the production image from the builder stage.
COPY --from=builder --chown=$USER:$GROUP /app/runner .
COPY --from=builder --chown=$USER:$GROUP /app/run.sh .

RUN chmod +x run.sh
RUN chown $USER:$GROUP ./* -R

# Run the web service on container startup.
EXPOSE 3000
ENTRYPOINT ["/bin/sh", "-c", "./run.sh"]

