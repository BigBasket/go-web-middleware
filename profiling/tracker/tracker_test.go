package tracker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xtgo/uuid"
)

type trackerTestSuite struct {
	suite.Suite
}

func (t *trackerTestSuite) SetupSuite() {

}

func TestXTrackerMiddlewareSuite(t *testing.T) {
	t.Parallel()
	xSuite := &trackerTestSuite{}
	suite.Run(t, xSuite)
}

func (t *trackerTestSuite) TestExistTrackerMiddleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		xTracker := FetchXTracker(r.Context())
		if _, err := uuid.Parse(xTracker); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	req.Header.Add(XTracker, uuid.NewRandom().String())
	res := httptest.NewRecorder()
	xHandler := InstrumentXTracker(handler)
	xHandler.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusOK)
}

func (t *trackerTestSuite) TestNewTrackerMiddleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		xTracker := FetchXTracker(r.Context())
		if _, err := uuid.Parse(xTracker); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	res := httptest.NewRecorder()
	xHandler := InstrumentXTracker(handler)
	xHandler.ServeHTTP(res, req)
	t.Assert().Equal(res.Code, http.StatusOK)
}
