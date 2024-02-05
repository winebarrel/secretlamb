package secretlamb_test

import (
	"net/http"
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
						"LastModifiedDate": 1530018761.888,
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
				LastModifiedDate: 1530018761.888,
				Arn:              "arn:aws:ssm:us-east-2:111222333444:parameter/MyStringParameter",
				DataType:         "text",
			},
		},
		value,
	)
}
