FROM --platform=$BUILDPLATFORM golang:alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o app

FROM --platform=$TARGETPLATFORM alpine

COPY --from=builder /src/app ./

ENTRYPOINT ["./app"]
