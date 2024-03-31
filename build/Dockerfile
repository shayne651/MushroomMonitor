FROM ubuntu:latest

WORKDIR /src

RUN apt-get update &&\
	apt-get install wget -y &&\
	apt-get install gcc -y

RUN wget https://go.dev/dl/go1.22.1.linux-arm64.tar.gz &&\
	tar -xzvf go1.22.1.linux-arm64.tar.gz -C /

ENV PATH="${PATH}:/go/bin"
# What the container should run when it is started.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -o /bin/server ./cmd/mushroommonitor/main.go

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

RUN touch /src/mush.db

RUN chown -R appuser:appuser /src

USER appuser

COPY ./db/migrations /db/migrations

# Expose the port that the application listens on.
EXPOSE 7891

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]