FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -o /out/list-pods \
    ./list-pods-with-k8s

# RUN CGO_ENABLED=0 go build \
#     -o /out/create-pod \
#     ./create-pod

# RUN CGO_ENABLED=0 go build \
#     -o /out/delete-pod \
#     ./delete-pod


FROM gcr.io/distroless/static-debian12:nonroot AS list-pods

COPY --from=build /out/list-pods /list-pods
ENTRYPOINT ["/list-pods"]


# FROM gcr.io/distroless/static-debian12:nonroot AS create-pod

# COPY --from=build /out/create-pod /create-pod
# ENTRYPOINT ["/create-pod"]


# FROM gcr.io/distroless/static-debian12:nonroot AS delete-pod

# COPY --from=build /out/delete-pod /delete-pod
# ENTRYPOINT ["/delete-pod"]