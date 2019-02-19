

# api
`import "Smilo-blackbox/src/server/api"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func Api(w http.ResponseWriter, r *http.Request)](#Api)
* [func ConfigPeersGet(w http.ResponseWriter, r *http.Request)](#ConfigPeersGet)
* [func ConfigPeersPut(w http.ResponseWriter, r *http.Request)](#ConfigPeersPut)
* [func Delete(w http.ResponseWriter, r *http.Request)](#Delete)
* [func GetPartyInfo(w http.ResponseWriter, r *http.Request)](#GetPartyInfo)
* [func GetVersion(w http.ResponseWriter, r *http.Request)](#GetVersion)
* [func Metrics(w http.ResponseWriter, r *http.Request)](#Metrics)
* [func Push(w http.ResponseWriter, r *http.Request)](#Push)
* [func PushTransactionForOtherNodes(encryptedTransaction data.Encrypted_Transaction, recipient []byte)](#PushTransactionForOtherNodes)
* [func Receive(w http.ResponseWriter, r *http.Request)](#Receive)
* [func ReceiveRaw(w http.ResponseWriter, r *http.Request)](#ReceiveRaw)
* [func Resend(w http.ResponseWriter, r *http.Request)](#Resend)
* [func RetrieveAndDecryptPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) []byte](#RetrieveAndDecryptPayload)
* [func RetrieveJsonPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte)](#RetrieveJsonPayload)
* [func Send(w http.ResponseWriter, r *http.Request)](#Send)
* [func SendRaw(w http.ResponseWriter, r *http.Request)](#SendRaw)
* [func SetLogger(loggers *logrus.Entry)](#SetLogger)
* [func TransactionDelete(w http.ResponseWriter, r *http.Request)](#TransactionDelete)
* [func TransactionGet(w http.ResponseWriter, r *http.Request)](#TransactionGet)
* [func UnknownRequest(w http.ResponseWriter, r *http.Request)](#UnknownRequest)
* [func Upcheck(w http.ResponseWriter, r *http.Request)](#Upcheck)
* [type DeleteRequest](#DeleteRequest)
* [type PeerUrl](#PeerUrl)
* [type ReceiveRequest](#ReceiveRequest)
  * [func (e *ReceiveRequest) Parse() ([]byte, []byte, []string)](#ReceiveRequest.Parse)
* [type ReceiveResponse](#ReceiveResponse)
* [type ResendRequest](#ResendRequest)
* [type SendRequest](#SendRequest)
  * [func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string)](#SendRequest.Parse)
* [type SendResponse](#SendResponse)


#### <a name="pkg-files">Package files</a>
[common.go](/src/Smilo-blackbox/src/server/api/common.go) [log.go](/src/Smilo-blackbox/src/server/api/log.go) [private.go](/src/Smilo-blackbox/src/server/api/private.go) [public.go](/src/Smilo-blackbox/src/server/api/public.go) [types.go](/src/Smilo-blackbox/src/server/api/types.go) 





## <a name="Api">func</a> [Api](/src/target/common.go?s=1421:1469#L45)
``` go
func Api(w http.ResponseWriter, r *http.Request)
```
Request path "/api", response json rest api spec.



## <a name="ConfigPeersGet">func</a> [ConfigPeersGet](/src/target/public.go?s=9236:9295#L258)
``` go
func ConfigPeersGet(w http.ResponseWriter, r *http.Request)
```
TODO
ConfigPeersGet Receive a GET request with index on path and return Status Code 200 and Peer json containing url, Status Code 404 if not found.



## <a name="ConfigPeersPut">func</a> [ConfigPeersPut](/src/target/public.go?s=8667:8726#L242)
``` go
func ConfigPeersPut(w http.ResponseWriter, r *http.Request)
```
TODO
ConfigPeersPut It receives a PUT request with a json containing a Peer url and returns Status Code 200.



## <a name="Delete">func</a> [Delete](/src/target/public.go?s=6824:6875#L189)
``` go
func Delete(w http.ResponseWriter, r *http.Request)
```
Delete Deprecated API
It receives a POST request with a json containing a DeleteRequest with key and returns Status 200 if succeed, 404 otherwise.



## <a name="GetPartyInfo">func</a> [GetPartyInfo](/src/target/public.go?s=1252:1309#L41)
``` go
func GetPartyInfo(w http.ResponseWriter, r *http.Request)
```
TODO
GetPartyInfo It receives a POST request with a json containing url and key, returns local publicKeys and a proof that private key is known.



## <a name="GetVersion">func</a> [GetVersion](/src/target/common.go?s=1105:1160#L35)
``` go
func GetVersion(w http.ResponseWriter, r *http.Request)
```
Request path "/version", response plain text version ID



## <a name="Metrics">func</a> [Metrics](/src/target/public.go?s=10058:10110#L282)
``` go
func Metrics(w http.ResponseWriter, r *http.Request)
```
TODO
Metrics Receive a GET request and return Status Code 200 and server internal status information in plain text.



## <a name="Push">func</a> [Push](/src/target/public.go?s=2536:2585#L66)
``` go
func Push(w http.ResponseWriter, r *http.Request)
```
Push It receives a POST request with a payload and returns Status Code 201 with a payload generated hash, on error returns Status Code 500.



## <a name="PushTransactionForOtherNodes">func</a> [PushTransactionForOtherNodes](/src/target/common.go?s=2599:2699#L82)
``` go
func PushTransactionForOtherNodes(encryptedTransaction data.Encrypted_Transaction, recipient []byte)
```


## <a name="Receive">func</a> [Receive](/src/target/private.go?s=5736:5788#L170)
``` go
func Receive(w http.ResponseWriter, r *http.Request)
```
Receive is a Deprecated API
It receives a ReceiveRequest json with an encoded key (hash) and to values, returns decrypted payload



## <a name="ReceiveRaw">func</a> [ReceiveRaw](/src/target/public.go?s=3785:3840#L107)
``` go
func ReceiveRaw(w http.ResponseWriter, r *http.Request)
```
ReceiveRaw Receive a GET request with header params bb0x-key and bb0x-to, return unencrypted payload



## <a name="Resend">func</a> [Resend](/src/target/public.go?s=5483:5534#L153)
``` go
func Resend(w http.ResponseWriter, r *http.Request)
```
Resend It receives a POST request with a json ResendRequest containing type (INDIVIDUAL, ALL), publicKey and key(for individual requests),
it returns encoded payload for INDIVIDUAL or it does one push request for each payload and returns empty for type ALL.



## <a name="RetrieveAndDecryptPayload">func</a> [RetrieveAndDecryptPayload](/src/target/common.go?s=1921:2021#L62)
``` go
func RetrieveAndDecryptPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) []byte
```


## <a name="RetrieveJsonPayload">func</a> [RetrieveJsonPayload](/src/target/common.go?s=1571:1658#L53)
``` go
func RetrieveJsonPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte)
```


## <a name="Send">func</a> [Send](/src/target/private.go?s=4119:4168#L123)
``` go
func Send(w http.ResponseWriter, r *http.Request)
```
Send It receives json SendRequest with from, to and payload, returns Status Code 200 and json SendResponse with encoded key.



## <a name="SendRaw">func</a> [SendRaw](/src/target/private.go?s=1190:1242#L38)
``` go
func SendRaw(w http.ResponseWriter, r *http.Request)
```
SendRaw It receives headers "bb0x-from" and "bb0x-to", payload body and returns Status Code 200 and encoded key plain text.



## <a name="SetLogger">func</a> [SetLogger](/src/target/log.go?s=1031:1068#L30)
``` go
func SetLogger(loggers *logrus.Entry)
```
SetLogger set the logger



## <a name="TransactionDelete">func</a> [TransactionDelete](/src/target/public.go?s=7887:7949#L220)
``` go
func TransactionDelete(w http.ResponseWriter, r *http.Request)
```
TransactionDelete It receives a DELETE request with a key on path string and returns 204 if succeed, 404 otherwise.



## <a name="TransactionGet">func</a> [TransactionGet](/src/target/private.go?s=6463:6522#L195)
``` go
func TransactionGet(w http.ResponseWriter, r *http.Request)
```
TransactionGet it receives a GET request with a hash on path and query var "to" with encoded hash and to, returns decrypted payload



## <a name="UnknownRequest">func</a> [UnknownRequest](/src/target/common.go?s=1476:1535#L49)
``` go
func UnknownRequest(w http.ResponseWriter, r *http.Request)
```


## <a name="Upcheck">func</a> [Upcheck](/src/target/common.go?s=1271:1323#L40)
``` go
func Upcheck(w http.ResponseWriter, r *http.Request)
```
Request path "/upcheck", response plain text upcheck message.




## <a name="DeleteRequest">type</a> [DeleteRequest](/src/target/types.go?s=1446:1500#L47)
``` go
type DeleteRequest struct {
    Key string `json:"key"`
}
```









## <a name="PeerUrl">type</a> [PeerUrl](/src/target/types.go?s=1833:1881#L59)
``` go
type PeerUrl struct {
    Url string `json:"url"`
}
```









## <a name="ReceiveRequest">type</a> [ReceiveRequest](/src/target/types.go?s=1299:1378#L38)
``` go
type ReceiveRequest struct {
    Key string `json:"key"`
    To  string `json:"to"`
}
```









### <a name="ReceiveRequest.Parse">func</a> (\*ReceiveRequest) [Parse](/src/target/types.go?s=2665:2724#L85)
``` go
func (e *ReceiveRequest) Parse() ([]byte, []byte, []string)
```



## <a name="ReceiveResponse">type</a> [ReceiveResponse](/src/target/types.go?s=1380:1444#L43)
``` go
type ReceiveResponse struct {
    Payload string `json:"payload"`
}
```









## <a name="ResendRequest">type</a> [ResendRequest](/src/target/types.go?s=1502:1831#L51)
``` go
type ResendRequest struct {
    // Type is the resend request type. It should be either "all" or "individual" depending on if
    // you want to request an individual transaction, or all transactions associated with a node.
    Type      string `json:"type"`
    PublicKey string `json:"publicKey"`
    Key       string `json:"key,omitempty"`
}
```









## <a name="SendRequest">type</a> [SendRequest](/src/target/types.go?s=859:1167#L24)
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









### <a name="SendRequest.Parse">func</a> (\*SendRequest) [Parse](/src/target/types.go?s=1883:1949#L63)
``` go
func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string)
```



## <a name="SendResponse">type</a> [SendResponse](/src/target/types.go?s=1169:1297#L33)
``` go
type SendResponse struct {
    // Key is the key that can be used to retrieve the submitted transaction.
    Key string `json:"key"`
}
```













- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
