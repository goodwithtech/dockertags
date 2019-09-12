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
```
