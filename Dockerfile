FROM golang:1.22.0-bullseye AS builder
WORKDIR /builder

# Import the codebase
COPY . .

# Install dependencies
RUN go mod vendor

# Compile the source
RUN CGO_ENABLED=0 go build -mod vendor -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/build ./cmd/main.go


FROM scratch AS runner
# FROM golang:1.22.0-bullseye AS runner
WORKDIR /

# Copy binary from builder
COPY --from=builder /builder/bin/build ./bin/build

# Run the binary
CMD ["/bin/build"]
