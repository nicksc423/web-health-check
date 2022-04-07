FROM golang:alpine AS builder

RUN apk update && apk add --no-cache make bash gcc musl-dev libc-dev
RUN adduser -D -g '' appuser

WORKDIR /app
COPY . .

RUN make build

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/bin/packator-health-check /app/bin/packator-health-check

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/bin/web-health-check"]
