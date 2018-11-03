package api


type DeleteRequest struct {
	Key string `json:"key"`
}

type ResendRequest struct {
	Key string `json:"key"`
}

type SendRequest struct {
	Key string `json:"key"`
}

type ReceiveRequest struct {
	Key string `json:"key"`
}