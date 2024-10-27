FROM golang:alpine AS builderGo
RUN apk --no-cache -U upgrade && \
    apk --no-cache add --upgrade make build-base ca-certificates && \
    wget -qO- https://downloads.sqlc.dev/sqlc_1.27.0_linux_amd64.tar.gz | tar xvz -C /bin
# Set working directory
WORKDIR /go/src/github.com/JDinABox/yapa
# Prefetch downloads
COPY go.* ./
RUN go mod download
#COPY files
COPY ./Makefile ./Makefile
COPY ./cmd/ ./cmd/
COPY ./internal ./internal

#COPY *.go ./
RUN --mount=type=cache,target=/root/.cache/go-build make build

FROM alpine:latest

RUN apk --no-cache -U upgrade \
    && apk --no-cache add --upgrade ca-certificates \
    && wget -O /bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64 \
    && chmod +x /bin/dumb-init

COPY --from=builderGo /go/src/github.com/JDinABox/yapa/cmd/yapa/yapa.out /bin/yapa.out

ENTRYPOINT ["/bin/dumb-init", "--"]
CMD ["/bin/yapa.out"]