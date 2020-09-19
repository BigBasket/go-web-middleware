package tracker

import (
	"context"
	"net/http"

	"github.com/xtgo/uuid"
)

//XTracker - is the web request tracker to track the api request flow through your application.
const XTracker string = "x-tracker"

var ctxKey interface{} = XTracker

// InstrumentXTracker is the middleware which sets the tracker to your request context
func InstrumentXTracker(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		xTracker := r.Header.Get(XTracker) // fetches the x-tracker from the headers.
		_, err := uuid.Parse(xTracker)     // validate if the trakcer is a valid uuid
		if err != nil {
			xTracker = uuid.NewRandom().String() // if this is the first endpoint, will create the tracker.
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKey, xTracker) // set the tracker to the context to use it api flow
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
	return fn
}

// FetchXTracker - returns the xtracker set in the context
func FetchXTracker(ctx context.Context) string {
	tracker := ctx.Value(XTracker)
	if tracker == nil {
		return ""
	}
	return ctx.Value(XTracker).(string)
}
