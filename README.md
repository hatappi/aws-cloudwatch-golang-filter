# aws-cloudwatch-golang-filter
aws-cloudwatch-golang-filter sends the request metric to Amazon CloudWatch Metrics using Envoy Golang HTTP filter.

※ This project is built for istio/proxy v1.18.0.

## Usage (local)
### 1. Get AWS credential
This project sends metrics to CloudWatch Metrics. In order to do that, you need to get access key ID and secret access key. The following policy is minimal one to work this project.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "cloudwatch:PutMetricData",
            "Resource": "*"
        }
    ]
}
```

### 2. Start project

```console
$ export AWS_ACCESS_KEY_ID=[access key ID you got step1]
$ export AWS_SECRET_ACCESS_KEY=[secret access key you got step1]
$ docker-compose up
```

### 3. Test

```console
➜ curl localhost:18000/test -H "Host: example.com"
OK
```

After executing the command above, you can confirm metrics on CloudWatch metrics.

```console
$ aws cloudwatch list-metrics --namespace AWSCloudWatchGolangFilter-dev
```
