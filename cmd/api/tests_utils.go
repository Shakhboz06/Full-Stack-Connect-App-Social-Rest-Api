package main

import (
	"go-project/internal/auth"
	"go-project/internal/ratelimiter"
	"go-project/internal/store"
	"go-project/internal/store/cache"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg servConfig) *application {
	t.Helper()

	// logger := zap.NewNop().Sugar()
	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStorage := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	rateLimiter := ratelimiter.NewfixedWindowRateLimiter(
		cfg.ratelimiter.RequestPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	return &application{
		logger: logger,
		store: mockStore,
		cacheStorage: mockCacheStorage,
		authenticator: testAuth,
		ratelimiter: rateLimiter,
	}
}

func executor(req *http.Request, mux http.Handler) *httptest.ResponseRecorder{
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}



func checkResponseCode(t *testing.T, expected, actual int){
	if expected != actual{
		t.Errorf("Expected response code %d, obtained %d", expected, actual)
	}
}