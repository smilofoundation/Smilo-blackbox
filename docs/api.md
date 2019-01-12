

# api
`import "Smilo-blackbox/src/server/api"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
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
* [func Upcheck(w http.ResponseWriter, r *http.Request)](#Upcheck)
* [type DeleteRequest](#DeleteRequest)
* [type PartyInfoResponse](#PartyInfoResponse)
* [type ReceiveRequest](#ReceiveRequest)
  * [func (e *ReceiveRequest) Parse() ([]byte, []byte, []string)](#ReceiveRequest.Parse)
* [type ReceiveResponse](#ReceiveResponse)
* [type ResendRequest](#ResendRequest)
* [type SendRequest](#SendRequest)
  * [func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string)](#SendRequest.Parse)
* [type SendResponse](#SendResponse)
* [type UpdatePartyInfo](#UpdatePartyInfo)


#### <a name="pkg-files">Package files</a>
[common.go](/src/Smilo-blackbox/src/server/api/common.go) [log.go](/src/Smilo-blackbox/src/server/api/log.go) [private.go](/src/Smilo-blackbox/src/server/api/private.go) [public.go](/src/Smilo-blackbox/src/server/api/public.go) [types.go](/src/Smilo-blackbox/src/server/api/types.go) 


## <a name="pkg-constants">Constants</a>
``` go
const BlackBoxVersion = "Smilo Black Box 0.1.0"
```
``` go
const HeaderFrom = "c11n-from"
```
``` go
const HeaderKey = "c11n-key"
```
``` go
const HeaderTo = "c11n-to"
```
``` go
const UpcheckMessage = "I'm up!"
```



## <a name="Api">func</a> [Api](/src/target/common.go?s=695:743#L32)
``` go
func Api(w http.ResponseWriter, r *http.Request)
```
Request path "/api", response json rest api spec.



## <a name="ConfigPeersGet">func</a> [ConfigPeersGet](/src/target/public.go?s=5952:6011#L172)
``` go
func ConfigPeersGet(w http.ResponseWriter, r *http.Request)
```
TODO
Receive a GET request with index on path and return Status Code 200 and Peer json containing url, Status Code 500 otherwise



## <a name="ConfigPeersPut">func</a> [ConfigPeersPut](/src/target/public.go?s=5528:5587#L160)
``` go
func ConfigPeersPut(w http.ResponseWriter, r *http.Request)
```
TODO
It receives a PUT request with a json containing a Peer and returns Status Code 200 and the new peer URL.



## <a name="Delete">func</a> [Delete](/src/target/public.go?s=3901:3952#L117)
``` go
func Delete(w http.ResponseWriter, r *http.Request)
```
Deprecated API
It receives a POST request with a json containing a DeleteRequest with key and returns Status 200 if succeed, 404 otherwise.



## <a name="GetPartyInfo">func</a> [GetPartyInfo](/src/target/public.go?s=319:376#L21)
``` go
func GetPartyInfo(w http.ResponseWriter, r *http.Request)
```
TODO
It receives a POST request with a binary encoded PartyInfo, updates it and returns updated PartyInfo encoded.



## <a name="GetVersion">func</a> [GetVersion](/src/target/common.go?s=391:446#L22)
``` go
func GetVersion(w http.ResponseWriter, r *http.Request)
```
Request path "/version", response plain text version ID



## <a name="Metrics">func</a> [Metrics](/src/target/public.go?s=6284:6336#L183)
``` go
func Metrics(w http.ResponseWriter, r *http.Request)
```
TODO
Receive a GET request and return Status Code 200 and server internal status information in plain text.



## <a name="Push">func</a> [Push](/src/target/public.go?s=521:570#L26)
``` go
func Push(w http.ResponseWriter, r *http.Request)
```
It receives a POST request with a payload and returns Status Code 201 with a payload generated hash, on error returns Status Code 500.



## <a name="PushTransactionForOtherNodes">func</a> [PushTransactionForOtherNodes](/src/target/common.go?s=1698:1798#L61)
``` go
func PushTransactionForOtherNodes(encryptedTransaction data.Encrypted_Transaction, recipient []byte)
```


## <a name="Receive">func</a> [Receive](/src/target/private.go?s=3541:3593#L107)
``` go
func Receive(w http.ResponseWriter, r *http.Request)
```
Deprecated API
It receives a ReceiveRequest json with an encoded key (hash) and to values, returns decrypted payload



## <a name="ReceiveRaw">func</a> [ReceiveRaw](/src/target/public.go?s=1593:1648#L59)
``` go
func ReceiveRaw(w http.ResponseWriter, r *http.Request)
```
Receive a GET request with header params c11n-key and c11n-to, return unencrypted payload



## <a name="Resend">func</a> [Resend](/src/target/public.go?s=2691:2742#L87)
``` go
func Resend(w http.ResponseWriter, r *http.Request)
```
It receives a POST request with a json ResendRequest containing type (INDIVIDUAL, ALL), publicKey and key(for individual requests),
it returns encoded payload for INDIVIDUAL or it does one push request for each payload and returns empty for type ALL.



## <a name="RetrieveAndDecryptPayload">func</a> [RetrieveAndDecryptPayload](/src/target/common.go?s=1100:1200#L45)
``` go
func RetrieveAndDecryptPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte) []byte
```


## <a name="RetrieveJsonPayload">func</a> [RetrieveJsonPayload](/src/target/common.go?s=750:837#L36)
``` go
func RetrieveJsonPayload(w http.ResponseWriter, r *http.Request, key []byte, to []byte)
```


## <a name="Send">func</a> [Send](/src/target/private.go?s=2132:2181#L68)
``` go
func Send(w http.ResponseWriter, r *http.Request)
```
It receives json SendRequest with from, to and payload, returns Status Code 200 and json SendResponse with encoded key.



## <a name="SendRaw">func</a> [SendRaw](/src/target/private.go?s=317:369#L19)
``` go
func SendRaw(w http.ResponseWriter, r *http.Request)
```
It receives headers "c11n-from" and "c11n-to", payload body and returns Status Code 200 and encoded key plain text.



## <a name="SetLogger">func</a> [SetLogger](/src/target/log.go?s=223:260#L14)
``` go
func SetLogger(loggers *logrus.Entry)
```
SetLogger set the logger



## <a name="TransactionDelete">func</a> [TransactionDelete](/src/target/public.go?s=4826:4888#L142)
``` go
func TransactionDelete(w http.ResponseWriter, r *http.Request)
```
It receives a DELETE request with a key on path string and returns 204 if succeed, 404 otherwise.



## <a name="TransactionGet">func</a> [TransactionGet](/src/target/private.go?s=4173:4232#L128)
``` go
func TransactionGet(w http.ResponseWriter, r *http.Request)
```
it receives a GET request with a hash on path and query var "to" with encoded hash and to, returns decrypted payload



## <a name="Upcheck">func</a> [Upcheck](/src/target/common.go?s=551:603#L27)
``` go
func Upcheck(w http.ResponseWriter, r *http.Request)
```
Request path "/upcheck", response plain text upcheck message.




## <a name="DeleteRequest">type</a> [DeleteRequest](/src/target/types.go?s=638:692#L31)
``` go
type DeleteRequest struct {
    Key string `json:"key"`
}
```









## <a name="PartyInfoResponse">type</a> [PartyInfoResponse](/src/target/types.go?s=1198:1264#L49)
``` go
type PartyInfoResponse struct {
    Payload []byte `json:"payload"`
}
```









## <a name="ReceiveRequest">type</a> [ReceiveRequest](/src/target/types.go?s=491:570#L22)
``` go
type ReceiveRequest struct {
    Key string `json:"key"`
    To  string `json:"to"`
}
```









### <a name="ReceiveRequest.Parse">func</a> (\*ReceiveRequest) [Parse](/src/target/types.go?s=2054:2113#L75)
``` go
func (e *ReceiveRequest) Parse() ([]byte, []byte, []string)
```



## <a name="ReceiveResponse">type</a> [ReceiveResponse](/src/target/types.go?s=572:636#L27)
``` go
type ReceiveResponse struct {
    Payload string `json:"payload"`
}
```









## <a name="ResendRequest">type</a> [ResendRequest](/src/target/types.go?s=694:1023#L35)
``` go
type ResendRequest struct {
    // Type is the resend request type. It should be either "all" or "individual" depending on if
    // you want to request an individual transaction, or all transactions associated with a node.
    Type      string `json:"type"`
    PublicKey string `json:"publicKey"`
    Key       string `json:"key,omitempty"`
}
```









## <a name="SendRequest">type</a> [SendRequest](/src/target/types.go?s=51:359#L8)
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









### <a name="SendRequest.Parse">func</a> (\*SendRequest) [Parse](/src/target/types.go?s=1266:1332#L53)
``` go
func (e *SendRequest) Parse() ([]byte, []byte, [][]byte, []string)
```



## <a name="SendResponse">type</a> [SendResponse](/src/target/types.go?s=361:489#L17)
``` go
type SendResponse struct {
    // Key is the key that can be used to retrieve the submitted transaction.
    Key string `json:"key"`
}
```









## <a name="UpdatePartyInfo">type</a> [UpdatePartyInfo](/src/target/types.go?s=1025:1196#L43)
``` go
type UpdatePartyInfo struct {
    Url        string            `json:"url"`
    Recipients map[string][]byte `json:"recipients"`
    Parties    map[string]bool   `json:"parties"`
}
```













- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
