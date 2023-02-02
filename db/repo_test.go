package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

type DbTestSuite struct {
	suite.Suite
	repo StatusStorer
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

func (suite *DbTestSuite) TestGetAll() {
	t := suite.T()
	db, mock, err := sqlxmock.Newx()
	require.Nil(t, err)

	defer db.Close()

	expectedWebsites := []WebsiteStatus{
		{ID: 1, Link: "https://www.google.com", Status: "UP"},
		{ID: 2, Link: "https://www.facebook.com", Status: "UP"},
	}

	rows := sqlxmock.NewRows([]string{"id", "link", "status"}).AddRow(expectedWebsites[0].ID, expectedWebsites[0].Link, expectedWebsites[0].Status).AddRow(expectedWebsites[1].ID, expectedWebsites[1].Link, expectedWebsites[1].Status)

	mock.ExpectQuery("SELECT id,link,status FROM links").WillReturnRows(rows)
	suite.repo = New(db)
	websites, err := suite.repo.GetAll()

	require.Nil(t, err)

	assert.Equal(t, expectedWebsites, websites)

}

func (suite *DbTestSuite) TestGetWebsiteStatus() {

	t := suite.T()

	db, mock, err := sqlxmock.Newx()
	require.Nil(t, err)
	defer db.Close()

	expectedWebsites := []WebsiteStatus{
		{ID: 1, Link: "https://www.instagram.com", Status: "UP"},
		{ID: 2, Link: "https://www.instamart.com", Status: "UP"},
	}
	websitename := "insta"
	rows := sqlxmock.NewRows([]string{"id", "link", "status"}).AddRow(expectedWebsites[0].ID, expectedWebsites[0].Link, expectedWebsites[0].Status).AddRow(expectedWebsites[1].ID, expectedWebsites[1].Link, expectedWebsites[1].Status)

	mock.ExpectQuery("SELECT id,link,status FROM links WHERE link LIKE ").WithArgs("%" + websitename + "%").WillReturnRows(rows)
	suite.repo = New(db)
	websites, err := suite.repo.GetWebsiteStatus(websitename)

	require.Nil(t, err)

	assert.Equal(t, expectedWebsites, websites)

}

func (suite *DbTestSuite) TestInsertWebsite() {
	t := suite.T()
	db, mock, err := sqlxmock.Newx()
	require.Nil(t, err)

	defer db.Close()

	WebsiteInserted := WebsiteStatus{
		ID: 1, Link: "https://www.google.com", Status: "UP",
	}

	mock.ExpectExec("INSERT INTO links").WithArgs(WebsiteInserted.Link, WebsiteInserted.Status).WillReturnResult(sqlxmock.NewResult(1, 1))

	suite.repo = New(db)
	err = suite.repo.InsertWebsite(WebsiteInserted)
	require.Nil(t, err)

}

func (suite *DbTestSuite) TestUpdateWebsiteStatus() {
	t := suite.T()
	db, mock, err := sqlxmock.Newx()
	require.Nil(t, err)

	defer db.Close()

	WebsiteInserted := WebsiteStatus{
		ID: 1, Link: "https://www.google.com", Status: "UP",
	}

	mock.ExpectExec("UPDATE links").WithArgs(WebsiteInserted.Link, WebsiteInserted.Status).WillReturnResult(sqlxmock.NewResult(1, 1))

	suite.repo = New(db)
	err = suite.repo.UpdateWebsiteStatus(WebsiteInserted.Link, WebsiteInserted.Status)
	require.Nil(t, err)

}
