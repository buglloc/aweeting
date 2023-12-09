FROM golang:1.21.5 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/aweeting ./cmd/aweeting

FROM debian:bookworm-slim

COPY --from=build /go/bin/aweeting /usr/sbin/aweeting

ENTRYPOINT ["/usr/sbin/aweeting"]
