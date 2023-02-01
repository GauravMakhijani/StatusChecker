package statuschecker

import (
	"StatusChecker/db"
	"StatusChecker/statuschecker/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	service *mocks.Service
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) SetupTest() {
	s.service = &mocks.Service{}
}

func (s *HandlerTestSuite) TearDownTest() {
	t := s.T()
	s.service.AssertExpectations(t)
}

func (s *HandlerTestSuite) TestcreateWebsiteHandler() {

	t := s.T()
	message := "Websites added"

	t.Run("when post request is successful", func(t *testing.T) {
		body := []byte(`{"websites":["www.yahoo.com"]}`)

		r := httptest.NewRequest(http.MethodPost, "/websites", bytes.NewBuffer(body))

		w := httptest.NewRecorder()

		website := db.WebsiteStatus{
			Link:   "www.yahoo.com",
			Status: "UP",
		}

		s.service.On("GetStatus", mock.Anything, website.Link).Return("UP")
		s.service.On("Add", mock.Anything, website).Return(nil)

		// s.service.on("Add",mock.Anything,)

		CreateWebsiteHandler(w, r, s.service)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp string
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, message, resp)
	})

	t.Run("when body format is wrong and post request is unsuccessful", func(t *testing.T) {
		body := []byte(`{"Links" ["www.wikipedia.com"]}`)

		r := httptest.NewRequest(http.MethodPost, "/websites", bytes.NewBuffer(body))

		w := httptest.NewRecorder()

		CreateWebsiteHandler(w, r, s.service)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	})

	t.Run("when error occured while inserting in db in createWebsiteHandler", func(t *testing.T) {
		body := []byte(`{"websites":["www.google.com"]}`)

		r := httptest.NewRequest(http.MethodPost, "/websites", bytes.NewBuffer(body))

		w := httptest.NewRecorder()

		websites := []db.WebsiteStatus{
			{
				Link:   "www.google.com",
				Status: "UP",
			},
		}

		for _, website := range websites {
			s.service.On("GetStatus", mock.Anything, website.Link).Return(website.Status)
			s.service.On("Add", mock.Anything, website).Return(errors.New("error while inserting in db"))
		}

		CreateWebsiteHandler(w, r, s.service)

		var err string
		json.NewDecoder(w.Body).Decode(&err)
		t.Log(err)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	})

}

func (s *HandlerTestSuite) TestGetWebsiteHandler() {

	t := s.T()

	t.Run("when get request is successful", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodGet, "/websites", nil)

		w := httptest.NewRecorder()

		websites := []db.WebsiteStatus{
			{
				Link:   "www.youtube.com",
				Status: "UP",
			},
			{
				Link:   "www.codingninjas.com",
				Status: "UP",
			},
		}

		s.service.On("GetAll", mock.Anything).Return(websites, nil).Once()

		GetWebsiteHandler(w, r, s.service)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp []db.WebsiteStatus
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, websites, resp)
	})
	t.Run("when get request is unsuccessful", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodGet, "/websites", nil)

		w := httptest.NewRecorder()

		websites := []db.WebsiteStatus{}

		s.service.On("GetAll", mock.Anything).Return(websites, errors.New("error occured while fetching webites")).Once()

		GetWebsiteHandler(w, r, s.service)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	})

	t.Run("when passing name get request is successful", func(t *testing.T) {

		websiteName := "instagram"
		r := httptest.NewRequest(http.MethodGet, "/websites?name="+websiteName, nil)

		w := httptest.NewRecorder()

		websites := []db.WebsiteStatus{
			{
				Link:   "www.instagram.com",
				Status: "UP",
			},
		}

		s.service.On("GetSimilar", mock.Anything, websiteName).Return(websites, nil).Once()

		GetWebsiteHandler(w, r, s.service)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp []db.WebsiteStatus
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, websites, resp)
	})
	t.Run("when passing name get request is unsuccessful", func(t *testing.T) {

		websiteName := "instagram"
		r := httptest.NewRequest(http.MethodGet, "/websites?name="+websiteName, nil)

		w := httptest.NewRecorder()

		websites := []db.WebsiteStatus{}

		s.service.On("GetSimilar", mock.Anything, websiteName).Return(websites, errors.New("error occured while fetching webites")).Once()

		GetWebsiteHandler(w, r, s.service)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	})

}
