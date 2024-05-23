FROM golang:1.22-alpine@sha256:b8ded51bad03238f67994d0a6b88680609b392db04312f60c23358cc878d4902 AS builder

ARG KTMPL_VERSION

WORKDIR /src
ADD . .
RUN go build -ldflags "-X \"main.ktmplVersion=${KTMPL_VERSION}\" -X \"main.compileDate=$(date)\"" -o ./bin/ktmpl .

# ---

FROM debian:latest@sha256:fac2c0fd33e88dfd3bc88a872cfb78dcb167e74af6162d31724df69e482f886c as kubectl_downloader

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

# ---

FROM debian:latest@sha256:fac2c0fd33e88dfd3bc88a872cfb78dcb167e74af6162d31724df69e482f886c

RUN apt-get update && apt-get install -y curl xz-utils && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/bin/ktmpl /usr/bin/ktmpl
COPY --from=kubectl_downloader /kubectl /usr/bin/kubectl
