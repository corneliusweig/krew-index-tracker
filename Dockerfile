FROM golang:alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY app ./app/
COPY pkg ./pkg/

RUN go build -tags netgo -ldflags "-s -w" -o krew-index-tracker ./app/http

FROM alpine:3.10
LABEL maintainer=cornelius.weig@gmail.com
EXPOSE 8080

RUN adduser -S service-user

WORKDIR /app
RUN chown service-user /app && \
    apk add --no-cache git

USER service-user

RUN mkdir index && \
  git clone --depth 1 --branch master https://github.com/kubernetes-sigs/krew-index.git index

ENTRYPOINT ["/app/krew-index-tracker"]

COPY --from=builder /app/krew-index-tracker ./
