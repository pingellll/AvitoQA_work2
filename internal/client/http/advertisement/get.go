package advertisement

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetAdvertisements performs GET /advertisement with required pagination params.
func HttpGetAdvertisements(t *testing.T, token string, limit int) *http.Response {
	runner := apiRunner.GetRunner()
	if token != "" {
		runner = runner.Auth(token)
	}

	if limit <= 0 {
		limit = 50
	}

	return runner.Create().
		Get(path.AdvertisementPath).
		Query("limit", strconv.Itoa(limit)).
		ContentType("application/json").
		Expect(t).
		End().
		Response
}
