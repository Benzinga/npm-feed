FROM golang:alpine
COPY . /go/src/github.com/Benzinga/npm-feed
RUN go build -o /app github.com/Benzinga/npm-feed

FROM scratch
COPY --from=0 /app /app
ENTRYPOINT ["/app"]
