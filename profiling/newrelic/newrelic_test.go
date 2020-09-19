package newrelic

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/suite"
)

type newrelicTestSuite struct {
	suite.Suite
}

func (t *newrelicTestSuite) SetupSuite() {

}

func TestRelicMiddlewareSuite(t *testing.T) {
	t.Parallel()
	nSuite := &newrelicTestSuite{}
	suite.Run(t, nSuite)
}

func (t *newrelicTestSuite) TestValidNewRelicMiddleware() {
	var testkey string = "valid_key" //replace with your valid newrelic api key
	ctx := context.Background()
	handler := func(w http.ResponseWriter, r *http.Request) {
		tx := FetchTransaction(r.Context())
		if tx == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	res := httptest.NewRecorder()
	nConn, err := NewConnection(ctx, ConnectionConfig{
		newrelic.Config{AppName: "newrelic-test",
			License: testkey},
	})
	if err != nil {
		t.Fail("newrelic connection failed", err.Error())
		return
	}
	nHandler := nConn.InstrumentNewRelic(handler)
	nHandler.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusOK)
}

func (t *newrelicTestSuite) TestInValidNewRelicMiddleware() {
	var testkey string = "invalidkey"
	ctx := context.Background()
	handler := func(w http.ResponseWriter, r *http.Request) {
		tx := FetchTransaction(r.Context())
		if tx == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	res := httptest.NewRecorder()
	nConn, err := NewConnection(ctx, ConnectionConfig{
		newrelic.Config{AppName: "newrelic-test",
			License: testkey},
	})
	if err == nil {
		t.Fail("newrelic connection succeed", err.Error())
	}
	nHandler := nConn.InstrumentNewRelic(handler)
	nHandler.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusNotFound)
}
