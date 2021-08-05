FROM golang:1.16.3 as build-env

WORKDIR /opendata
COPY sentinel/go.mod sentinel/go.sum /opendata/
RUN go mod download && go mod verify

COPY ./sentinel /opendata/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app

FROM gcr.io/distroless/base
LABEL maintainer="Manu <manu@set.health>"

COPY --from=build-env /opendata/app /opendata_sentinel
CMD ["/opendata_sentinel"]
