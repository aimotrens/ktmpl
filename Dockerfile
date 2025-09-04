FROM golang:1.25-alpine@sha256:2ad042dec672e85d9e631feb0d2d72db86fd2a4e0cf8daaf2c19771a26df1062 AS builder

ARG KTMPL_VERSION

WORKDIR /src
ADD . .
RUN go build -ldflags "-X \"main.ktmplVersion=${KTMPL_VERSION}\" -X \"main.compileDate=$(date)\"" -o ./bin/ktmpl ./cmd/ktmpl

# ---

FROM debian:latest@sha256:6d87375016340817ac2391e670971725a9981cfc24e221c47734681ed0f6c0f5 AS kubectl_downloader

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

# ---

FROM debian:latest@sha256:6d87375016340817ac2391e670971725a9981cfc24e221c47734681ed0f6c0f5

RUN apt-get update && apt-get install -y curl xz-utils && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/bin/ktmpl /usr/bin/ktmpl
COPY --from=kubectl_downloader /kubectl /usr/bin/kubectl
