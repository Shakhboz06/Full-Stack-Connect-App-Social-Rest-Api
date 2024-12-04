package main

import (
	"log"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetUser(t *testing.T) {
	var config servConfig
	app := newTestApplication(t,config)
	mux := app.mount().(*chi.Mux)

	testToken, err :=  app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("should not allow unathenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil{
			t.Fatal(err)

		}
		rr := executor(req, mux)

		if rr.Code != http.StatusUnauthorized{
			t.Errorf("Expected response code to be %d and we got %d", http.StatusUnauthorized, rr.Code)
		}


	})


	t.Run("should allow authenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil{
			t.Fatal(err)

		}
		rr := executor(req, mux)


		req.Header.Set("Authorization", "Bearer " + testToken)
		checkResponseCode(t, http.StatusOK, rr.Code)

		log.Println(rr.Body)
	})
}