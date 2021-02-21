FROM alpine:3.13
COPY dockertags /dockertags
ENTRYPOINT ["/dockertags"]