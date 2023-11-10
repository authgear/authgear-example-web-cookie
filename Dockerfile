FROM --platform=$BUILDPLATFORM golang:1.21.4-bullseye AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETARCH
ARG TARGETOS
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go build -tags authgearlite -o /serve ./cmd/serve

FROM gcr.io/distroless/static-debian11
COPY templates/* /templates/
COPY --from=build /serve /
CMD ["/serve"]
