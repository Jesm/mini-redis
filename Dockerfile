FROM golang:1.12 AS build

WORKDIR /go/src/miniredis/

COPY . /go/src/miniredis/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/miniredis miniredis

FROM scratch
COPY --from=build /go/bin/miniredis /go/bin/miniredis
ENTRYPOINT ["/go/bin/miniredis"]
