FROM golang:1.22-alpine@sha256:32c85006b1edf29c097514e0c81a33334aa1450685a885c10657ec756dbb7703 AS builder

ARG KTMPL_VERSION

WORKDIR /src
ADD . .
RUN go build -ldflags "-X \"main.ktmplVersion=${KTMPL_VERSION}\" -X \"main.compileDate=$(date)\"" -o ./bin/ktmpl .

# ---

FROM debian:latest@sha256:a92ed51e0996d8e9de041ca05ce623d2c491444df6a535a566dabd5cb8336946 as kubectl_downloader

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

# ---

FROM debian:latest@sha256:a92ed51e0996d8e9de041ca05ce623d2c491444df6a535a566dabd5cb8336946

RUN apt-get update && apt-get install -y curl xz-utils && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/bin/ktmpl /usr/bin/ktmpl
COPY --from=kubectl_downloader /kubectl /usr/bin/kubectl
