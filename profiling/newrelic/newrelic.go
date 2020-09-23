package newrelic

import (
	"context"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Connection - newrelic connection params
type Connection struct {
	app *newrelic.Application
}

// ConnectionConfig - config to connect to newrelic
type ConnectionConfig struct {
	newrelic.Config
}

// NewConnection - connection parameters
func NewConnection(ctx context.Context, params ConnectionConfig) (*Connection, error) {
	nApp, err := newrelic.NewApplication(
		func(config *newrelic.Config) {
			*config = params.Config
		},
	)
	if err != nil {
		return nil, err
	}
	return &Connection{app: nApp}, nil
}

// InstrumentNewRelic is the middleware for newrelic
func (nApp *Connection) InstrumentNewRelic(next http.HandlerFunc) http.HandlerFunc {
	if nApp == nil { //NewRelic Connection is not setup
		return http.HandlerFunc(next.ServeHTTP)
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		txn := nApp.app.StartTransaction(r.Method + " " + r.URL.Path)
		defer txn.End()
		w = txn.SetWebResponse(w)
		txn.SetWebRequestHTTP(r)
		//instrument the newrelic transaction into the request context.
		r = newrelic.RequestWithTransactionContext(r, txn)
		next.ServeHTTP(w, r)
	}
	return fn
}

// FetchTransaction - returns the newrelic transaction from the context.
func FetchTransaction(ctx context.Context) *newrelic.Transaction {
	return newrelic.FromContext(ctx)
}
