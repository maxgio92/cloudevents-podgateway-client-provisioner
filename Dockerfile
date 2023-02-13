FROM golang:1.19 AS builder

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go .
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -v -o provisioner

FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /workspace/provisioner . 
USER nonroot:nonroot

ENTRYPOINT ["/provisioner"]
