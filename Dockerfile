FROM golang:1.24.5-alpine AS build
WORKDIR /app
RUN apk add --no-cache ca-certificates build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/trikal ./cmd/trikal

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/trikal /trikal
COPY configs/config.example.yaml /configs/config.yaml
ENV TRIKAL_CONFIG=/configs/config.yaml
USER 65532:65532
ENTRYPOINT ["/trikal"]
