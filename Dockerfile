FROM scratch
COPY dockertags /dockertags
ENTRYPOINT ["/dockertags"]