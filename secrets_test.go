package secretlamb_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/secretlamb"
)

func TestSecretsGet(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/secretsmanager/get?secretId=foo", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"ARN": "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
				"Name": "MyTestSecret",
				"VersionId": "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
				"SecretString": "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
				"VersionStages": [
						"AWSPREVIOUS"
				],
				"CreatedDate": "1523477145.713"
			}
		`), nil
	})

	client, err := secretlamb.NewSecrets()
	require.NoError(err)
	value, err := client.Get("foo")
	require.NoError(err)

	assert.Equal(
		&secretlamb.SecretOutput{
			Arn:           "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
			Name:          "MyTestSecret",
			VersionID:     "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
			SecretString:  "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
			VersionStages: []string{"AWSPREVIOUS"},
			CreatedDate:   "1523477145.713",
		},
		value,
	)
}

func TestTestSecretsGetErr(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/secretsmanager/get?secretId=foo", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusBadRequest, "not ready to serve traffic, please wait"), nil
	})

	client, err := secretlamb.NewSecrets()
	require.NoError(err)
	_, err = client.Get("foo")
	assert.ErrorContains(err, "failed to get secret - http request error: 400 Bad Request: not ready to serve traffic, please wait")
}

func TestSecretsGetWithPortEnv(t *testing.T) {
	t.Setenv("PARAMETERS_SECRETS_EXTENSION_HTTP_PORT", "17777")
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:17777/secretsmanager/get?secretId=foo", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"ARN": "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
				"Name": "MyTestSecret",
				"VersionId": "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
				"SecretString": "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
				"VersionStages": [
						"AWSPREVIOUS"
				],
				"CreatedDate": "1523477145.713"
			}
		`), nil
	})

	client, err := secretlamb.NewSecrets()
	require.NoError(err)
	value, err := client.Get("foo")
	require.NoError(err)

	assert.Equal(
		&secretlamb.SecretOutput{
			Arn:           "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
			Name:          "MyTestSecret",
			VersionID:     "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
			SecretString:  "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
			VersionStages: []string{"AWSPREVIOUS"},
			CreatedDate:   "1523477145.713",
		},
		value,
	)
}

func TestSecretsGetWithOptions(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/secretsmanager/get?secretId=foo&versionId=bar&versionStage=zoo", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"ARN": "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
				"Name": "MyTestSecret",
				"VersionId": "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
				"SecretString": "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
				"VersionStages": [
						"AWSPREVIOUS"
				],
				"CreatedDate": "1523477145.713"
			}
		`), nil
	})

	client, err := secretlamb.NewSecrets()
	require.NoError(err)
	value, err := client.Get("foo",
		secretlamb.SecretVersionId("bar"),
		secretlamb.SecretVersionStage("zoo"),
	)
	require.NoError(err)

	assert.Equal(
		&secretlamb.SecretOutput{
			Arn:           "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
			Name:          "MyTestSecret",
			VersionID:     "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
			SecretString:  "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
			VersionStages: []string{"AWSPREVIOUS"},
			CreatedDate:   "1523477145.713",
		},
		value,
	)
}

func TestSecretsGetWithEncode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/secretsmanager/get?secretId=%E3%81%82&versionId=foo%2Fbar&versionStage=zoo%26baz", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"ARN": "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
				"Name": "MyTestSecret",
				"VersionId": "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
				"SecretString": "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
				"VersionStages": [
						"AWSPREVIOUS"
				],
				"CreatedDate": "1523477145.713"
			}
		`), nil
	})

	client, err := secretlamb.NewSecrets()
	require.NoError(err)
	value, err := client.Get("あ",
		secretlamb.SecretVersionId("foo/bar"),
		secretlamb.SecretVersionStage("zoo&baz"),
	)
	require.NoError(err)

	assert.Equal(
		&secretlamb.SecretOutput{
			Arn:           "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
			Name:          "MyTestSecret",
			VersionID:     "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
			SecretString:  "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
			VersionStages: []string{"AWSPREVIOUS"},
			CreatedDate:   "1523477145.713",
		},
		value,
	)
}

func TestSecretsGetWithRetry(t *testing.T) {
	t.Setenv("PARAMETERS_SECRETS_EXTENSION_HTTP_PORT", "12773")

	assert := assert.New(t)
	require := require.New(t)

	try := 0

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if try < 2 {
			w.WriteHeader(http.StatusBadRequest)
			try++
		} else {
			fmt.Fprintln(w, `
				{
					"ARN": "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
					"Name": "MyTestSecret",
					"VersionId": "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
					"SecretString": "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
					"VersionStages": [
							"AWSPREVIOUS"
					],
					"CreatedDate": "1523477145.713"
				}
		`)
		}
	})

	l, _ := net.Listen("tcp", ":12773")
	ts := httptest.Server{
		Listener: l,
		Config:   &http.Server{Handler: handler},
	}
	ts.Start()
	defer ts.Close()

	s, err := secretlamb.NewSecrets()
	require.NoError(err)
	s = s.WithRetry(2)
	value, err := s.Get("foo")
	require.NoError(err)
	assert.Equal(2, try)

	assert.Equal(
		&secretlamb.SecretOutput{
			Arn:           "arn:aws:secretsmanager:us-west-2:123456789012:secret:MyTestSecret-a1b2c3",
			Name:          "MyTestSecret",
			VersionID:     "a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
			SecretString:  "{\"user\":\"diegor\",\"password\":\"PREVIOUS-EXAMPLE-PASSWORD\"}",
			VersionStages: []string{"AWSPREVIOUS"},
			CreatedDate:   "1523477145.713",
		},
		value,
	)
}
