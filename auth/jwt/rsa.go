package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

const (
	//BearerAuthHeader - authorization header name for bearer token
	BearerAuthHeader string = "Authorization"

	//BearerExpiry - key for token expiry
	BearerExpiry string = "expiry"
)

//EncodingType - types of encoding supported
type EncodingType string

const (
	//Base64 - base64 encoding/decoding will be used
	Base64 EncodingType = "base64"
)

// RSADecode - object which implements the jwt token
type RSADecode struct {
	PublicKey             string
	Algorithm             jwa.SignatureAlgorithm
	PublicKeyEncodingType EncodingType
}

//GetPKCS1PublicKeyFromPEMBytes - reads the public key from bytes
func GetPKCS1PublicKeyFromPEMBytes(inp []byte) (*rsa.PublicKey, error) {
	pemBlock, _ := pem.Decode(inp)
	publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err == nil {
		return publicKey.(*rsa.PublicKey), nil
	}
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}

//ParseBase64EncodedRSAPublicKey - parse base64 encoded rsa public key
func ParseBase64EncodedRSAPublicKey(key string) (*rsa.PublicKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return GetPKCS1PublicKeyFromPEMBytes(bytes)
}

// InstrumentJWTBearerToken - sets the jwt claims to the context headers.
func (decode RSADecode) InstrumentJWTBearerToken(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(BearerAuthHeader)
		if strings.TrimSpace(authHeader) == "" { //authorization header is empty
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("authorization header is empty"))
			return
		}
		headers := strings.Fields(authHeader) //auth header has two fields; 1. type of auth (bearer) 2. the token (token)
		if len(headers) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("authorization header length is invalid"))
			return
		}

		if strings.ToLower(headers[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("authorization header type is not bearer"))
			return
		}

		pKey, err := ParseBase64EncodedRSAPublicKey(decode.PublicKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			msg := fmt.Sprintf("public key decoding failed; err:%v", err.Error())
			w.Write([]byte(msg))
			return
		}
		token, err := jwt.Parse(strings.NewReader(headers[1]), jwt.WithVerify(decode.Algorithm, pKey))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			msg := fmt.Sprintf("jwt token verfication failed err:%v", err.Error())
			w.Write([]byte(msg))
			return
		}

		expiry, exists := token.Get(BearerExpiry)
		if exists == false { //token expiry claim doesn't exists.
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token expiry claim doesn't exists"))
			return
		}
		exp, err := time.Parse(time.RFC3339, expiry.(string))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token expiry is not of RFC3339 time format"))
			return
		}
		currTime := time.Now().UTC()     //checking on UTC timezone; ensure the client is using the same timezone.
		if currTime.After(exp) == true { // currTime passes token expiry;
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token is expired"))
			return
		}
		next.ServeHTTP(w, r)
	}
	return fn
}
