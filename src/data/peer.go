package data

type Peer struct {
	publicKey []byte `storm:"id"`
	url       string
}

func NewPeer(pKey []byte, nodeURL string) *Peer {
	p := Peer{publicKey: pKey, url: nodeURL}
	return &p
}

func Update(pKey []byte, nodeURL string) *Peer {
	p, err := FindPeer(pKey)
	if err != nil {
		p = NewPeer(pKey, nodeURL)
	} else {
		p.url = nodeURL
	}
	p.Save()
	return p
}

func FindPeer(publicKey []byte) (*Peer, error) {
	var p Peer
	err := db.One("publicKey", publicKey, &p)
	if err != nil {
		log.Error("Unable to find Peer.")
		return nil, err
	}
	return &p, nil
}

func (p *Peer) Save() error {
	return db.Save(p)
}

func (p *Peer) Delete() error {
	return db.DeleteStruct(p)
}
