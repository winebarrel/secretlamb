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

func TestParametersGet(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/systemsmanager/parameters/get/?name=foo", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"Parameter": {
						"Name": "MyStringParameter",
						"Type": "String",
						"Value": "Veni",
						"Version": 1,
						"LastModifiedDate": "1530018761.888",
						"ARN": "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
						"DataType": "text"
				}
			}
		`), nil
	})

	p, err := secretlamb.NewParameters()
	require.NoError(err)
	value, err := p.Get("foo")
	require.NoError(err)

	assert.Equal(
		&secretlamb.ParameterOutput{
			Parameter: secretlamb.ParameterOutputParameter{
				Name:             "MyStringParameter",
				Type:             "String",
				Value:            "Veni",
				Version:          1,
				LastModifiedDate: "1530018761.888",
				Arn:              "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
				DataType:         "text",
			},
		},
		value,
	)
}

func TestParametersGetErr(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/systemsmanager/parameters/get/?name=foo", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusBadRequest, "not ready to serve traffic, please wait"), nil
	})

	p, err := secretlamb.NewParameters()
	require.NoError(err)
	_, err = p.Get("foo")
	assert.ErrorContains(err, "failed to get parameter - http request error: 400: not ready to serve traffic, please wait")
}

func TestParametersGetWithPortEnv(t *testing.T) {
	t.Setenv("PARAMETERS_SECRETS_EXTENSION_HTTP_PORT", "7777")
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:7777/systemsmanager/parameters/get/?name=foo", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"Parameter": {
						"Name": "MyStringParameter",
						"Type": "String",
						"Value": "Veni",
						"Version": 1,
						"LastModifiedDate": "1530018761.888",
						"ARN": "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
						"DataType": "text"
				}
			}
		`), nil
	})

	p, err := secretlamb.NewParameters()
	require.NoError(err)
	value, err := p.Get("foo")
	require.NoError(err)

	assert.Equal(
		&secretlamb.ParameterOutput{
			Parameter: secretlamb.ParameterOutputParameter{
				Name:             "MyStringParameter",
				Type:             "String",
				Value:            "Veni",
				Version:          1,
				LastModifiedDate: "1530018761.888",
				Arn:              "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
				DataType:         "text",
			},
		},
		value,
	)
}

func TestParametersGetWithDecryption(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/systemsmanager/parameters/get/?name=foo&withDecryption=true", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"Parameter": {
						"Name": "MyStringParameter",
						"Type": "String",
						"Value": "Veni",
						"Version": 1,
						"LastModifiedDate": "1530018761.888",
						"ARN": "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
						"DataType": "text"
				}
			}
		`), nil
	})

	p, err := secretlamb.NewParameters()
	require.NoError(err)
	value, err := p.GetWithDecryption("foo")
	require.NoError(err)

	assert.Equal(
		&secretlamb.ParameterOutput{
			Parameter: secretlamb.ParameterOutputParameter{
				Name:             "MyStringParameter",
				Type:             "String",
				Value:            "Veni",
				Version:          1,
				LastModifiedDate: "1530018761.888",
				Arn:              "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
				DataType:         "text",
			},
		},
		value,
	)
}

func TestParametersWithOptions(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/systemsmanager/parameters/get/?label=zoo&name=foo&version=1&withDecryption=true", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"Parameter": {
						"Name": "MyStringParameter",
						"Type": "String",
						"Value": "Veni",
						"Version": 1,
						"LastModifiedDate": "1530018761.888",
						"ARN": "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
						"DataType": "text"
				}
			}
		`), nil
	})

	p, err := secretlamb.NewParameters()
	require.NoError(err)
	value, err := p.Get("foo",
		secretlamb.ParameterVersion(1),
		secretlamb.ParameterLabel("zoo"),
		secretlamb.ParameterWithDecryption(),
	)
	require.NoError(err)

	assert.Equal(
		&secretlamb.ParameterOutput{
			Parameter: secretlamb.ParameterOutputParameter{
				Name:             "MyStringParameter",
				Type:             "String",
				Value:            "Veni",
				Version:          1,
				LastModifiedDate: "1530018761.888",
				Arn:              "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
				DataType:         "text",
			},
		},
		value,
	)
}

func TestParametersWithEncode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost:2773/systemsmanager/parameters/get/?label=foo%2Fvar&name=%E3%81%82&version=1&withDecryption=true", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, `
			{
				"Parameter": {
						"Name": "MyStringParameter",
						"Type": "String",
						"Value": "Veni",
						"Version": 1,
						"LastModifiedDate": "1530018761.888",
						"ARN": "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
						"DataType": "text"
				}
			}
		`), nil
	})

	p, err := secretlamb.NewParameters()
	require.NoError(err)
	value, err := p.Get("あ",
		secretlamb.ParameterVersion(1),
		secretlamb.ParameterLabel("foo/var"),
		secretlamb.ParameterWithDecryption(),
	)
	require.NoError(err)

	assert.Equal(
		&secretlamb.ParameterOutput{
			Parameter: secretlamb.ParameterOutputParameter{
				Name:             "MyStringParameter",
				Type:             "String",
				Value:            "Veni",
				Version:          1,
				LastModifiedDate: "1530018761.888",
				Arn:              "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
				DataType:         "text",
			},
		},
		value,
	)
}

func TestParametersGetWithRetry(t *testing.T) {
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
					"Parameter": {
							"Name": "MyStringParameter",
							"Type": "String",
							"Value": "Veni",
							"Version": 1,
							"LastModifiedDate": "1530018761.888",
							"ARN": "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
							"DataType": "text"
					}
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

	p, err := secretlamb.NewParameters()
	require.NoError(err)
	p = p.WithRetry(2)
	value, err := p.Get("foo")
	require.NoError(err)
	assert.Equal(2, try)

	assert.Equal(
		&secretlamb.ParameterOutput{
			Parameter: secretlamb.ParameterOutputParameter{
				Name:             "MyStringParameter",
				Type:             "String",
				Value:            "Veni",
				Version:          1,
				LastModifiedDate: "1530018761.888",
				Arn:              "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
				DataType:         "text",
			},
		},
		value,
	)
}
