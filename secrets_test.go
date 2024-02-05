package secretlamb_test

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/secretlamb"
)

func TestSecrets(t *testing.T) {
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
				"CreatedDate": 1523477145.713
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
			CreatedDate:   1523477145.713,
		},
		value,
	)
}
