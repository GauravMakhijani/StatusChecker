package statuschecker

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"StatusChecker/db"
	"StatusChecker/db/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	service Service
	repo    *mocks.StatusStorer
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.repo = &mocks.StatusStorer{}
	suite.service = New(suite.repo)
}

func (suite *ServiceTestSuite) TearDownSuite() {
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestStatusCheckerAppend() {
	t := suite.T()
	t.Run("when append is successful", func(t *testing.T) {
		website := db.WebsiteStatus{
			Link:   "https://www.google.com",
			ID:     1,
			Status: "",
		}

		suite.repo.On("InsertWebsite", website).Return(nil)

		err := suite.service.Add(context.Background(), website)
		require.NoError(t, err)
	})
}

func (suite *ServiceTestSuite) TestStatusCheckerGetAll() {
	t := suite.T()
	t.Run("when GetAll is successful", func(t *testing.T) {
		res := []db.WebsiteStatus{
			{ID: 1, Link: "google", Status: "UP"},
		}

		suite.repo.On("GetAll").Return(res, nil)

		testres, err := suite.service.GetAll(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, testres)
	})
}

func (suite *ServiceTestSuite) TestStatusCheckerGetSimilar() {
	t := suite.T()
	t.Run("when GetSimilar is successful", func(t *testing.T) {

		inputWebsite := "google.com"
		res := []db.WebsiteStatus{
			{ID: 1, Link: "google.com", Status: "UP"},
		}

		suite.repo.On("GetWebsiteStatus", inputWebsite).Return(res, nil)

		testres, err := suite.service.GetSimilar(context.Background(), inputWebsite)
		require.NoError(t, err)
		assert.NotNil(t, testres)
		assert.Equal(t, inputWebsite, testres[0].Link)
	})

	t.Run("when GetSimilar is unsuccessful", func(t *testing.T) {

		inputWebsite := "google"
		res := []db.WebsiteStatus{
			{ID: 1, Link: "google.com", Status: "UP"},
		}

		suite.repo.On("GetWebsiteStatus", inputWebsite).Return(res, errors.New("error occured in database"))

		testres, err := suite.service.GetSimilar(context.Background(), inputWebsite)
		require.Error(t, err)
		assert.Empty(t, testres)
	})

}

func (suite *ServiceTestSuite) TestStatusCheckerGetStatus() {
	t := suite.T()

	//postive testcase
	t.Run("when GetStatus is successful", func(t *testing.T) {
		inputUrl := "www.google.com"
		expected := "UP"

		resp, err := http.Get("http://" + inputUrl) // doubt

		testres := suite.service.GetStatus(context.Background(), inputUrl)

		require.NoError(t, err)

		assert.Equal(t, testres, expected)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
	})

	//negative testcase

	t.Run("when GetStatus is Unsuccessful", func(t *testing.T) {

		inputUrl := "www.fakewebite1.com"
		expected := "DOWN"

		resp, err := http.Get("http://" + inputUrl) // doubt

		testres := suite.service.GetStatus(context.Background(), inputUrl)

		require.Error(t, err)
		assert.Equal(t, testres, expected)
		assert.Nil(t, resp)
	})
}

// func (suite *ServiceTestSuite) TestStatusCheckerCheckStatus(){
// 	t := suite.T()
// 	ticker := time.NewTicker(time.Second)
// 	go suite.service.CheckStatus(ticker)

// 	suite.repo.UpdateWebsiteStatus()
// }
