FROM golang:1.22 AS builder

WORKDIR /usr/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/api/ cmd/api/
COPY internal/ internal/
COPY generated/ generated/
COPY sql/ sql/
RUN CGO_ENABLED=0 go build -o api cmd/api/main.go

FROM scratch

COPY --from=builder /usr/src/app/api /api

ENTRYPOINT ["/api"]
