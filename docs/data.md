

# data
`import "Smilo-blackbox/src/data"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func SetFilename(filename string)](#SetFilename)
* [func SetLogger(loggers *logrus.Entry)](#SetLogger)
* [func Start()](#Start)
* [type EncryptedRawTransaction](#EncryptedRawTransaction)
  * [func FindEncryptedRawTransaction(hash []byte) (*EncryptedRawTransaction, error)](#FindEncryptedRawTransaction)
  * [func NewEncryptedRawTransaction(encodedPayload []byte, sender []byte) *EncryptedRawTransaction](#NewEncryptedRawTransaction)
  * [func (et *EncryptedRawTransaction) Delete() error](#EncryptedRawTransaction.Delete)
  * [func (et *EncryptedRawTransaction) Save() error](#EncryptedRawTransaction.Save)
* [type EncryptedTransaction](#EncryptedTransaction)
  * [func CreateEncryptedTransaction(hash []byte, encodedPayload []byte, timestamp time.Time) *EncryptedTransaction](#CreateEncryptedTransaction)
  * [func FindEncryptedTransaction(hash []byte) (*EncryptedTransaction, error)](#FindEncryptedTransaction)
  * [func NewEncryptedTransaction(encodedPayload []byte) *EncryptedTransaction](#NewEncryptedTransaction)
  * [func (et *EncryptedTransaction) Delete() error](#EncryptedTransaction.Delete)
  * [func (et *EncryptedTransaction) Save() error](#EncryptedTransaction.Save)
* [type Peer](#Peer)
  * [func FindPeer(publicKey []byte) (*Peer, error)](#FindPeer)
  * [func NewPeer(pKey []byte, nodeURL string) *Peer](#NewPeer)
  * [func Update(pKey []byte, nodeURL string) *Peer](#Update)
  * [func (p *Peer) Delete() error](#Peer.Delete)
  * [func (p *Peer) Save() error](#Peer.Save)


#### <a name="pkg-files">Package files</a>
[db.go](/src/Smilo-blackbox/src/data/db.go) [encrypted_raw_transaction.go](/src/Smilo-blackbox/src/data/encrypted_raw_transaction.go) [encrypted_transaction.go](/src/Smilo-blackbox/src/data/encrypted_transaction.go) [log.go](/src/Smilo-blackbox/src/data/log.go) [peer.go](/src/Smilo-blackbox/src/data/peer.go) 





## <a name="SetFilename">func</a> [SetFilename](/src/target/db.go?s=933:966#L30)
``` go
func SetFilename(filename string)
```
SetFilename set filename



## <a name="SetLogger">func</a> [SetLogger](/src/target/log.go?s=1019:1056#L30)
``` go
func SetLogger(loggers *logrus.Entry)
```
SetLogger set the logger



## <a name="Start">func</a> [Start](/src/target/db.go?s=1018:1030#L35)
``` go
func Start()
```
Start will start the db




## <a name="EncryptedRawTransaction">type</a> [EncryptedRawTransaction](/src/target/encrypted_raw_transaction.go?s=79:242#L6)
``` go
type EncryptedRawTransaction struct {
    Hash           []byte `storm:"id"`
    EncodedPayload []byte
    Sender         []byte
    Timestamp      time.Time `storm:"index"`
}

```
EncryptedRawTransaction holds hash and payload







### <a name="FindEncryptedRawTransaction">func</a> [FindEncryptedRawTransaction](/src/target/encrypted_raw_transaction.go?s=712:791#L25)
``` go
func FindEncryptedRawTransaction(hash []byte) (*EncryptedRawTransaction, error)
```
FindEncryptedRawTransaction will find a encrypted transaction for a hash


### <a name="NewEncryptedRawTransaction">func</a> [NewEncryptedRawTransaction](/src/target/encrypted_raw_transaction.go?s=344:438#L14)
``` go
func NewEncryptedRawTransaction(encodedPayload []byte, sender []byte) *EncryptedRawTransaction
```
NewEncryptedRawTransaction will create a new encrypted transaction based on the provided payload





### <a name="EncryptedRawTransaction.Delete">func</a> (\*EncryptedRawTransaction) [Delete](/src/target/encrypted_raw_transaction.go?s=1081:1130#L41)
``` go
func (et *EncryptedRawTransaction) Delete() error
```
Delete delete it on the db




### <a name="EncryptedRawTransaction.Save">func</a> (\*EncryptedRawTransaction) [Save](/src/target/encrypted_raw_transaction.go?s=979:1026#L36)
``` go
func (et *EncryptedRawTransaction) Save() error
```
Save saves into db




## <a name="EncryptedTransaction">type</a> [EncryptedTransaction](/src/target/encrypted_transaction.go?s=918:1055#L26)
``` go
type EncryptedTransaction struct {
    Hash           []byte `storm:"id"`
    EncodedPayload []byte
    Timestamp      time.Time `storm:"index"`
}

```
EncryptedTransaction holds hash and payload







### <a name="CreateEncryptedTransaction">func</a> [CreateEncryptedTransaction](/src/target/encrypted_transaction.go?s=1560:1670#L48)
``` go
func CreateEncryptedTransaction(hash []byte, encodedPayload []byte, timestamp time.Time) *EncryptedTransaction
```
CreateEncryptedTransaction will encrypt the transaction


### <a name="FindEncryptedTransaction">func</a> [FindEncryptedTransaction](/src/target/encrypted_transaction.go?s=1886:1959#L58)
``` go
func FindEncryptedTransaction(hash []byte) (*EncryptedTransaction, error)
```
FindEncryptedTransaction will find a encrypted transaction for a hash


### <a name="NewEncryptedTransaction">func</a> [NewEncryptedTransaction](/src/target/encrypted_transaction.go?s=1154:1227#L33)
``` go
func NewEncryptedTransaction(encodedPayload []byte) *EncryptedTransaction
```
NewEncryptedTransaction will create a new encrypted transaction based on the provided payload





### <a name="EncryptedTransaction.Delete">func</a> (\*EncryptedTransaction) [Delete](/src/target/encrypted_transaction.go?s=2243:2289#L74)
``` go
func (et *EncryptedTransaction) Delete() error
```
Delete delete it on the db




### <a name="EncryptedTransaction.Save">func</a> (\*EncryptedTransaction) [Save](/src/target/encrypted_transaction.go?s=2144:2188#L69)
``` go
func (et *EncryptedTransaction) Save() error
```
Save saves into db




## <a name="Peer">type</a> [Peer](/src/target/peer.go?s=858:927#L20)
``` go
type Peer struct {
    // contains filtered or unexported fields
}

```
Peer holds url and pub for a peer







### <a name="FindPeer">func</a> [FindPeer](/src/target/peer.go?s=1416:1462#L47)
``` go
func FindPeer(publicKey []byte) (*Peer, error)
```
FindPeer will find a peer


### <a name="NewPeer">func</a> [NewPeer](/src/target/peer.go?s=975:1022#L26)
``` go
func NewPeer(pKey []byte, nodeURL string) *Peer
```
NewPeer create new peer based on pk and url


### <a name="Update">func</a> [Update](/src/target/peer.go?s=1109:1155#L32)
``` go
func Update(pKey []byte, nodeURL string) *Peer
```
Update will update a peer





### <a name="Peer.Delete">func</a> (\*Peer) [Delete](/src/target/peer.go?s=1721:1750#L63)
``` go
func (p *Peer) Delete() error
```
Delete delete a peer on db




### <a name="Peer.Save">func</a> (\*Peer) [Save](/src/target/peer.go?s=1640:1667#L58)
``` go
func (p *Peer) Save() error
```
Save save a peer into db








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
