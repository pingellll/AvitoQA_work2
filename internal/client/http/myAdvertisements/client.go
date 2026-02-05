package myAdvertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpGetMyAdvertisements(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.MyAdvertisementsPath).
		ContentType("application/json").
		Query("limit", "50").
		Expect(t).
		End().Response
}
