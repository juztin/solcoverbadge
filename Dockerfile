FROM golang as builder
COPY . /src
WORKDIR /src
RUN go get -d ./... && go build solcoverbadge.go

FROM scratch
COPY --from=builder /src/solcoverbadge /bin/solcoverbadge
WORKDIR /src
ENTRYPOINT ["/bin/solcoverbadge"]
