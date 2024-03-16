FROM golang:1.22-alpine@sha256:0466223b8544fb7d4ff04748acc4d75a608234bf4e79563bff208d2060c0dd79 AS builder

ARG KTMPL_VERSION

WORKDIR /src
ADD . .
RUN go build -ldflags "-X \"main.ktmplVersion=${KTMPL_VERSION}\" -X \"main.compileDate=$(date)\"" -o ./bin/ktmpl .

# ---

FROM debian:latest@sha256:e97ee92bf1e11a2de654e9f3da827d8dce32b54e0490ac83bfc65c8706568116 as kubectl_downloader

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

# ---

FROM debian:latest@sha256:e97ee92bf1e11a2de654e9f3da827d8dce32b54e0490ac83bfc65c8706568116

RUN apt-get update && apt-get install -y curl xz-utils && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/bin/ktmpl /usr/bin/ktmpl
COPY --from=kubectl_downloader /kubectl /usr/bin/kubectl
