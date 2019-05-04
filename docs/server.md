

# server
`import "Smilo-blackbox/src/server"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func AddPeer(node *discover.Node) error](#AddPeer)
* [func GetExternalIP() (string, error)](#GetExternalIP)
* [func GetPeerListSend(peer *p2p.Peer, rw p2p.MsgReadWriter)](#GetPeerListSend)
* [func GetPeerNodeID(id string) string](#GetPeerNodeID)
* [func InitP2PPeers(peers []model.PeerNode)](#InitP2PPeers)
* [func InitP2PServer(bootstrapNodes []*discover.Node) (*p2p.Server, error)](#InitP2PServer)
* [func InitP2p()](#InitP2p)
* [func InitRouting() (*mux.Router, *mux.Router)](#InitRouting)
* [func IsPeerAlreadyAdded(parsedPeer *discover.Node) bool](#IsPeerAlreadyAdded)
* [func NewServer(Port string) (*http.Server, *http.Server)](#NewServer)
* [func PeerList(p2pMessage Message)](#PeerList)
* [func SendMsg(peer *p2p.Peer, rw p2p.MsgReadWriter, err error, outmsg Message)](#SendMsg)
* [func SetLogger(loggers *logrus.Entry)](#SetLogger)
* [func StartServer()](#StartServer)
* [func SubscribeP2P()](#SubscribeP2P)
* [type Message](#Message)
* [type Peer](#Peer)


#### <a name="pkg-files">Package files</a>
[p2p.go](/src/Smilo-blackbox/src/server/p2p.go) [p2p_func.go](/src/Smilo-blackbox/src/server/p2p_func.go) [server.go](/src/Smilo-blackbox/src/server/server.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (

    //StormDBPeers is the main object for peers db
    StormDBPeers *storm.DB

    //DefaultExpirationTime is the default expiration time used on the database
    DefaultExpirationTime = &buntdb.SetOptions{Expires: false} // never expire

    //PUBLIC_SERVER_READ_TIMEOUT_STR will be used to hold env var
    PUBLIC_SERVER_READ_TIMEOUT_STR = os.Getenv("PUBLIC_SERVER_READ_TIMEOUT")

    //PUBLIC_SERVER_WRITE_TIMEOUT_STR will be used to hold env var
    PUBLIC_SERVER_WRITE_TIMEOUT_STR = os.Getenv("PUBLIC_SERVER_WRITE_TIMEOUT")

    //PRIVATE_SERVER_READ_TIMEOUT_STR will be used to hold env var
    PRIVATE_SERVER_READ_TIMEOUT_STR = os.Getenv("PRIVATE_SERVER_READ_TIMEOUT")

    //PRIVATE_SERVER_WRITE_TIMEOUT_STR will be used to hold env var
    PRIVATE_SERVER_WRITE_TIMEOUT_STR = os.Getenv("PRIVATE_SERVER_WRITE_TIMEOUT")
    //PUBLIC_SERVER_READ_TIMEOUT will be used to hold env var
    PUBLIC_SERVER_READ_TIMEOUT = 120
    //PUBLIC_SERVER_WRITE_TIMEOUT will be used to hold env var
    PUBLIC_SERVER_WRITE_TIMEOUT = 120
    //PRIVATE_SERVER_READ_TIMEOUT   will be used to hold env var
    PRIVATE_SERVER_READ_TIMEOUT = 60
    //PRIVATE_SERVER_WRITE_TIMEOUT  will be used to hold env var
    PRIVATE_SERVER_WRITE_TIMEOUT = 60
)
```


## <a name="AddPeer">func</a> [AddPeer](/src/target/p2p_func.go?s=5296:5335#L220)
``` go
func AddPeer(node *discover.Node) error
```
AddPeer will add a peer



## <a name="GetExternalIP">func</a> [GetExternalIP](/src/target/p2p_func.go?s=11209:11245#L458)
``` go
func GetExternalIP() (string, error)
```
GetExternalIP Get the external IP



## <a name="GetPeerListSend">func</a> [GetPeerListSend](/src/target/p2p_func.go?s=9289:9347#L371)
``` go
func GetPeerListSend(peer *p2p.Peer, rw p2p.MsgReadWriter)
```
GetPeerListSend will get peers



## <a name="GetPeerNodeID">func</a> [GetPeerNodeID](/src/target/p2p_func.go?s=10131:10167#L408)
``` go
func GetPeerNodeID(id string) string
```
GetPeerNodeID will get peer node



## <a name="InitP2PPeers">func</a> [InitP2PPeers](/src/target/p2p_func.go?s=5708:5749#L238)
``` go
func InitP2PPeers(peers []model.PeerNode)
```
InitP2PPeers will init peers



## <a name="InitP2PServer">func</a> [InitP2PServer](/src/target/p2p_func.go?s=1236:1308#L41)
``` go
func InitP2PServer(bootstrapNodes []*discover.Node) (*p2p.Server, error)
```
InitP2PServer will init p2p peers



## <a name="InitP2p">func</a> [InitP2p](/src/target/p2p_func.go?s=2637:2651#L99)
``` go
func InitP2p()
```
InitP2p will init p2p



## <a name="InitRouting">func</a> [InitRouting](/src/target/server.go?s=7219:7264#L254)
``` go
func InitRouting() (*mux.Router, *mux.Router)
```
InitRouting will init routing



## <a name="IsPeerAlreadyAdded">func</a> [IsPeerAlreadyAdded](/src/target/p2p_func.go?s=10612:10667#L430)
``` go
func IsPeerAlreadyAdded(parsedPeer *discover.Node) bool
```
IsPeerAlreadyAdded check if peer already connected



## <a name="NewServer">func</a> [NewServer](/src/target/server.go?s=4267:4323#L139)
``` go
func NewServer(Port string) (*http.Server, *http.Server)
```
NewServer will create a new http server instance -- pub and private



## <a name="PeerList">func</a> [PeerList](/src/target/p2p_func.go?s=8830:8863#L353)
``` go
func PeerList(p2pMessage Message)
```
PeerList will init all peers provided on the p2p message



## <a name="SendMsg">func</a> [SendMsg](/src/target/p2p.go?s=3862:3939#L174)
``` go
func SendMsg(peer *p2p.Peer, rw p2p.MsgReadWriter, err error, outmsg Message)
```
SendMsg will send a message



## <a name="SetLogger">func</a> [SetLogger](/src/target/server.go?s=4049:4086#L129)
``` go
func SetLogger(loggers *logrus.Entry)
```
SetLogger set the logger



## <a name="StartServer">func</a> [StartServer](/src/target/server.go?s=4839:4857#L157)
``` go
func StartServer()
```
StartServer will start the server



## <a name="SubscribeP2P">func</a> [SubscribeP2P](/src/target/p2p_func.go?s=4102:4121#L171)
``` go
func SubscribeP2P()
```
SubscribeP2P will subscribe p2p




## <a name="Message">type</a> [Message](/src/target/p2p.go?s=951:1035#L28)
``` go
type Message struct {
    Header string `json:"content"`
    Body   string `json:"body"`
}
```
Message holds header and body










## <a name="Peer">type</a> [Peer](/src/target/p2p.go?s=1068:1142#L34)
``` go
type Peer struct {
    ID         string
    Dest       string
    SourcePort int
}
```
Peer is the main peer struct














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
