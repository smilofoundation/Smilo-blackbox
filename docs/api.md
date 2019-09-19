

# api
`import "Smilo-blackbox/src/server/api"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func API(w http.ResponseWriter, r *http.Request)](#API)
* [func ConfigPeersGet(w http.ResponseWriter, r *http.Request)](#ConfigPeersGet)
* [func ConfigPeersPut(w http.ResponseWriter, r *http.Request)](#ConfigPeersPut)
* [func Delete(w http.ResponseWriter, r *http.Request)](#Delete)
* [func GetPartyInfo(w http.ResponseWriter, r *http.Request)](#GetPartyInfo)
* [func GetVersion(w http.ResponseWriter, r *http.Request)](#GetVersion)
* [func Metrics(w http.ResponseWriter, r *http.Request)](#Metrics)
* [func Push(w http.ResponseWriter, r *http.Request)](#Push)
* [func PushTransactionForOtherNodes(encryptedTransaction data.EncryptedTransaction, recipient []byte)](#PushTransactionForOtherNodes)
* [func Receive(w http.ResponseWriter, r *http.Request)](#Receive)
* [func ReceiveRaw(w http.ResponseWriter, r *http.Request)](#ReceiveRaw)
* [func Resend(w http.ResponseWriter, r *http.Request)](#Resend)
* [func RetrieveAndDecryptPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) []byte](#RetrieveAndDecryptPayload)
* [func RetrieveJSONPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte)](#RetrieveJSONPayload)
* [func Send(w http.ResponseWriter, r *http.Request)](#Send)
* [func SendRaw(w http.ResponseWriter, r *http.Request)](#SendRaw)
* [func SendSignedTx(w http.ResponseWriter, r *http.Request)](#SendSignedTx)
* [func SetLogger(loggers *logrus.Entry)](#SetLogger)
* [func StoreRaw(w http.ResponseWriter, r *http.Request)](#StoreRaw)
* [func TransactionDelete(w http.ResponseWriter, r *http.Request)](#TransactionDelete)
* [func TransactionGet(w http.ResponseWriter, r *http.Request)](#TransactionGet)
* [func UnknownRequest(w http.ResponseWriter, r *http.Request)](#UnknownRequest)
* [func Upcheck(w http.ResponseWriter, r *http.Request)](#Upcheck)
* [type KeyJSON](#KeyJSON)
* [type PeerURL](#PeerURL)
* [type ReceiveRequest](#ReceiveRequest)
  * [func (e *ReceiveRequest) Parse() ([]byte, []byte, []string)](#ReceiveRequest.Parse)
* [type ReceiveResponse](#ReceiveResponse)
* [type ResendRequest](#ResendRequest)
* [type SendRequest](#SendRequest)
  * [func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string)](#SendRequest.Parse)
* [type StoreRawRequest](#StoreRawRequest)


#### <a name="pkg-files">Package files</a>
[common.go](/src/Smilo-blackbox/src/server/api/common.go) [log.go](/src/Smilo-blackbox/src/server/api/log.go) [private.go](/src/Smilo-blackbox/src/server/api/private.go) [public.go](/src/Smilo-blackbox/src/server/api/public.go) [types.go](/src/Smilo-blackbox/src/server/api/types.go) 





## <a name="API">func</a> [API](/src/target/common.go?s=1600:1648#L51)
``` go
func API(w http.ResponseWriter, r *http.Request)
```
API Request path "/api", response json rest api spec.



## <a name="ConfigPeersGet">func</a> [ConfigPeersGet](/src/target/public.go?s=14078:14137#L420)
``` go
func ConfigPeersGet(w http.ResponseWriter, r *http.Request)
```
ConfigPeersGet Receive a GET request with index on path and return Status Code 200 and Peer json containing url, Status Code 404 if not found.



## <a name="ConfigPeersPut">func</a> [ConfigPeersPut](/src/target/public.go?s=13494:13553#L403)
``` go
func ConfigPeersPut(w http.ResponseWriter, r *http.Request)
```
ConfigPeersPut It receives a PUT request with a json containing a Peer url and returns Status Code 200.



## <a name="Delete">func</a> [Delete](/src/target/public.go?s=8777:8828#L254)
``` go
func Delete(w http.ResponseWriter, r *http.Request)
```
Delete Deprecated API
It receives a POST request with a json containing a DeleteRequest with key and returns Status 200 if succeed, 404 otherwise.



## <a name="GetPartyInfo">func</a> [GetPartyInfo](/src/target/public.go?s=1245:1302#L40)
``` go
func GetPartyInfo(w http.ResponseWriter, r *http.Request)
```
GetPartyInfo It receives a POST request with a json containing url and key, returns local publicKeys and a proof that private key is known.



## <a name="GetVersion">func</a> [GetVersion](/src/target/common.go?s=1115:1170#L35)
``` go
func GetVersion(w http.ResponseWriter, r *http.Request)
```
GetVersion Request path "/version", response plain text version ID



## <a name="Metrics">func</a> [Metrics](/src/target/public.go?s=15235:15287#L455)
``` go
func Metrics(w http.ResponseWriter, r *http.Request)
```
Metrics Receive a GET request and return Status Code 200 and server internal status information in plain text.



## <a name="Push">func</a> [Push](/src/target/public.go?s=3385:3434#L93)
``` go
func Push(w http.ResponseWriter, r *http.Request)
```
Push It receives a POST request with a payload and returns Status Code 201 with a payload generated hash, on error returns Status Code 500.



## <a name="PushTransactionForOtherNodes">func</a> [PushTransactionForOtherNodes](/src/target/common.go?s=3113:3212#L95)
``` go
func PushTransactionForOtherNodes(encryptedTransaction data.EncryptedTransaction, recipient []byte)
```
PushTransactionForOtherNodes will push encrypted transaction to other nodes



## <a name="Receive">func</a> [Receive](/src/target/private.go?s=9019:9071#L276)
``` go
func Receive(w http.ResponseWriter, r *http.Request)
```
Receive is a Deprecated API
It receives a ReceiveRequest json with an encoded key (hash) and to values, returns decrypted payload



## <a name="ReceiveRaw">func</a> [ReceiveRaw](/src/target/public.go?s=5040:5095#L150)
``` go
func ReceiveRaw(w http.ResponseWriter, r *http.Request)
```
ReceiveRaw Receive a GET request with header params bb0x-key and bb0x-to, return unencrypted payload



## <a name="Resend">func</a> [Resend](/src/target/public.go?s=6827:6878#L199)
``` go
func Resend(w http.ResponseWriter, r *http.Request)
```
Resend It receives a POST request with a json ResendRequest containing type (INDIVIDUAL, ALL), publicKey and key(for individual requests),
it returns encoded payload for INDIVIDUAL or it does one push request for each payload and returns empty for type ALL.



## <a name="RetrieveAndDecryptPayload">func</a> [RetrieveAndDecryptPayload](/src/target/common.go?s=2366:2466#L74)
``` go
func RetrieveAndDecryptPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) []byte
```
RetrieveAndDecryptPayload will retrieve and decrypt the payload



## <a name="RetrieveJSONPayload">func</a> [RetrieveJSONPayload](/src/target/common.go?s=1852:1939#L61)
``` go
func RetrieveJSONPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte)
```
RetrieveJSONPayload will retrieve payload based on request



## <a name="Send">func</a> [Send](/src/target/private.go?s=4407:4456#L134)
``` go
func Send(w http.ResponseWriter, r *http.Request)
```
Send It receives json SendRequest with from, to and payload, returns Status Code 200 and json KeyJSON with encoded key.



## <a name="SendRaw">func</a> [SendRaw](/src/target/private.go?s=1206:1258#L39)
``` go
func SendRaw(w http.ResponseWriter, r *http.Request)
```
SendRaw It receives headers "bb0x-from" and "bb0x-to", payload body and returns Status Code 200 and encoded key plain text.



## <a name="SendSignedTx">func</a> [SendSignedTx](/src/target/private.go?s=5553:5610#L172)
``` go
func SendSignedTx(w http.ResponseWriter, r *http.Request)
```
SendSignedTx It receives header "bb0x-to" and raw transaction hash body and returns Status Code 200 and transaction hash.



## <a name="SetLogger">func</a> [SetLogger](/src/target/log.go?s=1017:1054#L30)
``` go
func SetLogger(loggers *logrus.Entry)
```
SetLogger set the logger



## <a name="StoreRaw">func</a> [StoreRaw](/src/target/public.go?s=11340:11393#L331)
``` go
func StoreRaw(w http.ResponseWriter, r *http.Request)
```
StoreRaw It receives a POST request with a payload and returns Status Code 201 with a payload generated hash, on error returns Status Code 500.



## <a name="TransactionDelete">func</a> [TransactionDelete](/src/target/public.go?s=10416:10478#L307)
``` go
func TransactionDelete(w http.ResponseWriter, r *http.Request)
```
TransactionDelete It receives a DELETE request with a key on path string and returns 204 if succeed, 404 otherwise.



## <a name="TransactionGet">func</a> [TransactionGet](/src/target/private.go?s=9853:9912#L306)
``` go
func TransactionGet(w http.ResponseWriter, r *http.Request)
```
TransactionGet it receives a GET request with a hash on path and query var "to" with encoded hash and to, returns decrypted payload



## <a name="UnknownRequest">func</a> [UnknownRequest](/src/target/common.go?s=1696:1755#L56)
``` go
func UnknownRequest(w http.ResponseWriter, r *http.Request)
```
UnknownRequest will debug unknown reqs



## <a name="Upcheck">func</a> [Upcheck](/src/target/common.go?s=1369:1421#L43)
``` go
func Upcheck(w http.ResponseWriter, r *http.Request)
```
Upcheck Request path "/upcheck", response plain text upcheck message.




## <a name="KeyJSON">type</a> [KeyJSON](/src/target/types.go?s=1260:1383#L35)
``` go
type KeyJSON struct {
    // Key is the key that can be used to retrieve the submitted transaction.
    Key string `json:"key"`
}

```
KeyJSON marshal/unmarshal a key










## <a name="PeerURL">type</a> [PeerURL](/src/target/types.go?s=2215:2263#L67)
``` go
type PeerURL struct {
    URL string `json:"url"`
}

```
PeerURL will marshal/unmarshal url










## <a name="ReceiveRequest">type</a> [ReceiveRequest](/src/target/types.go?s=1431:1510#L41)
``` go
type ReceiveRequest struct {
    Key string `json:"key"`
    To  string `json:"to"`
}

```
ReceiveRequest marshal/unmarshal key and to










### <a name="ReceiveRequest.Parse">func</a> (\*ReceiveRequest) [Parse](/src/target/types.go?s=3120:3179#L95)
``` go
func (e *ReceiveRequest) Parse() ([]byte, []byte, []string)
```
Parse will process receiving parsing




## <a name="ReceiveResponse">type</a> [ReceiveResponse](/src/target/types.go?s=1561:1625#L47)
``` go
type ReceiveResponse struct {
    Payload string `json:"payload"`
}

```
ReceiveResponse will marshal/unmarshal payload










## <a name="ResendRequest">type</a> [ResendRequest](/src/target/types.go?s=1683:2012#L52)
``` go
type ResendRequest struct {
    // Type is the resend request type. It should be either "all" or "individual" depending on if
    // you want to request an individual transaction, or all transactions associated with a node.
    Type      string `json:"type"`
    PublicKey string `json:"publicKey"`
    Key       string `json:"key,omitempty"`
}

```
ResendRequest will marshal/unmarshal type, pub and pk










## <a name="SendRequest">type</a> [SendRequest](/src/target/types.go?s=916:1224#L25)
``` go
type SendRequest struct {
    // Payload is the transaction payload data we wish to store.
    Payload string `json:"payload"`
    // From is the sender node identification.
    From string `json:"from"`
    // To is a list of the recipient nodes that should be privy to this transaction payload.
    To []string `json:"to"`
}

```
SendRequest will marshal/unmarshal payload from and to










### <a name="SendRequest.Parse">func</a> (\*SendRequest) [Parse](/src/target/types.go?s=2299:2365#L72)
``` go
func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string)
```
Parse will process send parsing




## <a name="StoreRawRequest">type</a> [StoreRawRequest](/src/target/types.go?s=2072:2176#L61)
``` go
type StoreRawRequest struct {
    Payload string `json:"payload"`
    From    string `json:"from,omitempty"`
}

```
StoreRawRequest will marshal/unmarshal payload and from














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
