package server

import (
	"net/http"
	"testing"

	"Smilo-blackbox/src/server/api"

	"bytes"
	"encoding/json"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/require"
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
			contentType      string
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

				if test.followUp {

					var err error
					var sendRequest api.SendRequest
					var sendResponse api.SendResponse

					err = json.Unmarshal([]byte(test.body), &sendRequest)
					require.Empty(t, err)

					err = json.Unmarshal(response.RawBody, &sendResponse)
					require.Empty(t, err)

					receiveRequest := api.ReceiveRequest{Key: sendResponse.Key, To: sendRequest.To[0]}

					targetObject, err := json.Marshal(receiveRequest)
					require.Empty(t, err)

					var followUpResponse *testflight.Response

					if test.followUpMethod == "GET+BODY" {

						return

						//TODO: Fix this test
						targetBody := string(targetObject)

						newrequest, err := http.NewRequest("GET", test.followUpEndpoint, bytes.NewBuffer([]byte(targetBody)))
						newrequest.Header.Set("Content-Type", "application/json")

						require.Empty(t, err)
						require.NotEmpty(t, newrequest)

						followUpResponse = r.Do(newrequest)
					}

					require.NotEmpty(t, followUpResponse)
					require.NotEmpty(t, followUpResponse.StatusCode)
					require.Equal(t, sendRequest.Payload, followUpResponse.Body)
				}

			})
		}

	})

}
