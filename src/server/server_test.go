package server

import (
	"net/http"
	"testing"

	"Smilo-blackbox/src/server/api"

	"bytes"
	"encoding/json"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/require"
	"encoding/base64"
)

func TestPublicAPI(t *testing.T) {

	public, _ := InitRouting()

	testflight.WithServer(public, func(r *testflight.Requester) {

		testCases := []struct {
			name        string
			endpoint    string
			method      string
			body        string
			contentType string
			response    string
			statusCode  int
			expectedErr error
		}{
			{
				name:        "test upcheck",
				endpoint:    "/upcheck",
				method:      "GET",
				contentType: "application/json",
				response:    "I'm up!",
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test version",
				endpoint:    "/version",
				method:      "GET",
				contentType: "application/json",
				response:    api.BlackBoxVersion,
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test transaction delete",
				endpoint:    "/transaction/1",
				method:      "DELETE",
				contentType: "application/json",
				response:    "",
				statusCode:  204,
				expectedErr: nil,
			},
		}

		for _, test := range testCases {

			t.Run(test.name, func(t *testing.T) {

				var response *testflight.Response
				if test.method == "GET" {
					response = r.Get(test.endpoint)
				} else if test.method == "POST" {
					response = r.Post(test.endpoint, test.contentType, test.body)
				} else if test.method == "DELETE" {
					response = r.Delete(test.endpoint, test.contentType, test.body)

				}

				if test.response != "" {
					require.NotEmpty(t, response)
					require.NotEmpty(t, response.StatusCode)
					require.NotEmpty(t, response.RawBody)
					require.Equal(t, test.response, response.Body)
				}

				require.Equal(t, test.statusCode, response.StatusCode)

			})
		}

	})

}

func TestPrivateAPI(t *testing.T) {

	_, private := InitRouting()

	testflight.WithServer(private, func(r *testflight.Requester) {

		testCases := []struct {
			name             string
			endpoint         string
			method           string
			body             string
			bodyRaw          []byte
			contentType      string
			headers          http.Header
			response         string
			statusCode       int
			expectedErr      error
			followUp         bool
			followUpEndpoint string
			followUpMethod   string
		}{
			{
				name:        "test upcheck",
				endpoint:    "/upcheck",
				method:      "GET",
				contentType: "application/json",
				response:    "I'm up!",
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test version",
				endpoint:    "/version",
				method:      "GET",
				contentType: "application/json",
				response:    api.BlackBoxVersion,
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test delete",
				endpoint:    "/delete",
				method:      "POST",
				contentType: "application/json",
				body:        `{"key": "123456"}`,
				response:    "",
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:             "test send receive",
				endpoint:         "/send",
				method:           "POST",
				contentType:      "application/json",
				body:             `{"payload":"MTIzNDU2Nzg5MGFiY2RlZmdoaWprbG1ub3BxcnM=","from":"MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk=","to":["OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="]}`,
				response:         "",
				statusCode:       200,
				expectedErr:      nil,
				followUp:         true,
				followUpEndpoint: "/receive",
				followUpMethod:   "GET+BODY",
			},

			{
				name:     "test send raw & get transaction",
				endpoint: "/sendraw",
				method:   "CUSTOM",
				body:     string([]byte(base64.StdEncoding.EncodeToString([]byte("1234567890abcdefghijklmnopqrs")))),
				headers: http.Header{"c11n-from": []string{"MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk="},
					"c11n-to": []string{"OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="}},
				response:    "",
				statusCode:  200,
				expectedErr: nil,

				followUp:         false,
				followUpEndpoint: "/transaction",
				followUpMethod:   "GET",
			},
		}

		for _, test := range testCases {

			t.Run(test.name, func(t *testing.T) {

				var response *testflight.Response
				if test.method == "GET" {
					response = r.Get(test.endpoint)
				} else if test.method == "POST" {
					response = r.Post(test.endpoint, test.contentType, test.body)
				} else if test.method == "DELETE" {
					response = r.Delete(test.endpoint, test.contentType, test.body)
				} else if test.method == "CUSTOM" {

					newrequest, err := http.NewRequest("POST", test.endpoint, bytes.NewBuffer([]byte(test.body)))
					newrequest.Header = test.headers
					require.Empty(t, err)
					require.NotEmpty(t, newrequest)

					response = r.Do(newrequest)
				}

				require.NotEmpty(t, response)
				require.NotEmpty(t, response.StatusCode)
				require.NotEmpty(t, response.RawBody)
				if test.response != "" {
					require.Equal(t, test.response, response.Body)
				}

				require.Equal(t, test.statusCode, response.StatusCode)

				var err error
				var sendRequest api.SendRequest
				var sendResponse api.SendResponse
				var followUpResponse *testflight.Response

				if test.followUpEndpoint == "/receive" {

					t.SkipNow()

					err = json.Unmarshal([]byte(test.body), &sendRequest)
					require.Empty(t, err)

					err = json.Unmarshal(response.RawBody, &sendResponse)
					require.Empty(t, err)

					receiveRequest := api.ReceiveRequest{Key: sendResponse.Key, To: sendRequest.To[0]}

					targetObject, err := json.Marshal(receiveRequest)
					require.Empty(t, err)


					//TODO: Fix this test
					targetBody := string(targetObject)

					newrequest, err := http.NewRequest("GET", test.followUpEndpoint, bytes.NewBuffer([]byte(targetBody)))
					newrequest.Header.Set("Content-Type", "application/json")

					require.Empty(t, err)
					require.NotEmpty(t, newrequest)

					followUpResponse = r.Do(newrequest)

				} else if test.followUpEndpoint == "/transaction" {


					//TODO: Fix this test

					t.SkipNow()

					key, err := base64.StdEncoding.DecodeString(response.Body)
					if err != nil {
						t.Fail()
					}
					urlEncodedKey := base64.URLEncoding.EncodeToString(key)
					log.Debug("Send Response: ", response)
					toBytes, err := base64.StdEncoding.DecodeString("OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w")
					if err != nil {
						t.Fail()
					}
					urlEncodedTo := base64.URLEncoding.EncodeToString(toBytes)

					targetURL := "/transaction/" + urlEncodedKey + "?to=" + urlEncodedTo
					followUpResponse = r.Get(targetURL)

				} else {
					return
				}
				require.NotEmpty(t, followUpResponse)
				require.NotEmpty(t, followUpResponse.StatusCode)
				require.Equal(t, sendRequest.Payload, followUpResponse.Body)

			})
		}

	})

}
