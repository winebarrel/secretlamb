package secretlamb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
)

type Parameters struct {
	*client
}

type ParameterOutputParameter struct {
	Name             string `json:"Name"`
	Type             string `json:"Type"`
	Value            string `json:"Value"`
	Version          int64  `json:"Version"`
	LastModifiedDate string `json:"LastModifiedDate"`
	Arn              string `json:"ARN"`
	DataType         string `json:"DataType"`
}

type ParameterOutput struct {
	Parameter ParameterOutputParameter `json:"Parameter"`
}

type ParameterOption struct {
	Key   string
	Value string
}

func ParameterVersion(version int) *ParameterOption {
	return &ParameterOption{
		Key:   "version",
		Value: strconv.Itoa(version),
	}
}

func ParameterLabel(label string) *ParameterOption {
	return &ParameterOption{
		Key:   "label",
		Value: label,
	}
}

func ParameterWithDecryption() *ParameterOption {
	return &ParameterOption{
		Key:   "withDecryption",
		Value: "true",
	}
}

func NewParameters() (*Parameters, error) {
	client, err := newClient("/systemsmanager/parameters/get/")
	return &Parameters{client: client}, err
}

func MustNewParameters() *Parameters {
	client, err := NewParameters()

	if err != nil {
		panic("MustNewParameters(): " + err.Error())
	}

	return client
}

func (p *Parameters) WithRetry(retryMax int) *Parameters {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retryMax
	retryClient.CheckRetry = retryPolicy
	p.HTTPClient = retryClient.StandardClient()
	return p
}

func (p *Parameters) Get(name string, options ...*ParameterOption) (*ParameterOutput, error) {
	return p.GetWithContext(context.Background(), name, options...)
}

func (p *Parameters) GetWithContext(ctx context.Context, name string, options ...*ParameterOption) (*ParameterOutput, error) {
	query := &url.Values{}
	query.Add("name", name)

	for _, opt := range options {
		query.Add(opt.Key, opt.Value)
	}

	body, err := p.get(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get parameter - http request error: %w", err)
	}

	output := &ParameterOutput{}
	err = json.Unmarshal(body, output)

	if err != nil {
		return nil, fmt.Errorf("failed to get parameter - json unmarshal error: %w", err)
	}

	return output, nil
}

func (p *Parameters) GetWithDecryption(name string) (*ParameterOutput, error) {
	return p.Get(name, ParameterWithDecryption())
}
