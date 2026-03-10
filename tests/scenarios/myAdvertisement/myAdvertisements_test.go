package myAdvertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK)
	})

	var items gjson.Result
	s.Run("Проверяем, что у нас есть несколько добавленных ранее объявлений", func() {
		items = gjson.Get(advertisementsBody, "items")
		require.NotEmpty(s.T(), items.Array())
	})

	s.Run("Проверяем, что наши объявления имеют принадлежность к нашему пользователю и объявления имеют НЕ пустые параметры", func() {
		for _, item := range items.Array() {
			require.Equal(s.T(), s.loginData.User.Id, gjson.Get(item.String(), "user_id").String(), "user_id не совпадает")

			require.NotEmpty(s.T(), gjson.Get(item.String(), "id").String(), "id пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "created_at").String(), "created_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "updated_at").String(), "updated_at пустое")
			require.Greater(s.T(), len(gjson.Get(item.String(), "title").String()), 0, "Название пустое")
			require.Greater(s.T(), len(gjson.Get(item.String(), "description").String()), 0, "Описание пустое")
			require.True(s.T(), len(gjson.Get(item.String(), "photos").Array()) > 0, "Присутствуют пустые фотографии")
		}
	})
}

func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
