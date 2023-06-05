package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/handler"
)

func (s *TestSuite) TestShortOk() {
	code, rawBody := s.doRequest(s.shortRequest("http://ya.ru"))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Len(shortUrl.Path, 18)

	s.Require().Len(response.SecretKey, 16)
}

func (s *TestSuite) TestShortInvalidLongUrl() {
	// url without schema
	code, _ := s.doRequest(s.shortRequest("ya.ru"))
	s.Require().Equal(http.StatusBadRequest, code)
}

func (s *TestSuite) TestShortDuplicate() {
	// main request
	code, rawBody := s.doRequest(s.shortRequest("http://ya.ru"))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Len(shortUrl.Path, 18)

	s.Require().Len(response.SecretKey, 16)

	// duplicated request
	code, rawBody = s.doRequest(s.shortRequest("http://ya.ru"))
	s.Require().Equal(http.StatusOK, code)

	response, err = s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err = url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Len(shortUrl.Path, 18)

	s.Require().Len(response.SecretKey, 0)
}

func (s *TestSuite) shortRequest(longUrl string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/make_shorter", bytes.NewReader([]byte(fmt.Sprintf(`{"long_url": "%s"}`, longUrl))))
	return req
}

func (s *TestSuite) shortResponseFromBody(body string) (handler.ShortLinkResponse, error) {
	var res handler.ShortLinkResponse
	err := json.Unmarshal([]byte(body), &res)
	return res, err
}
