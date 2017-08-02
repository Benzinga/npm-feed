FROM golang:alpine
COPY . /go/src/github.com/Benzinga/npm-feed
ENV CGO_ENABLED=0
RUN go build -o /app github.com/Benzinga/npm-feed

FROM scratch
COPY --from=0 /app /app
ENTRYPOINT ["/app"]
