FROM golang:1.23-alpine@sha256:2b975675437f0758b6cde49c7f59e9e30a811ed27783101b980c78f74b2345ce AS builder

ARG KTMPL_VERSION

WORKDIR /src
ADD . .
RUN go build -ldflags "-X \"main.ktmplVersion=${KTMPL_VERSION}\" -X \"main.compileDate=$(date)\"" -o ./bin/ktmpl .

# ---

FROM debian:latest@sha256:4abf773f2a570e6873259c4e3ba16de6c6268fb571fd46ec80be7c67822823b3 as kubectl_downloader

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

# ---

FROM debian:latest@sha256:4abf773f2a570e6873259c4e3ba16de6c6268fb571fd46ec80be7c67822823b3

RUN apt-get update && apt-get install -y curl xz-utils && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/bin/ktmpl /usr/bin/ktmpl
COPY --from=kubectl_downloader /kubectl /usr/bin/kubectl
