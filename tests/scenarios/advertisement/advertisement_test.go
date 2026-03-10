package advertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	adv "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TestAdvertisementCreatePositive() {
	base.Precondition("Создаём объявление, проверяем что оно видно в списке и что фото доступно")

	title := "auto-test-title-" + utils.RandomString(10)

	var createdID string
	s.Run("Создаём объявление (201)", func() {
		createBody := adv.Create(
			s.T(),
			s.loginData.Token,
			adv.CreateRequest{
				Title:       title,
				Description: "auto-test-description",
				Price:       "1000",
				Quantity:    "1",
				PhotoPath:   "testdata/9.jpg",
			},
			http.StatusCreated,
		)

		createdID = gjson.Get(createBody, "id").String()
		require.NotEmpty(s.T(), createdID, "id пустой в ответе create")
	})

	s.Run("Проверяем, что объявление доступно через GET /my/advertisements", func() {
		myBody := myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK)

		items := gjson.Get(myBody, "items")
		require.True(s.T(), items.IsArray(), "ожидали массив items в ответе")

		found := false
		for _, item := range items.Array() {
			if gjson.Get(item.String(), "id").String() == createdID {
				found = true
				require.Equal(s.T(), title, gjson.Get(item.String(), "title").String(), "title не совпал")
				break
			}
		}
		require.True(s.T(), found, "созданное объявление не найдено в /my/advertisements")
	})

	s.Run("Проверяем, что фото доступно через /advertisements/{id}/photos", func() {
		photoBytes := adv.GetAdvertisementPhotos(s.T(), s.loginData.Token, createdID, http.StatusOK)
		require.Greater(s.T(), len(photoBytes), 0, "ответ ручки photos пустой")
	})
}

func (s *TestSuite) TestAdvertisementCreateNegativeMissingRequired() {
	base.Precondition("Негативный тест: создаём объявление без обязательных параметров")

	s.Run("Пытаемся создать объявление без title (ожидаем 400)", func() {
		body := adv.Create(
			s.T(),
			s.loginData.Token,
			adv.CreateRequest{
				Title:       "",
				Description: "auto-test-description",
				Price:       "1000",
				Quantity:    "1",
				PhotoPath:   "testdata/9.jpg",
			},
			http.StatusBadRequest,
		)

		require.NotEmpty(s.T(), body)
	})
}

func (s *TestSuite) TestAdvertisementCreateNegativeWithoutAuth() {
	base.Precondition("Негативный тест: создаём объявление без авторизации")

	s.Run("Пытаемся создать объявление без токена (ожидаем 401)", func() {
		body := adv.Create(
			s.T(),
			"",
			adv.CreateRequest{
				Title:       "auto-test-title-" + utils.RandomString(10),
				Description: "auto-test-description",
				Price:       "1000",
				Quantity:    "1",
				PhotoPath:   "testdata/9.jpg",
			},
			http.StatusUnauthorized,
		)
		require.NotEmpty(s.T(), body)
	})
}
