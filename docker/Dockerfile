FROM golang:1.23.2 AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/openlane-cloud -a .

FROM gcr.io/distroless/static:nonroot

# `nonroot` coming from distroless
USER 65532:65532

# Copy the binary that goreleaser built
COPY --from=builder /go/bin/openlane-cloud /bin/openlane-cloud

# Run the web service on container startup.
ENTRYPOINT [ "/bin/openlane-cloud" ]
CMD ["serve"]
