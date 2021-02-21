FROM alpine:3.13
RUN apk --no-cache add ca-certificates shadow
COPY dockertags /usr/local/bin/dockertags
RUN addgroup -S docker && adduser -S -G docker dockertags
USER dockertags
ENTRYPOINT ["dockertags"]