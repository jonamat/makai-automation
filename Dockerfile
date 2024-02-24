FROM golang:1.17.0-bullseye AS builder
WORKDIR /build

# Import the codebase
COPY . .

# Create server binary, fetch & convert icons, generate init KML
RUN go build -mod vendor -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/makai-automation ./cmd/main.go


FROM scratch AS runner
WORKDIR /app

# Server binary from builder
COPY --from=builder /build/bin/makai-automation ./bin/makai-automation

# Run the server
ENTRYPOINT ["/app/bin/makai-automation"]
