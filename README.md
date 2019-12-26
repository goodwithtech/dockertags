# dockertags
Show information of container images ordered by recently updated. <br />
Now supporting Docker Hub, GCR (Google Container Registry) and Amazon ECR (Elastic Container Registry).

<img src="assets/usage.gif" width="1200">



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
+----------+------+----------------------+-------------+
|   TAG    | SIZE |      CREATED AT      | UPLOADED AT |
+----------+------+----------------------+-------------+
| 3        | 2.7M | 2019-12-24T20:40:57Z | NULL        |
| 3.11     |      |                      |             |
| latest   |      |                      |             |
| 3.11.2   |      |                      |             |
+----------+------+----------------------+-------------+
| edge     | 2.7M | 2019-12-20T00:41:30Z | NULL        |
| 20191219 |      |                      |             |
+----------+------+----------------------+-------------+
| 3.11.0   | 2.7M | 2019-12-20T00:41:21Z | NULL        |
+----------+------+----------------------+-------------+
| 20191114 | 2.7M | 2019-11-14T22:41:11Z | NULL        |
+----------+------+----------------------+-------------+
| 3.10     | 2.7M | 2019-10-21T18:41:18Z | NULL        |
| 3.10.3   |      |                      |             |
+----------+------+----------------------+-------------+
| 20190925 | 2.7M | 2019-09-25T22:40:50Z | NULL        |
+----------+------+----------------------+-------------+
| 3.10.2   | 2.7M | 2019-08-20T21:40:57Z | NULL        |
+----------+------+----------------------+-------------+
| 3.8      | 2.1M | 2019-08-20T06:41:01Z | NULL        |
| 3.8.4    |      |                      |             |
+----------+------+----------------------+-------------+
| 20190809 | 2.7M | 2019-08-09T21:41:13Z | NULL        |
+----------+------+----------------------+-------------+
| 3.10.1   | 2.7M | 2019-07-11T22:41:17Z | NULL        |
+----------+------+----------------------+-------------+



# You can set limit, filter and format
$ dockertags  -limit 1 -contain latest -format json alpine
[
  {
    "tags": [
      "latest",
      "3.11.2",
      "3.11",
      "3"
    ],
    "byte": 2801778,
    "created_at": "2019-12-24T20:40:57.918177Z",
    "uploaded_at": "0001-01-01T00:00:00Z"
  }
]
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

