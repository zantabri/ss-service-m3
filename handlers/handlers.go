package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zantabri/ss-service/store"
)

type JsonError struct {
	Message string `json:"message"`
}

type Secret struct {
	Secret string
	Seen   bool
}

type AddSecretRequest struct {
	Secret string `json:"plain_text"`
}

type AddSecretResponse struct {
	Id string `json:"id"`
}

type GetSecretResponse struct {
	Secret string `json:"data"`
}

type Handler struct {
	store *store.SecretStore
}

func New(store *store.SecretStore) (handle Handler) {

	handle = Handler{store: store}
	return

}

func (handler Handler) HealthCheck(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	writer.Write([]byte("ok"))

}

func (handler Handler) GetSecret(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	id := request.URL.Query().Get("id")

	if len(id) == 0 {

		writer.WriteHeader(400)
		raw, _ := json.Marshal(JsonError{Message: "invalid or missing id"})
		fmt.Fprintf(writer, "%s", raw)
		return

	}

	var store store.SecretStore = *handler.store

	resp := GetSecretResponse{Secret: store.RetriveSecret(id)}
	raw, err := json.Marshal(resp)

	if err != nil {

		writer.WriteHeader(500)
		raw, _ := json.Marshal(JsonError{Message: err.Error()})
		fmt.Fprintf(writer, "%s", string(raw))
		return

	}

	writer.WriteHeader(200)
	fmt.Fprintf(writer, "%s", string(raw))

}

func (handler Handler) AddSecret(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {

	bodyArr := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bodyArr)

	if err != nil && err != io.EOF {
		writer.WriteHeader(400)
		raw, _ := json.Marshal(JsonError{Message: err.Error()})
		fmt.Fprintf(writer, "%s", string(raw))
		return
	}

	payload := AddSecretRequest{}
	err = json.Unmarshal(bodyArr, &payload)

	if err != nil && err != io.EOF {
		writer.WriteHeader(400)
		raw, _ := json.Marshal(JsonError{Message: err.Error()})
		fmt.Fprintf(writer, "%s", string(raw))
		return
	}

	var store = *handler.store
	id := store.StoreSecret(payload.Secret)

	resp, err := json.Marshal(AddSecretResponse{Id: id})

	if err != nil {

		writer.WriteHeader(500)
		raw, _ := json.Marshal(JsonError{Message: err.Error()})
		fmt.Fprintf(writer, "%s", string(raw))
		return
	}

	writer.WriteHeader(201)
	fmt.Fprintf(writer, "%s", string(resp))

}
