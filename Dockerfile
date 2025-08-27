FROM --platform=$BUILDPLATFORM golang:1.22 AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o server .

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/server .
ENTRYPOINT ["./server"]

