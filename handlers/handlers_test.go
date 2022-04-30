package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/zantabri/ss-service/store"
)


func TestMain(m *testing.M) {

	exitCode := m.Run()
	os.Exit(exitCode)

}

type MockStore struct {

	StoreSecretCall int
	RetriveSecretCall int

}

func (store *MockStore) StoreSecret(key string) string {

	id := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	store.StoreSecretCall++
	return id

}

func (store *MockStore) RetriveSecret(id string) string {

	store.RetriveSecretCall++
	return "val"

}

func TestHealthCheck(t *testing.T) {

	mux := httprouter.New()
	var store store.SecretStore = &MockStore{}
	handlers := New(&store)
	mux.GET("/health", handlers.HealthCheck)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/health", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %d", writer.Code)
	}

	resp := writer.Body.String()
	if resp != "ok" {
		t.Errorf("response body is %s", resp)
	}

}

func TestGetSecret(t *testing.T) {

	mux := httprouter.New()
	var store store.SecretStore = &MockStore{}

	Handlers := New(&store)

	mux.GET("/",Handlers.GetSecret)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET","/?id=234", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("response is %d", writer.Code)
	}

	resp := writer.Body.String()
	t.Logf("response is %s", resp)

	mock, _ := store.(*MockStore)

	if mock.RetriveSecretCall != 1 {
		t.Error("retriever secret not called")
	}

}

func TestGetSecretWhereIdNotSpecified(t *testing.T) {

	mux := httprouter.New()
	var store store.SecretStore = &MockStore{}

	Handlers := New(&store)

	mux.GET("/",Handlers.GetSecret)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET","/", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 400 {
		t.Fatalf("response is %d", writer.Code)
	}

	resp := writer.Body.String()
	t.Logf("response is %s", resp)

	mock, _ := store.(*MockStore)

	if mock.RetriveSecretCall != 0 {
		t.Fatal("retriever secret called")
	}

}

func TestAddSecret(t *testing.T) {
	mux := httprouter.New()
	var store store.SecretStore = &MockStore{}

	Handlers := New(&store)

	mux.POST("/",Handlers.AddSecret)
	writer := httptest.NewRecorder()
	payload := "{\"plain_text\" : \"lubalu2323232balu\"}"
	
	request, _ := http.NewRequest("POST","/", strings.NewReader(payload))
	mux.ServeHTTP(writer, request)

	if writer.Code != 201 {
		t.Errorf("response is %d", writer.Code)
	}

	resp := writer.Body.String()
	t.Logf("response is %s", resp)

	var responseBody AddSecretResponse
	json.Unmarshal(writer.Body.Bytes(), &responseBody)

	if responseBody.Id != "1f2b78f8b8067dfd47df852f12697c69" {
		t.Errorf("returned id is %s", responseBody.Id)
	}

	mock, _ := store.(*MockStore)

	if mock.StoreSecretCall != 1 {
		t.Errorf("store secret call is %d", mock.StoreSecretCall)
	}

}
