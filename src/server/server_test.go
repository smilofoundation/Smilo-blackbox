// Copyright 2019 The Smilo-blackbox Authors
// This file is part of the Smilo-blackbox library.
//
// The Smilo-blackbox library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Smilo-blackbox library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Smilo-blackbox library. If not, see <http://www.gnu.org/licenses/>.

package server

import (
	"fmt"
	"net/http"
	"testing"

	"Smilo-blackbox/src/server/api"

	"bytes"
	"encoding/json"

	"encoding/base64"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server/encoding"
	"Smilo-blackbox/src/server/syncpeer"
	"Smilo-blackbox/src/utils"
)

var testEncryptedTransaction = createEncryptedTransaction()
var nonce = make([]byte, 24)

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
				response:    utils.BlackBoxVersion,
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test push",
				endpoint:    "/push",
				method:      "POST",
				body:        base64.StdEncoding.EncodeToString(testEncryptedTransaction.EncodedPayload),
				contentType: "application/octet-stream",
				response:    base64.StdEncoding.EncodeToString(testEncryptedTransaction.Hash),
				statusCode:  201,
				expectedErr: nil,
			},
			{
				name:        "test storeraw",
				endpoint:    "/storeraw",
				method:      "POST",
				body:        `{"payload":"MTIzNDU2Nzg5MGFiY2RlZmdoaWprbG1ub3BxcnM=","from":"MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk="}`,
				contentType: "application/json",
				response:    "{\"key\":\"gGW65vu3FhCdiMvCHEUoWfTllYeKP4HhVCZaC22Fl+192LnA6C1Rt4T9aGEHWSHUc+5fHwnX6KtjXaHM3h9nBw==\"}\n",
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test resend individual",
				endpoint:    "/resend",
				method:      "POST",
				body:        "{ \"type\": \"Individual\", \"publicKey\": \"" + base64.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012")) + "\", \"key\": \"" + base64.StdEncoding.EncodeToString(testEncryptedTransaction.Hash) + "\" }",
				contentType: "application/json",
				response:    base64.StdEncoding.EncodeToString(testEncryptedTransaction.EncodedPayload),
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test transaction delete",
				endpoint:    "/transaction/" + base64.URLEncoding.EncodeToString(createEncryptedTransactionForTest().Hash),
				method:      "DELETE",
				contentType: "application/json",
				response:    "",
				statusCode:  204,
				expectedErr: nil,
			},
			{
				name:        "test party info",
				endpoint:    "/partyinfo",
				method:      "POST",
				body:        "{ \"url\": \"http://localhost:9000\", \"key\": \"MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk=\", \"nonce\": \"" + base64.StdEncoding.EncodeToString(nonce) + "\" }",
				contentType: "application/json",
				response:    "",
				statusCode:  200,
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
					//fmt.Println(test.endpoint, string(createEncryptedTransactionForDeletion().Hash))
					response = r.Delete(test.endpoint, test.contentType, test.body)
				}

				if test.response != "" {
					require.NotEmpty(t, response)
					require.NotEmpty(t, response.StatusCode)
					require.NotEmpty(t, response.RawBody)
					if test.endpoint == "/storeraw" {
						var respJSON api.KeyJSON
						err := json.Unmarshal([]byte(response.Body), &respJSON)
						if err != nil {
							t.Logf("Invalid json response. %v", response)
							t.Fail()
						}
						key, err := base64.StdEncoding.DecodeString(respJSON.Key)
						if err != nil {
							t.Logf("Cannot decode key from json. %v", response)
							t.Fail()
						}
						encRawTrans, err := data.FindEncryptedRawTransaction(key)
						if err != nil {
							t.Logf("Raw transaction not found. %v", response)
							t.Fail()
						}
						require.Equal(t, encRawTrans.Hash, key)
					} else {
						require.Equal(t, test.response, response.Body)
					}
				}

				require.Equal(t, test.statusCode, response.StatusCode)

				if test.endpoint == "/partyinfo" {
					var respJSON syncpeer.PartyInfoResponse
					err := json.Unmarshal([]byte(response.Body), &respJSON)
					if err == nil {
						t.Logf("Public Key: %s Proof: %s", respJSON.PublicKeys[0].Key, respJSON.PublicKeys[0].Proof)
						require.Equal(t, respJSON.PublicKeys[0].Key, "MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk=")
						pubKey, _ := base64.StdEncoding.DecodeString("MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk=")
						proof, _ := base64.StdEncoding.DecodeString(respJSON.PublicKeys[0].Proof)
						ret := crypt.DecryptPayload(crypt.ComputeSharedKey(crypt.GetPrivateKey(pubKey), pubKey), proof, nonce)
						require.NotEmpty(t, ret)
						t.Logf("Unboxed Proof: %s", ret)
					} else {
						t.Logf("Invalid json response. %v", response)
						t.Fail()
					}
				}
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
				response:    utils.BlackBoxVersion,
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:        "test delete",
				endpoint:    "/delete",
				method:      "POST",
				contentType: "application/json",
				body:        `{"key": "` + base64.StdEncoding.EncodeToString(createEncryptedTransactionForTest().Hash) + `"}`,
				response:    "Delete successful",
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:     "test send signed tx",
				endpoint: "/sendsignedtx",
				method:   "CUSTOM",
				body:     string([]byte(base64.StdEncoding.EncodeToString(createEncryptedRawTransactionForTest().Hash))),
				headers: http.Header{
					utils.HeaderTo: []string{"OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="}},
				response:    "",
				statusCode:  200,
				expectedErr: nil,
			},
			{
				name:     "test send signed tx",
				endpoint: "/sendsignedtx",
				method:   "CUSTOM",
				body:     string([]byte(base64.StdEncoding.EncodeToString(createEncryptedRawTransactionForTest().Hash))),
				headers: http.Header{
					utils.HeaderTo: []string{"OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="}},
				response:    "",
				statusCode:  200,
				expectedErr: nil,
				followUp:         false,
				followUpEndpoint: "/transaction",
				followUpMethod:   "GET",
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
				headers: http.Header{utils.HeaderFrom: []string{"MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk="},
					utils.HeaderTo: []string{"OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="}},
				response:    "",
				statusCode:  200,
				expectedErr: nil,

				followUp:         false,
				followUpEndpoint: "/transaction",
				followUpMethod:   "GET",
			},

			{
				name:     "test send raw go-smilo test payload",
				endpoint: "/sendraw",
				method:   "CUSTOM",
				body:     string([]byte(base64.StdEncoding.EncodeToString([]byte("`\x80`@R4\x80\x15a\x00\x10W`\x00\x80\xfd[P`@Q` \x80a\x01a\x839\x81\x01\x80`@R\x81\x01\x90\x80\x80Q\x90` \x01\x90\x92\x91\x90PPP\x80`\x00\x81\x90UPPa\x01\x17\x80a\x00J`\x009`\x00\xf3\x00`\x80`@R`\x046\x10`SW`\x005|\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x90\x04c\xff\xff\xff\xff\x16\x80c*\x1a\xfc\xd9\x14`XW\x80c`\xfeG\xb1\x14`\x80W\x80cmL\xe6<\x14`\xaaW[`\x00\x80\xfd[4\x80\x15`cW`\x00\x80\xfd[P`j`\xd2V[`@Q\x80\x82\x81R` \x01\x91PP`@Q\x80\x91\x03\x90\xf3[4\x80\x15`\x8bW`\x00\x80\xfd[P`\xa8`\x04\x806\x03\x81\x01\x90\x80\x805\x90` \x01\x90\x92\x91\x90PPP`\xd8V[\x00[4\x80\x15`\xb5W`\x00\x80\xfd[P`\xbc`\xe2V[`@Q\x80\x82\x81R` \x01\x91PP`@Q\x80\x91\x03\x90\xf3[`\x00T\x81V[\x80`\x00\x81\x90UPPV[`\x00\x80T\x90P\x90V\x00\xa1ebzzr0X q\xec\xf8MD\xfa_μ\xb7\xb7S\xaa\x9avy\xb7\xbe\x04\xb5\xeb\xb89\xdbp\xc8$_G\xbf\xfc\x9c\x00)\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\xa4")))),
				//body:     string([]byte("`\x80`@R4\x80\x15a\x00\x10W`\x00\x80\xfd[P`@Q` \x80a\x01a\x839\x81\x01\x80`@R\x81\x01\x90\x80\x80Q\x90` \x01\x90\x92\x91\x90PPP\x80`\x00\x81\x90UPPa\x01\x17\x80a\x00J`\x009`\x00\xf3\x00`\x80`@R`\x046\x10`SW`\x005|\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x90\x04c\xff\xff\xff\xff\x16\x80c*\x1a\xfc\xd9\x14`XW\x80c`\xfeG\xb1\x14`\x80W\x80cmL\xe6<\x14`\xaaW[`\x00\x80\xfd[4\x80\x15`cW`\x00\x80\xfd[P`j`\xd2V[`@Q\x80\x82\x81R` \x01\x91PP`@Q\x80\x91\x03\x90\xf3[4\x80\x15`\x8bW`\x00\x80\xfd[P`\xa8`\x04\x806\x03\x81\x01\x90\x80\x805\x90` \x01\x90\x92\x91\x90PPP`\xd8V[\x00[4\x80\x15`\xb5W`\x00\x80\xfd[P`\xbc`\xe2V[`@Q\x80\x82\x81R` \x01\x91PP`@Q\x80\x91\x03\x90\xf3[`\x00T\x81V[\x80`\x00\x81\x90UPPV[`\x00\x80T\x90P\x90V\x00\xa1ebzzr0X q\xec\xf8MD\xfa_μ\xb7\xb7S\xaa\x9avy\xb7\xbe\x04\xb5\xeb\xb89\xdbp\xc8$_G\xbf\xfc\x9c\x00)\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\xa4")),

				headers:     http.Header{utils.HeaderTo: []string{"OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w="}},
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
				var sendResponse api.KeyJSON
				var followUpResponse *testflight.Response

				if test.followUpEndpoint == "/receive" {

					err = json.Unmarshal([]byte(test.body), &sendRequest)
					require.Empty(t, err)

					err = json.Unmarshal(response.RawBody, &sendResponse)
					require.NoError(t, err)

					require.Empty(t, err)

					receiveRequest := api.ReceiveRequest{Key: sendResponse.Key, To: sendRequest.To[0]}

					targetObject, err := json.Marshal(receiveRequest)
					require.Empty(t, err)

					targetBody := string(targetObject)

					newrequest, err := http.NewRequest("GET", test.followUpEndpoint, bytes.NewBuffer([]byte(targetBody)))
					require.NoError(t, err)
					newrequest.Header.Set("Content-Type", "application/json")

					require.Empty(t, err)
					require.NotEmpty(t, newrequest)

					followUpResponse = r.Do(newrequest)
					var responseJSON api.ReceiveResponse
					err = json.NewDecoder(bytes.NewBuffer(followUpResponse.RawBody)).Decode(&responseJSON)
					require.NoError(t, err)

					require.Equal(t, sendRequest.Payload, responseJSON.Payload)

				} else if test.followUpEndpoint == "/transaction" {

					key, err := base64.StdEncoding.DecodeString(response.Body)
					if err != nil {
						t.Fail()
					}
					urlEncodedKey := base64.URLEncoding.EncodeToString(key)
					t.Log("Send Response: ", response)
					toBytes, err := base64.StdEncoding.DecodeString("OeVDzTdR95fhLKIgpBLxqdDNXYzgozgi7dnnS125A3w=")
					if err != nil {
						t.Fail()
					}
					urlEncodedTo := base64.URLEncoding.EncodeToString(toBytes)

					targetURL := "/transaction/" + urlEncodedKey + "?to=" + urlEncodedTo
					followUpResponse = r.Get(targetURL)
					var responseJSON api.ReceiveResponse
					err = json.NewDecoder(bytes.NewBuffer(followUpResponse.RawBody)).Decode(&responseJSON)
					require.NoError(t, err)
					if test.endpoint == "/sendsignedtx" {
						decodedBody, _ := base64.StdEncoding.DecodeString(responseJSON.Payload)
						require.Equal(t, string(decodedBody), "12345")
					} else {
						require.Equal(t, test.body, responseJSON.Payload)
					}
				} else {
					return
				}
				require.NotEmpty(t, followUpResponse)
				require.NotEmpty(t, followUpResponse.StatusCode)

			})
		}

	})

}

func createEncryptedTransactionForTest() *data.EncryptedTransaction {
	encTrans := createEncryptedTransaction()
	err := encTrans.Save()
	if err != nil {
		fmt.Println("Could not createEncryptedTransactionForDeletion")
	}
	return encTrans
}

func createEncryptedTransaction() *data.EncryptedTransaction {
	toValues := make([][]byte, 1)
	toValues[0] = []byte("09876543210987654321098765432109")
	fromValue := []byte("12345678901234567890123456789012")
	payload, _ := crypt.NewRandomNonce()
	encPayloadData, _ := encoding.EncodePayloadData(payload, fromValue, toValues)
	encTrans := data.NewEncryptedTransaction(*encPayloadData.Serialize())
	return encTrans
}

func createEncryptedRawTransactionForTest() *data.EncryptedRawTransaction {
	encTrans := createEncryptedRawTransaction()
	err := encTrans.Save()
	if err != nil {
		fmt.Println("Could not createEncryptedTransactionForDeletion")
	}
	return encTrans
}

func createEncryptedRawTransaction() *data.EncryptedRawTransaction {
	toValues := make([][]byte, 1)
	pubkey, _ := base64.StdEncoding.DecodeString("MD3fapkkHUn86h/W7AUhiD4NiDFkuIxtuRr0Nge27Bk=")
	toValues[0] = pubkey
	fromValue := pubkey
	payload := []byte("12345")
	encPayloadData, _ := encoding.EncodePayloadData(payload, fromValue, toValues)
	encTrans := data.NewEncryptedRawTransaction(*encPayloadData.Serialize(),encPayloadData.Sender)
	return encTrans
}
