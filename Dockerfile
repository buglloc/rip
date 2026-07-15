FROM golang:1.26 AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o /go/bin/rip

FROM debian:bookworm-slim

COPY --from=build /go/bin/rip /usr/sbin/rip

ENTRYPOINT ["/usr/sbin/rip"]
