FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 go build -o /out/linkcuter ./cmd/linkcuter

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /out/linkcuter /app/linkcuter
COPY configs /app/configs

EXPOSE 8080

ENTRYPOINT ["/app/linkcuter"]
