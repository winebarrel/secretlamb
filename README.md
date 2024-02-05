# secretlamb


Golang library for using AWS Parameters and Secrets Lambda Extension.

- [Using Parameter Store parameters in AWS Lambda functions - AWS Systems Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/ps-integration-lambda-extensions.html)
- [Use AWS Secrets Manager secrets in AWS Lambda functions - AWS Secrets Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/retrieving-secrets_lambda.html)

## Installation

```sh
go get github.com/winebarrel/secretlamb
```

## Usage

### Parameter Store

```go
client := secretlamb.MustNewParameters()
v, err := client.Get("foo")
//v, err := client.GetWithDecryption("foo")

if err != nil {
	panic(err)
}

fmt.Println(v.Parameter.Value)
```

### Secrets Manager

```go
client := secretlamb.MustNewParameters()
v, err := client.Get("foo")

if err != nil {
	panic(err)
}

fmt.Println(v.SecretString)
```
