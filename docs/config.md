

# config
`import "Smilo-blackbox/src/server/config"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func Init(app *cli.App)](#Init)
* [func LoadConfig(configPath string) error](#LoadConfig)
* [func ReadPrimaryKey(pkFile string) ([]byte, error)](#ReadPrimaryKey)
* [func ReadPublicKey(pubFile string) ([]byte, error)](#ReadPublicKey)
* [type Config](#Config)
* [type Key](#Key)
* [type Keys](#Keys)
* [type Peer](#Peer)
* [type PrivateKey](#PrivateKey)
* [type PrivateKeyBytes](#PrivateKeyBytes)
* [type Server](#Server)


#### <a name="pkg-files">Package files</a>
[config.go](/src/Smilo-blackbox/src/server/config/config.go) [types.go](/src/Smilo-blackbox/src/server/config/types.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (

    //GenerateKeys (cli) uses it for key pair
    GenerateKeys = cli.StringFlag{Name: "generate-keys", Value: "", Usage: "Generate a new keypair"}
    //ConfigFile (cli) uses it for config file name
    ConfigFile = cli.StringFlag{Name: "configfile", Value: "blackbox.conf", Usage: "Config file name"}
    //DBFile (cli) uses it for db file name
    DBFile = cli.StringFlag{Name: "dbfile", Value: "blackbox.db", Usage: "DB file name"}
    //PeersDBFile (cli) uses it for peer db file
    PeersDBFile = cli.StringFlag{Name: "peersdbfile", Value: "blackbox-peers.db", Usage: "Peers DB file name"}
    //Port (cli) uses it for local api public port
    Port = cli.StringFlag{Name: "port", Value: "9000", Usage: "Local port to the Public API"}
    //Socket (cli) uses it for socket
    Socket = cli.StringFlag{Name: "socket", Value: "blackbox.ipc", Usage: "IPC socket to the Private API"}
    //OtherNodes (cli) uses it for other nodes
    OtherNodes = cli.StringFlag{Name: "othernodes", Value: "", Usage: "\"Boot nodes\" to connect"}
    //PublicKeys (cli) uses it for  pub
    PublicKeys = cli.StringFlag{Name: "publickeys", Value: "", Usage: "Public keys"}
    //PrivateKeys (cli) uses it for pk
    PrivateKeys = cli.StringFlag{Name: "privatekeys", Value: "", Usage: "Private keys"}
    //Storage (cli) uses it for  db name
    Storage = cli.StringFlag{Name: "storage", Value: "blackbox.db", Usage: "Database file name"}
    //HostName (cli) uses it for hostname
    HostName = cli.StringFlag{Name: "hostname", Value: "http://localhost", Usage: "HostName for public API"}

    //WorkDir (cli) uses it for work dir
    WorkDir = cli.StringFlag{Name: "workdir", Value: "../../", Usage: ""}
    //IsTLS (cli) uses it for enable/disable https
    IsTLS = cli.BoolFlag{Name: "tls", Usage: "Enable HTTPs communication"}
    //ServCert (cli) uses it for cert
    ServCert = cli.StringFlag{Name: "serv_cert", Value: "", Usage: ""}
    //ServKey (cli) uses it for key
    ServKey = cli.StringFlag{Name: "serv_key", Value: "", Usage: ""}

    //P2PDestination (cli) uses it for p2p dest
    P2PDestination = cli.StringFlag{Name: "p2p_dest", Value: "", Usage: ""}
    //P2PPort (cli) uses it for p2p port
    P2PPort = cli.StringFlag{Name: "p2p_port", Value: "", Usage: ""}
    //CPUProfiling (cli) uses it for CPU profiling data filename
    CPUProfiling = cli.StringFlag{Name: "cpuprofile", Value: "", Usage: "CPU profiling data filename"}
    //P2PEnabled (cli) uses it for enable / disable p2p
    P2PEnabled = cli.BoolFlag{Name: "p2p", Usage: "Enable p2p communication"}
    //RootCert  (cli) uses it for certs
    RootCert = cli.StringFlag{Name: "root_cert", Value: "", Usage: ""}
)
```


## <a name="Init">func</a> [Init](/src/target/config.go?s=3924:3947#L97)
``` go
func Init(app *cli.App)
```
Init will init cli and logs



## <a name="LoadConfig">func</a> [LoadConfig](/src/target/config.go?s=4263:4303#L107)
``` go
func LoadConfig(configPath string) error
```
LoadConfig will load cfg



## <a name="ReadPrimaryKey">func</a> [ReadPrimaryKey](/src/target/config.go?s=6225:6275#L183)
``` go
func ReadPrimaryKey(pkFile string) ([]byte, error)
```
ReadPrimaryKey will read pk



## <a name="ReadPublicKey">func</a> [ReadPublicKey](/src/target/config.go?s=6734:6784#L204)
``` go
func ReadPublicKey(pubFile string) ([]byte, error)
```
ReadPublicKey will read pub




## <a name="Config">type</a> [Config](/src/target/types.go?s=1688:2047#L57)
``` go
type Config struct {
    Server      Server   `json:"server"`
    HostName    string   `json:"hostName"`
    RootCA      []string `json:"rootCA,omitempty"`
    Peers       []Peer   `json:"peer"`
    Keys        Keys     `json:"keys"`
    UnixSocket  string   `json:"socket"`
    DBFile      string   `json:"dbfile,omitempty"`
    PeersDBFile string   `json:"peersdbfile,omitempty"`
}

```
Config json file specification










## <a name="Key">type</a> [Key](/src/target/types.go?s=1424:1524#L45)
``` go
type Key struct {
    PrivateKeyFile string `json:"config"`
    PublicKeyFile  string `json:"publicKey"`
}

```
Key json file specification










## <a name="Keys">type</a> [Keys](/src/target/types.go?s=1557:1653#L51)
``` go
type Keys struct {
    Passwords []string `json:"passwords"`
    KeyData   []Key    `json:"keyData"`
}

```
Keys json file specification










## <a name="Peer">type</a> [Peer](/src/target/types.go?s=1347:1392#L40)
``` go
type Peer struct {
    URL string `json:"url"`
}

```
Peer json file specification










## <a name="PrivateKey">type</a> [PrivateKey](/src/target/types.go?s=982:1080#L25)
``` go
type PrivateKey struct {
    Data PrivateKeyBytes `json:"data"`
    Type string          `json:"type"`
}

```
PrivateKey holds data and type










## <a name="PrivateKeyBytes">type</a> [PrivateKeyBytes](/src/target/types.go?s=887:947#L20)
``` go
type PrivateKeyBytes struct {
    Bytes string `json:"bytes"`
}

```
PrivateKeyBytes Start of Private Key json file specification










## <a name="Server">type</a> [Server](/src/target/types.go?s=1177:1314#L33)
``` go
type Server struct {
    Port    int    `json:"port"`
    TLSCert string `json:"tlscert,omitempty"`
    TLSKey  string `json:"tlskey,omitempty"`
}

```
Server Start of Config json file specification














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
