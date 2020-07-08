# Compile stage
FROM golang:alpine AS build-env
RUN apk add --no-cache git

ARG COMMIT="n/a"
ARG TAG="n/a"
ARG CONFIG="master"

#ENV GOPROXY=direct
ENV GO111MODULE=on
ENV GOPRIVATE=*.inn4science.com,gitlab.com

WORKDIR /service
ADD . .
COPY ./env/${CONFIG}.config.yaml /config.yaml
RUN go mod download && go build -ldflags "-X main.Build=$COMMIT -X main.Tag=$TAG" -o /app .

# Final stage
FROM alpine:3.7

# Port 8080 belongs to our application
EXPOSE 8080

# Allow delve to run on Alpine based containers.
RUN apk add --no-cache ca-certificates bash

WORKDIR /

COPY --from=build-env /app /
COPY --from=build-env /config.yaml /

# Run delve
CMD ["/app", "serve"]
