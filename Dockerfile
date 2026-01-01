
FROM golang:1.24

WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o ts-derp-verifier ./cmd/ts-derp-verifier

FROM alpine:3.14
COPY --from=0 /build/ts-derp-verifier /bin/

ENTRYPOINT [ "/bin/ts-derp-verifier" ]