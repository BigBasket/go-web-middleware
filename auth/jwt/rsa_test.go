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
	t.Parallel()
	tSuite := &rsaTestSuite{}
	suite.Run(t, tSuite)
}

func (t *rsaTestSuite) SetupSuite() {
	t.token = "Bearer <token>"
	t.publickey = "Base64EncodedPublicKey"
}

func (t *rsaTestSuite) TearDownSuite() {
}

func (t *rsaTestSuite) TestRSAValidTokenMiddleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	req.Header.Add(BearerAuthHeader, t.token)
	res := httptest.NewRecorder()
	jwtDecode := &RSADecode{
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
	jwtDecode := &RSADecode{
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
	jwtDecode := &RSADecode{
		PublicKey: "",
		Algorithm: jwa.RS256,
	}
	jwtm := jwtDecode.InstrumentJWTBearerToken(handler)
	jwtm.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusUnauthorized)
}
