package sync

type PartyInfoRequest struct {
SenderKey string `json:"key"`
}

type PartyInfoResponse struct {
PublicKeys []ProvenPublicKey `json:"publicKeys"`
}

type ProvenPublicKey struct {
Key string `json:"key"`
Proof string `json:"proof"`
}
