package jwt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/stretchr/testify/suite"
)

type rsaTestSuite struct {
	suite.Suite
	token     string
	publickey string
}

func TestRSAMiddlewareSuite(t *testing.T) {
	tSuite := &rsaTestSuite{}
	suite.Run(t, tSuite)
}

func (t *rsaTestSuite) SetupSuite() {
	t.token = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOiIyMDMwLTA4LTE4VDEwOjI3OjIxLjg0NjgwNzg5WiIsInVzZXJpZCI6ImtvbWFsIiwidXNlcm5hbWUiOiJrb21hbCJ9.q7bymi6hXMtI-UEopespC8THSSTtdk6NEwwJCZgGWBAUP7laUVXn8KTNF9fcuryvyYEACm27mqVxt7zowqbKJw"
	t.publickey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTC9sRjFyWmdOdGxwUWRaN1R6NW51cFc3MG1lVEw3LwpjOG1JZU03alM1alVHU0M3c3EwYkNqYWRSK1c1Q3lSNHhOWTZ0czNsSWd3OTM0b2pOUitKRkQwQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
}

func (t *rsaTestSuite) TearDownSuite() {
}

func (t *rsaTestSuite) TestRSAValidTokenMiddleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	req.Header.Add(BearerAuthHeader, t.token)
	res := httptest.NewRecorder()
	handler(res, req)
	jwtDecode := &Decode{
		PublicKey: t.publickey,
		Algorithm: jwa.RS256,
	}
	jwtm := jwtDecode.InstrumentJWTBearerToken(handler)
	jwtm.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusOK)
}

func (t *rsaTestSuite) TestRSAMissingAuthHeaderMiddleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	res := httptest.NewRecorder()
	handler(res, req)
	jwtDecode := &Decode{
		PublicKey: t.publickey,
		Algorithm: jwa.RS256,
	}
	jwtm := jwtDecode.InstrumentJWTBearerToken(handler)
	jwtm.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusUnauthorized)
}

func (t *rsaTestSuite) TestRSAInvalidPublicKeyMiddleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	res := httptest.NewRecorder()
	handler(res, req)
	jwtDecode := &Decode{
		PublicKey: "",
		Algorithm: jwa.RS256,
	}
	jwtm := jwtDecode.InstrumentJWTBearerToken(handler)
	jwtm.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusUnauthorized)
}
