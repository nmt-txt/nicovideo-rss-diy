FROM golang:1.25-alpine AS builder
WORKDIR /app
#COPY go.mod go.sum ./ #現在何もインストールしていない
#RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o niconico-rss-diy main.go


FROM gcr.io/distroless/static-debian12 AS runner
WORKDIR /app
COPY --from=builder /app/niconico-rss-diy .

EXPOSE 8080
USER nonroot
ENTRYPOINT ["./niconico-rss-diy"]
CMD ["/config/"] 