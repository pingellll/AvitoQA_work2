package advertisement

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	httpAdv "api-tests-template/internal/client/http/advertisement"
)

func GetAdvertisements(t *testing.T, token string, limit int, expectedStatusCode int) string {
	resp := httpAdv.HttpGetAdvertisements(t, token, limit)
	require.Equalf(t, expectedStatusCode, resp.StatusCode, "HTTP status code должен быть %d", expectedStatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(bodyBytes)
}
