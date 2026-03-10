package advertisement

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisement"
)

func GetAdvertisementPhotos(t *testing.T, token string, advertisementID string, expectedStatusCode int) []byte {
	resp := advertisement.HttpGetAdvertisementPhotos(t, token, advertisementID)
	require.Equalf(t, expectedStatusCode, resp.StatusCode, "HTTP status code должен быть %d", expectedStatusCode)

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return data
}
