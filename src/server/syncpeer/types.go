package syncpeer

import (
	"io"
	"net/http"
)

//PartyInfoRequest used to marshal/unmarshal json
type PartyInfoRequest struct {
	SenderURL   string `json:"url"`
	SenderKey   string `json:"key"`
	SenderNonce string `json:"nonce"`
}

//PartyInfoResponse used to marshal/unmarshal json
type PartyInfoResponse struct {
	PublicKeys []ProvenPublicKey `json:"publicKeys"`
	PeerURLs   []string          `json:"peers"`
}

//ProvenPublicKey used to marshal/unmarshal json
type ProvenPublicKey struct {
	Key   string `json:"key"`
	Proof string `json:"proof"`
}

type HTTPClientWrapper struct {
	http.Client
	RequestResponseFunction func(req *http.Request) (*http.Response, error)
	PostResponseFunction    func(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

func (h *HTTPClientWrapper) Do(req *http.Request) (*http.Response, error) {
	return h.RequestResponseFunction(req)
}

func (h *HTTPClientWrapper) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return h.PostResponseFunction(url, contentType, body)
}
