FROM golang:1.25-alpine@sha256:98e6cffc31ccc44c7c15d83df1d69891efee8115a5bb7ede2bf30a38af3e3c92 AS builder

ARG KTMPL_VERSION

WORKDIR /src
ADD . .
RUN go build -ldflags "-X \"main.ktmplVersion=${KTMPL_VERSION}\" -X \"main.compileDate=$(date)\"" -o ./bin/ktmpl ./cmd/ktmpl

# ---

FROM debian:latest@sha256:2c91e484d93f0830a7e05a2b9d92a7b102be7cab562198b984a84fdbc7806d91 AS kubectl_downloader

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

# ---

FROM debian:latest@sha256:2c91e484d93f0830a7e05a2b9d92a7b102be7cab562198b984a84fdbc7806d91

RUN apt-get update && apt-get install -y curl xz-utils && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/bin/ktmpl /usr/bin/ktmpl
COPY --from=kubectl_downloader /kubectl /usr/bin/kubectl
