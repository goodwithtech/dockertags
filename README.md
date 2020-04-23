<img src="assets/dockertags.png" width="300" />

Show information of container images ordered by recently updated. <br />
Now supporting Docker Hub, GCR (Google Container Registry) and Amazon ECR (Elastic Container Registry).

## Quick Start

```bash
$ brew install goodwithtech/r/dockertags
$ dockertags [IMAGE_NAME]

or 

$ docker run --rm goodwithtech/dockertags [IMAGENAME]
```
## When to Use

Make easy to fetch target tag in scheduled operation.
 
```.env
$ dockertags -limit 1 -format json <imagename> | jq -r .[0].tags[0]
...output tag...

# Scan a latest container image with https://github.com/aquasecurity/trivy
$ export IMAGENAME=<imagename>
$ trivy $IMAGENAME:$(dockertags -limit 1 -format json $IMAGENAME | jq -r .[0].tags[0])
```

## Examples

```bash

$ dockertags alpine
+--------------------------------+--------------------------------+--------------------------------+--------------------------------+------------+----------------------+
|              TAG               |              SIZE              |             DIGEST             |            OS/ARCH             | CREATED AT |     UPLOADED AT      |
+--------------------------------+--------------------------------+--------------------------------+--------------------------------+------------+----------------------+
| 3.9                            | 2.7M                           | fa5361fbf636                   | linux/ppc64le                  | NULL       | 2020-01-23T17:42:13Z |
| 3.9.5                          | 2.6M                           | c7b3e8392e08                   | linux/386                      |            |                      |
|                                | 2.6M                           | cae6522b6a35                   | linux/arm64                    |            |                      |
|                                | 2.2M                           | 0c6b515386fd                   | linux/arm                      |            |                      |
|                                | 2.4M                           | 97e9e9a15ef9                   | linux/s390x                    |            |                      |
|                                | 2.4M                           | 5292cebaf695                   | linux/arm                      |            |                      |
|                                | 2.6M                           | ab3fe83c0696                   | linux/amd64                    |            |                      |
+--------------------------------+--------------------------------+--------------------------------+--------------------------------+------------+----------------------+
| 3.8                            | 2M                             | e802987f152d                   | linux/arm64                    | NULL       | 2020-01-23T17:41:40Z |
| 3.8.5                          | 2.1M                           | 402d21757a03                   | linux/ppc64le                  |            |                      |
|                                | 2.2M                           | cf35b4fa14e2                   | linux/386                      |            |                      |
|                                | 2M                             | dabea2944dcc                   | linux/arm                      |            |                      |
|                                | 2.1M                           | 954b378c375d                   | linux/amd64                    |            |                      |
|                                | 2.2M                           | 514ec80ffbe1                   | linux/s390x                    |            |                      |
+--------------------------------+--------------------------------+--------------------------------+--------------------------------+------------+----------------------+
| 3.10                           | 2.7M                           | 33158d51a7a5                   | linux/ppc64le                  | NULL       | 2020-01-23T17:41:10Z |
| 3.10.4                         | 2.7M                           | de78803598bc                   | linux/amd64                    |            |                      |
|                                | 2.5M                           | 9afbfccb8066                   | linux/arm                      |            |                      |
|                                | 2.7M                           | 747f335d2f68                   | linux/386                      |            |                      |
|                                | 2.5M                           | 216161924b52                   | linux/s390x                    |            |                      |
|                                | 2.3M                           | 2632d6288d34                   | linux/arm                      |            |                      |
|                                | 2.6M                           | 4491fd429b8a                   | linux/arm64                    |            |                      |
+--------------------------------+--------------------------------+--------------------------------+--------------------------------+------------+----------------------+
| edge                           | 2.6M                           | fb7bea212348                   | linux/arm64                    | NULL       | 2020-01-23T00:41:09Z |
| 20200122                       | 2.5M                           | e3e522f13253                   | linux/arm                      |            |                      |
|                                | 2.7M                           | 7b5953e862c9                   | linux/ppc64le                  |            |                      |
|                                | 2.3M                           | e137ff293fcc                   | linux/arm                      |            |                      |
|                                | 2.7M                           | 5f60bf03cace                   | linux/386                      |            |                      |
|                                | 2.5M                           | cb3bf0adee89                   | linux/s390x                    |            |                      |
|                                | 2.7M                           | 9898e9a51db3                   | linux/amd64                    |            |                      |
+--------------------------------+--------------------------------+--------------------------------+--------------------------------+------------+----------------------+
| latest                         | 2.6M                           | 4d5c59516695                   | linux/arm64                    | NULL       | 2020-01-18T02:41:10Z |
| 3.11.3                         | 2.3M                           | 2c26a655f6e3                   | linux/arm                      |            |                      |
|                                | 2.5M                           | 401f030aa35e                   | linux/arm                      |            |                      |
|                                | 2.7M                           | ddba4d27a7ff                   | linux/amd64                    |            |                      |
|                                | 2.5M                           | ef20eb43010a                   | linux/s390x                    |            |                      |
|                                | 2.7M                           | c40c013324aa                   | linux/386                      |            |                      |
|                                | 2.7M                           | ff8a6adf5859                   | linux/ppc64le                  |            |                      |
+--------------------------------+--------------------------------+--------------------------------+--------------------------------+------------+----------------------+

# You can set limit, filter and format
$ dockertags  -limit 1 -contain latest -format json alpine
[
  {
    "tags": [
      "latest",
      "3.11.3",
      "3.11",
      "3"
    ],
    "data": [
      {
        "Os": "linux",
        "Arch": "ppc64le",
        "Digest": "sha256:ff8a6adf5859433869343296f1b06e0a7bdf4fc836b08d5854221e351baf6929",
        "byte": 2822125
      },
      {
        "Os": "linux",
        "Arch": "arm64",
        "Digest": "sha256:4d5c5951669588e23881c158629ae6bac4ab44866d5b4d150c3f15d91f26682b",
        "byte": 2723075
      },
      {
        "Os": "linux",
        "Arch": "s390x",
        "Digest": "sha256:ef20eb43010abda2d7944e0cd11ef00a961ff7a7f953671226fbf8747895417d",
        "byte": 2582031
      },
      {
        "Os": "linux",
        "Arch": "arm",
        "Digest": "sha256:401f030aa35e86bafd31c6cc292b01659cbde72d77e8c24737bd63283837f02c",
        "byte": 2617562
      },
      {
        "Os": "linux",
        "Arch": "386",
        "Digest": "sha256:c40c013324aa73f430d33724d8030c34b1881e96b23f44ec616f1caf8dbf445f",
        "byte": 2806560
      },
      {
        "Os": "linux",
        "Arch": "amd64",
        "Digest": "sha256:ddba4d27a7ffc3f86dd6c2f92041af252a1f23a8e742c90e6e1297bfa1bc0c45",
        "byte": 2802957
      },
      {
        "Os": "linux",
        "Arch": "arm",
        "Digest": "sha256:2c26a655f6e38294e859edac46230210bbed3591d6ff57060b8671cda09756d4",
        "byte": 2419554
      }
    ],
    "created_at": "0001-01-01T00:00:00Z",
    "uploaded_at": "2020-01-18T02:41:10.850638Z"
  }
]
```

## Usage

```
NAME:
  dockertags - fetch docker image tags
  dockertags [options] image_name
OPTIONS:
  --limit value, -l value        set max tags count. if exist no tag image will be short numbers. limit=0 means fetch all tags (default: 10)
  --contain value, -c value      contains target string. multiple string allows.
  --format value, -f value       target format table or json, default table (default: "table")
  --output value, -o value       output file name, default output to stdout
  --authurl value, --auth value  GetURL when fetch authentication
  --timeout value, -t value      e.g)5s, 1m (default: 10s)
  --username value, -u value     Username
  --password value, -p value     Using -password via CLI is insecure. Be careful.
  --debug, -d                    Show debug logs
  --help, -h                     show help
  --version, -v                  print the version
```

## GitHub Actions

You can scan target image everyday recently updated.<br />
This actions also notify results if trivy detects vulnerabilities.

```
name: Scan the target image with trivy
on:
  schedule:
      - cron:  '0 0 * * *'
jobs:
  scan:
    name: Scan via trivy
    runs-on: ubuntu-latest
    env:
      IMAGE: goodwithtech/dockle # target image name
      FILTER: v0.2    # pattern : /*v0.2*/
    steps:
      - name: detect a target image tag
        id: target
        run: echo ::set-output name=ver::$(
            docker run --rm goodwithtech/dockertags -contain $FILTER -limit 1 -format json $IMAGE
            | jq -r .[0].tags[0]
            )
      - name: detect a trivy image tag
        id: trivy
        run: echo ::set-output name=ver::$(
            docker run --rm goodwithtech/dockertags -limit 1 -format json aquasec/trivy
            | jq -r .[0].tags[0]
            )
      - name: check tags
        run: |
          echo trivy ${{ steps.trivy.outputs.ver }}
          echo $IMAGE ${{ steps.target.outputs.ver }}
      - name: scan the image with trivy
        run: docker run aquasec/trivy:${{ steps.trivy.outputs.ver }}
          --cache-dir /var/lib/trivy --exit-code 1 --no-progress
          $IMAGE:${{ steps.target.outputs.ver }}
      - name: notify to slack
        if: failure()
        uses: rtCamp/action-slack-notify@master
        env:
          SLACK_CHANNEL: channel  # target channel
          SLACK_MESSAGE: 'failed : trivy detects vulnerabilities'
          SLACK_TITLE: trivy-scan-notifier
          SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
```

## Authentication

### Docker Hub

You can use `--username` and `--password` of Docker Hub.

```bash
dockertags -u goodwithtech -p xxxx goodwithtech/privateimage
```

### Amazon ECR (Elastic Container Registry)

Use [AWS CLI's ENVIRONMENT variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

```bash
AWS_PROFILE={PROFILE_NAME}
AWS_DEFAULT_REGION={REGION}
```

### GCR (Google Container Registry)

If you'd like to use the target project's repository, you can settle via `GOOGLE_APPLICATION_CREDENTIAL`.

```bash
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credential.json
```

