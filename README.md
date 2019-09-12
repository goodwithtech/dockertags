# dockertags
Show tags information by container image name.

Now supporting Docker Hub, GCR (Google Container Registry) and Amazon ECR (Elastic Container Registry).

## Quick Start

```bash
$ brew install goodwithtech/r/dockertags
$ dockertags [IMAGE_NAME]
```

## Examples

```bash
$ dockertags goodwithtech/dockle
+-----------------------------+---------+-------+----------------------+-------------+
|            FULL             |   TAG   | SIZE  |      CREATED AT      | UPLOADED AT |
+-----------------------------+---------+-------+----------------------+-------------+
| goodwithtech/dockle:v0.1.16 | v0.1.16 | 20.5M | 2019-08-25T06:35:40Z | NULL        |
| goodwithtech/dockle:latest  | latest  | 20.5M | 2019-08-25T06:15:39Z | NULL        |
| goodwithtech/dockle:v0.1.15 | v0.1.15 | 20.5M | 2019-08-16T00:01:43Z | NULL        |
| goodwithtech/dockle:v0.1.14 | v0.1.14 | 20.5M | 2019-07-11T15:11:08Z | NULL        |
| goodwithtech/dockle:v0.1.13 | v0.1.13 | 20.5M | 2019-06-19T13:32:17Z | NULL        |
| goodwithtech/dockle:v0.1.12 | v0.1.12 | 20.5M | 2019-06-16T18:49:37Z | NULL        |
| goodwithtech/dockle:v0.1.11 | v0.1.11 | 20.5M | 2019-06-16T17:58:23Z | NULL        |
| goodwithtech/dockle:v0.1.10 | v0.1.10 | 20.5M | 2019-06-15T14:38:20Z | NULL        |
| goodwithtech/dockle:v0.1.9  | v0.1.9  | 20.5M | 2019-06-15T14:11:43Z | NULL        |
| goodwithtech/dockle:v0.1.8  | v0.1.8  | 20.5M | 2019-06-14T19:20:04Z | NULL        |
+-----------------------------+---------+-------+----------------------+-------------+

$ dockertags --limit 20 debian
+-------------------------------+------------------------+-------+----------------------+-------------+
|             FULL              |          TAG           | SIZE  |      CREATED AT      | UPLOADED AT |
+-------------------------------+------------------------+-------+----------------------+-------------+
| debian:rc-buggy               | rc-buggy               | 43.3M | 2019-09-12T01:07:58Z | NULL        |
| debian:experimental-20190910  | experimental-20190910  | 43.3M | 2019-09-12T01:07:36Z | NULL        |
| debian:experimental           | experimental           | 43.3M | 2019-09-12T01:07:32Z | NULL        |
| debian:bullseye-20190910-slim | bullseye-20190910-slim | 43.3M | 2019-09-12T01:07:12Z | NULL        |
| debian:bullseye-20190910      | bullseye-20190910      | 43.3M | 2019-09-12T01:07:07Z | NULL        |
| debian:bullseye               | bullseye               | 43.3M | 2019-09-12T01:07:03Z | NULL        |
| debian:9.11-slim              | 9.11-slim              | 43.3M | 2019-09-12T01:06:52Z | NULL        |
| debian:9.11                   | 9.11                   | 43.3M | 2019-09-12T01:06:48Z | NULL        |
| debian:9-slim                 | 9-slim                 | 43.3M | 2019-09-12T01:06:44Z | NULL        |
| debian:9                      | 9                      | 43.3M | 2019-09-12T01:06:39Z | NULL        |
| debian:stable-20190910-slim   | stable-20190910-slim   | 43.3M | 2019-09-11T23:56:41Z | NULL        |
| debian:stable-20190910        | stable-20190910        | 43.3M | 2019-09-11T23:56:09Z | NULL        |
| debian:stable                 | stable                 | 43.3M | 2019-09-11T23:56:02Z | NULL        |
| debian:sid-slim               | sid-slim               | 43.3M | 2019-09-11T23:55:29Z | NULL        |
| debian:sid-20190910-slim      | sid-20190910-slim      | 43.3M | 2019-09-11T23:55:24Z | NULL        |
| debian:sid-20190910           | sid-20190910           | 43.3M | 2019-09-11T23:54:50Z | NULL        |
| debian:sid                    | sid                    | 43.3M | 2019-09-11T23:54:45Z | NULL        |
| debian:rc-buggy-20190910      | rc-buggy-20190910      | 43.3M | 2019-09-11T23:54:20Z | NULL        |
| debian:oldstable-slim         | oldstable-slim         | 43.3M | 2019-09-11T23:53:32Z | NULL        |
| debian:oldstable-backports    | oldstable-backports    | 43.3M | 2019-09-11T23:53:28Z | NULL        |
+-------------------------------+------------------------+-------+----------------------+-------------+
```

## Options

```
USAGE:
  dockertags [options] image_name
OPTIONS:
  --all                          fetch all tagged image information
  --limit value, -l value        Set max fetch count (default: 50)
  --timeout value, -t value      e.g)5s, 1m (default: 10s)
  --username value, -u value     Username
  --password value, -p value     Using -password via CLI is insecure. Be careful.
  --authurl value, --auth value  Url when fetch authentication
  --debug, -d                    Show debug logs
  --help, -h                     show help
  --version, -v                  print the version
```