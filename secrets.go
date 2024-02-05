package secretlamb

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Secrets struct {
	client *client
}

type SecretOutput struct {
	Arn           string   `json:"ARN"`
	Name          string   `json:"Name"`
	VersionID     string   `json:"VersionId"`
	SecretString  string   `json:"SecretString"`
	VersionStages []string `json:"VersionStages"`
	CreatedDate   float64  `json:"CreatedDate"`
}

type SecretOption struct {
	Key   string
	Value string
}

func SecretVersionId(versionId string) *SecretOption {
	return &SecretOption{
		Key:   "versionId",
		Value: versionId,
	}
}

func SecretVersionStage(versionStage string) *SecretOption {
	return &SecretOption{
		Key:   "versionStage",
		Value: versionStage,
	}
}

func NewSecrets() (*Secrets, error) {
	client, err := newClient("/secretsmanager/get")
	return &Secrets{client: client}, err
}

func MustNewSecrets() *Secrets {
	client, err := NewSecrets()

	if err != nil {
		panic("NewSecrets(): " + err.Error())
	}

	return client
}

func (s *Secrets) Get(secretId string, options ...*SecretOption) (*SecretOutput, error) {
	query := &url.Values{}
	query.Add("secretId", secretId)

	for _, opt := range options {
		query.Add(opt.Key, opt.Value)
	}

	body, err := s.client.get(query)

	if err != nil {
		return nil, fmt.Errorf("failed to get secret - http request error: %w", err)
	}

	output := &SecretOutput{}
	err = json.Unmarshal(body, output)

	if err != nil {
		return nil, fmt.Errorf("failed to get secret - json unmarshal error: %w", err)
	}

	return output, nil
}
