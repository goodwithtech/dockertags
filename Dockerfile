FROM alpine:3.11
RUN apk --no-cache add ca-certificates shadow
COPY dockertags /usr/local/bin/dockertags

# for use docker daemon via mounted /var/run/docker.sock
RUN addgroup -S docker && adduser -S -G docker dockertags
USER dockertags
ENTRYPOINT ["dockertags"]