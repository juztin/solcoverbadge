FROM golang as golang
COPY . /src
WORKDIR /src
RUN go get -d ./...
RUN CGO_ENABLED=0 go build solcoverbadge.go

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
COPY --from=golang /src/solcoverbadge /bin/solcoverbadge
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /src
ENTRYPOINT ["/bin/solcoverbadge"]
