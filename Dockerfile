FROM golang:1.19.5-alpine3.17 AS golang

ARG BUILD_REF
ARG LOC=/builds/go/src/github.com/betorvs
RUN apk add --no-cache git
RUN mkdir -p $LOC
ENV GOPATH /go
COPY main.go go.mod go.sum $LOC/
COPY internal $LOC/internal
WORKDIR $LOC
ENV CGO_ENABLED 0
RUN go test ./... && go build -ldflags "-X main.build=${BUILD_REF}"

FROM alpine:3.17
ARG LOC=/builds/go/src/github.com/betorvs
WORKDIR /
VOLUME /tmp
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN mkdir -p /app
RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -G app -S -D -h /app app && \
    chmod 755 /app
COPY --from=golang $LOC/plhello /app

EXPOSE 9090
RUN chmod +x /app/plhello
WORKDIR /app    
USER app
CMD ["/app/plhello"]

LABEL org.opencontainers.image.title="plhello" \
      org.opencontainers.image.authors="Roberto Scudeller <beto.rvs@gmail.com>" \
      org.opencontainers.image.source="https://github.com/betorvs/plhello" \
      org.opencontainers.image.revision="${BUILD_REF}"