FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -o /out/list-pods \
    ./list-pods-with-k8s

RUN CGO_ENABLED=0 go build \
    -o /out/crud-pods \
    ./crud-pods

RUN CGO_ENABLED=0 go build \
    -o /out/watch-pods \
    ./watch-pods


FROM gcr.io/distroless/static-debian12:nonroot AS list-pods

COPY --from=build /out/list-pods /list-pods
ENTRYPOINT ["/list-pods"]


FROM gcr.io/distroless/static-debian12:nonroot AS crud-pods

COPY --from=build /out/crud-pods /crud-pods
ENTRYPOINT ["/crud-pods"]


FROM gcr.io/distroless/static-debian12:nonroot AS watch-pods

COPY --from=build /out/watch-pods /watch-pods
ENTRYPOINT ["/watch-pods"]