# secretlamb

[![CI](https://github.com/winebarrel/secretlamb/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/secretlamb/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/winebarrel/secretlamb.svg)](https://pkg.go.dev/github.com/winebarrel/secretlamb)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/winebarrel/secretlamb)](https://github.com/winebarrel/secretlamb/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/winebarrel/secretlamb)](https://goreportcard.com/report/github.com/winebarrel/secretlamb)

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
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/winebarrel/secretlamb"
)

func HandleRequest(ctx context.Context, event any) (*string, error) {
	client := secretlamb.MustNewParameters() // .WithRetry(3)

	v, err := client.Get("foo")
	//v, err := client.GetWithDecryption("foo")

	if err != nil {
		return nil, err
	}

	fmt.Println(v.Parameter.Value)
	return nil, nil
}

func main() {
	lambda.Start(HandleRequest)
}
```

### Secrets Manager

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/winebarrel/secretlamb"
)

func HandleRequest(ctx context.Context, event any) (*string, error) {
	client := secretlamb.MustNewSecrets() // .WithRetry(3)
	v, err := client.Get("foo")

	if err != nil {
		return nil, err
	}

	fmt.Println(v.SecretString)
	return nil, nil
}

func main() {
	lambda.Start(HandleRequest)
}
```
