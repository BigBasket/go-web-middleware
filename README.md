# gopkg

# Import Package
go get github.com/BigBasket/go-web-middleware

The Objective is to manage common middlewares needed for golang web applications.

## JWT-RSA
web apis that are using rsa-jwt token for security can instrument the middleware to parse and validate the token.

```go
import (
    "github.com/BigBasket/go-web-middleware/auth/jwt"
    "github.com/gorilla/mux"
    "github.com/lestrrat-go/jwx/jwa"
)

func testHandler (w http.ResponseWriter, r *http.Request) {
    w.Write(http.StatusOK)
}

func main() {
    router := mux.NewRouter()
    jwtDecode := &RSADecode{
        PublicKey : "key" // base64 encoded rsa public key
        Algorithm: jwa.RS256
    }
    router.Use(jwtDecode.InstrumentJWTBearerToken)
    router.HandleFunc("/order/1", testHandler).Methods(http.MethodGet)
}
```

## NewRelic
```go
import (
    "github.com/BigBasket/go-web-middleware/profiling/newrelic"
    log "github.com/sirupsen/logrus"
)

func testHandler (w http.ResponseWriter, r *http.Request) {
    tx := newrelic.FetchTransaction(r.Context())
    // use the newrelic transaction for creating segments etc..
    w.Write(http.StatusOK)
}

func main() {
    router := mux.NewRouter()
    nConn, err := NewConnection(ctx, ConnectionConfig{
		newrelic.Config{AppName: "newrelic-test",
			License: "testkey"},
	})
	if err != nil {
        // handle errors
		return
    }
    router.Use(nConn.InstrumentNewRelic)
    router.HandleFunc("/order/1", testHandler).Methods(http.MethodGet)
}
```

## XTracker

```go
import (
    "github.com/BigBasket/go-web-middleware/profiling/tracker"
)

func testHandler (w http.ResponseWriter, r *http.Request) {
    tracker := FetchXTracker(r.Context())
    w.Write(http.StatusOK)
}

func main() {
    router := mux.NewRouter()
    router.Use(profiling.InstrumentXTracker)
    router.HandleFunc("/order/1", testHandler).Methods(http.MethodGet)
}
```