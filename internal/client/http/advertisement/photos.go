package advertisement

import (
	"fmt"
	"net/http"
	"testing"

	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetAdvertisementPhotos performs GET /advertisements/{id}/photos.
func HttpGetAdvertisementPhotos(t *testing.T, token string, advertisementID string) *http.Response {
	runner := apiRunner.GetRunner()
	if token != "" {
		runner = runner.Auth(token)
	}

	photosPath := fmt.Sprintf("/advertisements/%s/photos", advertisementID)
	return runner.Create().Get(photosPath).
		Expect(t).
		End().Response
}
