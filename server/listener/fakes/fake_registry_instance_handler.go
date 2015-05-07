package fakes

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type agentSettingsResponse struct {
	Settings string `json:"settings"`
	Status   string `json:"status"`
}

type FakeRegistryInstanceHandler struct {
	Username         string
	Password         string
	InstanceSettings []byte
}

func NewFakeRegistryInstanceHandler(username string, password string) *FakeRegistryInstanceHandler {
	return &FakeRegistryInstanceHandler{
		Username:         username,
		Password:         password,
		InstanceSettings: []byte{},
	}
}

func (s *FakeRegistryInstanceHandler) HandleFunc(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		if s.InstanceSettings != nil {
			response := agentSettingsResponse{
				Settings: string(s.InstanceSettings),
				Status:   "ok",
			}
			responseJSON, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Error marshalling response", http.StatusBadRequest)
				return
			}
			w.Write(responseJSON)
			return
		}

		http.NotFound(w, req)
		return
	}

	if req.Method == "PUT" {
		if !s.isAuthorized(req) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		reqBody, _ := ioutil.ReadAll(req.Body)
		s.InstanceSettings = reqBody

		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method == "DELETE" {
		if !s.isAuthorized(req) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		s.InstanceSettings = []byte{}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func (s *FakeRegistryInstanceHandler) isAuthorized(req *http.Request) bool {
	if s.Username != "" && s.Password != "" {
		auth := s.Username + ":" + s.Password
		expectedAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

		return expectedAuthorizationHeader == req.Header.Get("Authorization")
	}

	return true
}
