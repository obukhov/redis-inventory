package server

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ServerTestSuite struct {
	suite.Suite
}

func (suite *ServerTestSuite) TestRenderParam() {
	srv := NewServer(zerolog.Nop())

	go srv.Serve(9999, "test content to serve")
	time.Sleep(time.Millisecond)

	req, _ := http.NewRequest("GET", "/", nil)

	recorder := httptest.NewRecorder()
	srv.ServeHTTP(recorder, req)

	body, _ := ioutil.ReadAll(recorder.Result().Body)
	suite.Assert().Equal("test content to serve", string(body))
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
