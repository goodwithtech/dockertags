FROM alpine:3.13 as builder
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY dockertags /dockertags
ENTRYPOINT ["/dockertags"]