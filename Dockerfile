FROM golang:1.16-alpine AS builder
COPY go.mod go.sum /app/
WORKDIR /app/
RUN apk --no-cache add git
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /dockertags cmd/dockertags/main.go

FROM alpine:3.11
COPY --from=builder /dockertags /usr/local/bin/dockertags
RUN chmod +x /usr/local/bin/dockertags
RUN apk --no-cache add ca-certificates shadow

# for use docker daemon via mounted /var/run/docker.sock
RUN addgroup -S docker && adduser -S -G docker dockertags
USER dockertags

ENTRYPOINT ["dockertags"]