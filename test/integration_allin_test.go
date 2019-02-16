package test

import (
	"testing"
	"io/ioutil"
	"fmt"
	"github.com/stretchr/testify/require"
	"time"
	osExec "os/exec"

	"path/filepath"
	"os"
	"net"
	"errors"
	"strings"
)


var (
	waitingErr = errors.New("unix socket dial failed")
	upcheckErr = errors.New("http upcheck failed")
	doneErr    = errors.New("done")
)

func checkFunc(tmIPCFile string) error {
	conn, err := net.Dial("unix", tmIPCFile)
	if err != nil {
		return waitingErr
	}
	if _, err := conn.Write([]byte("GET /upcheck HTTP/1.0\r\n\r\n")); err != nil {
		return upcheckErr
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil || string(result) == "I'm up!" {
		return doneErr
	}
	return upcheckErr
}


func runBlackbox(targetNode string) (*osExec.Cmd, error){
	here, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tempdir, err := ioutil.TempDir("", "blackbox")
	if err != nil {
		return nil, err
	}


	cmdStatusChan := make(chan error)
	blackboxCMD := filepath.Join(here, "..", "blackbox")

	blackboxConfigFile := filepath.Join(here, fmt.Sprintf("test%s.conf", targetNode))
	blackboxIPC := filepath.Join(here, fmt.Sprintf("blackbox%s.ipc", targetNode))


	blackboxDBFile := filepath.Join(tempdir, fmt.Sprintf("blackbox%s.db", targetNode))

	cmd := osExec.Command(blackboxCMD, "-configfile", blackboxConfigFile, "-dbfile", blackboxDBFile)
	// run blackbox
	go func() {
		err := cmd.Start()
		cmdStatusChan <- err
	}()
	// wait 30s for blackbox to come up
	var started bool
	go func() {

		for i := 0; i < 10; i++ {
			time.Sleep(3 * time.Second)
			if err := checkFunc(blackboxIPC); err != nil && err == doneErr {
				cmdStatusChan <- err
			} else {
				fmt.Println("Waiting for blackbox to start", "err", err)
			}
		}
		if !started {
			panic("Blackbox never managed to start!")
		}
	}()

	if err := <-cmdStatusChan; err != nil {
		return nil, err
	}
	return cmd, nil

}

func checkblackboxstarted(t *testing.T, err error){
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			t.Fatal(err)
		} else {
			t.Fatal(err)
		}
	}

}

func TestIntegrationAllInSendAll(t *testing.T) {

	t.SkipNow()

	blackboxCmd1, err1 := runBlackbox("1")
	checkblackboxstarted(t, err1)
	defer blackboxCmd1.Process.Kill()

	blackboxCmd2, err2 := runBlackbox("2")
	checkblackboxstarted(t, err2)
	defer blackboxCmd2.Process.Kill()

	blackboxCmd3, err3 := runBlackbox("3")
	checkblackboxstarted(t, err3)
	defer blackboxCmd3.Process.Kill()

	blackboxCmd4, err4 := runBlackbox("4")
	checkblackboxstarted(t, err4)
	defer blackboxCmd4.Process.Kill()

	blackboxCmd5, err5 := runBlackbox("5")
	checkblackboxstarted(t, err5)
	defer blackboxCmd5.Process.Kill()


	//Init()

	waitNodesUp([]int{int(9001),int(9002),int(9003),int(9004),int(9005)})
	time.Sleep(1 * time.Minute)
	to := make([]string, 4)
	to[0] = testServers[1].PublicKey
	to[1] = testServers[2].PublicKey
	to[2] = testServers[3].PublicKey
	to[3] = testServers[4].PublicKey
	sendResponse := sendTestPayload(t, testServers[0], to)

	for i:=1; i<5; i++ {
		receiveResponse := receiveTestPayload(t, testServers[i], sendResponse.Key)
		if receiveResponse.Payload != TEST_PAYLOAD {
			require.Equal(t, TEST_PAYLOAD, receiveResponse.Payload,"Payload not received on Server "+fmt.Sprint(i))
		}
	}
}
