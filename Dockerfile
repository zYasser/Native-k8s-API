# syntax=docker/dockerfile:1.7

# Shared base: download dependencies and copy the source
FROM golang:1.26-alpine AS build-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .


FROM build-base AS compile-list-pods

RUN --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    CGO_ENABLED=0 go build \
    -o /out/list-pods \
    ./list-pods-with-k8s


FROM build-base AS compile-crud-pods

RUN --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    CGO_ENABLED=0 go build \
    -o /out/crud-pods \
    ./crud-pods


FROM build-base AS compile-watch-pods

RUN --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    CGO_ENABLED=0 go build \
    -o /out/watch-pods \
    ./watch-pods


FROM build-base AS compile-watch-pods-and-modify-it

RUN --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    CGO_ENABLED=0 go build \
    -o /out/watch-pods-and-modify-it \
    ./watch-pods-and-modify-it


FROM gcr.io/distroless/static-debian12:nonroot AS list-pods

COPY --from=compile-list-pods /out/list-pods /list-pods

ENTRYPOINT ["/list-pods"]


FROM gcr.io/distroless/static-debian12:nonroot AS crud-pods

COPY --from=compile-crud-pods /out/crud-pods /crud-pods

ENTRYPOINT ["/crud-pods"]


FROM gcr.io/distroless/static-debian12:nonroot AS watch-pods

COPY --from=compile-watch-pods /out/watch-pods /watch-pods

ENTRYPOINT ["/watch-pods"]


FROM gcr.io/distroless/static-debian12:nonroot AS watch-pods-and-modify-it

COPY --from=compile-watch-pods-and-modify-it \
    /out/watch-pods-and-modify-it \
    /watch-pods-and-modify-it

ENTRYPOINT ["/watch-pods-and-modify-it"]