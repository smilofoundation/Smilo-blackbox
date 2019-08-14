package data

import "time"

// EncryptedTransaction holds hash and payload
type EncryptedRawTransaction struct {
	Hash           []byte `storm:"id"`
	EncodedPayload []byte
	Sender         []byte
	Timestamp      time.Time `storm:"index"`
}

// NewEncryptedTransaction will create a new encrypted transaction based on the provided payload
func NewEncryptedRawTransaction(encodedPayload []byte, sender []byte) *EncryptedRawTransaction {
	trans := EncryptedRawTransaction{
		Hash:           calculateHash(encodedPayload),
		EncodedPayload: encodedPayload,
		Sender:         sender,
		Timestamp:      time.Now(),
	}
	return &trans
}

// FindEncryptedTransaction will find a encrypted transaction for a hash
func FindEncryptedRawTransaction(hash []byte) (*EncryptedRawTransaction, error) {
	var t EncryptedRawTransaction
	err := db.One("Hash", hash, &t)
	if err != nil {
		log.Error("Unable to find transaction.")
		return nil, err
	}
	return &t, nil
}

//Save saves into db
func (et *EncryptedRawTransaction) Save() error {
	return db.Save(et)
}

//Delete delete it on the db
func (et *EncryptedRawTransaction) Delete() error {
	return db.DeleteStruct(et)
}
