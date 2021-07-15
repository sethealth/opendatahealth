FROM golang:1.16.3 as build-env

ARG COMMIT_SHA=""

WORKDIR /opendata
COPY go.mod go.sum /opendata/
RUN go mod download && go mod verify

COPY . /opendata/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.Revision=$COMMIT_SHA" -o app

FROM gcr.io/distroless/base
COPY --from=build-env /opendata/app /
CMD ["/app"]
