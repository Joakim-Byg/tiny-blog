FROM --platform=${BUILDPLATFORM} golang:1-alpine3.19 as build

ARG BUILDPLATFORM=linux/arm64
ARG TARGETARCH=arm64

RUN addgroup --gid 10001 -S appgroup && adduser --uid 10001 -S appuser -G appgroup -s "/sbin/nologin" -H

WORKDIR /build
COPY main.go .
COPY go.mod .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o ./bin/

FROM --platform=${TARGETPLATFORM:-linux/arm64} scratch as app

# Copy ssl, user and group data
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
#COPY --from=build /etc/shadow /etc/shadow
#COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

# Copy app binaries and static web-content
COPY --chown=appuser:appgroup --from=build /build/bin /app
COPY --chown=appuser:appgroup static /app/static
COPY --chown=appuser:appgroup favicon.ico /app/

WORKDIR "/app"
USER appuser:appgroup
ENTRYPOINT ["./tiny-blog"]
