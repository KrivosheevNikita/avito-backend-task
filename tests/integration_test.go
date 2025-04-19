package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

func getToken(t *testing.T, role string) string {
	body := map[string]string{"role": role}
	data, _ := json.Marshal(body)

	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewReader(data))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var tokenResp struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	assert.NoError(t, err)

	return tokenResp.Token
}

func sendRequest(t *testing.T, method, url, token string, body any) (*http.Response, []byte) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		assert.NoError(t, err)
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()

	return resp, respBody
}

func TestIntegration(t *testing.T) {
	modToken := getToken(t, "moderator")
	empToken := getToken(t, "employee")

	pvzBody := map[string]string{
		"city": "Москва",
	}
	resp, body := sendRequest(t, "POST", baseURL+"/pvz", modToken, pvzBody)
	assert.Equal(t, 201, resp.StatusCode)

	var pvz struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal(body, &pvz)
	assert.NoError(t, err)

	receptionBody := map[string]string{
		"pvzId": pvz.ID,
	}
	resp, body = sendRequest(t, "POST", baseURL+"/receptions", empToken, receptionBody)
	assert.Equal(t, 201, resp.StatusCode)

	var reception struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(body, &reception)
	assert.NoError(t, err)

	for i := 0; i < 50; i++ {
		productBody := map[string]string{
			"type":  "электроника",
			"pvzId": pvz.ID,
		}
		resp, _ := sendRequest(t, "POST", baseURL+"/products", empToken, productBody)
		assert.Equal(t, 201, resp.StatusCode)
	}

	closeURL := fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvz.ID)
	resp, _ = sendRequest(t, "POST", closeURL, empToken, nil)
	assert.Equal(t, 200, resp.StatusCode)
}
