package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

func HandleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

var userconfig struct {
	Address        string
	PrivateKeyPath string
	BinPath        string
}

func init() {
	bytes, err := os.ReadFile("config.json")
	HandleError(err)
	HandleError(json.Unmarshal(bytes, &userconfig))
}

func generateRandomSigner() ssh.Signer {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	HandleError(err)
	signer, err := ssh.NewSignerFromKey(privkey)
	HandleError(err)
	return signer
}

var serverconfig *ssh.ServerConfig

func init() {
	serverconfig = &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
	}

	var signer ssh.Signer
	privkeybytes, err := os.ReadFile(userconfig.PrivateKeyPath)
	if err == nil {
		signer, err = ssh.ParsePrivateKey(privkeybytes)
		HandleError(err)
	} else {
		log.Println("The private key file either was not supplied or can not be opened, falling back to randomly generated private key.")
		signer = generateRandomSigner()
	}
	serverconfig.AddHostKey(signer)
}

func supplyEnvironmentVariables(_ ssh.Channel, cmd *exec.Cmd) {
	// Also supply password, client key, client username/id
	// but afaik that's not even possible with the current state of x/crypto/ssh
	cmd.Env = []string{"PATH=" + userconfig.BinPath}
}

func spawnProcess(channel ssh.Channel) {
	cmd := exec.Command(userconfig.BinPath)

	supplyEnvironmentVariables(channel, cmd)

	stdinp, err := cmd.StdinPipe()
	HandleError(err)
	stdoup, err := cmd.StdoutPipe()
	HandleError(err)
	go io.Copy(stdinp, channel)
	go io.Copy(channel, stdoup)

	log.Println("Spawning process: ", userconfig.BinPath)
	HandleError(cmd.Run())
	log.Println("Done")
}

func handleConnection(nConn net.Conn) {
	_, chans, reqs, err := ssh.NewServerConn(nConn, serverconfig)
	HandleError(err)

	go ssh.DiscardRequests(reqs)
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		HandleError(err)

		go func(in <-chan *ssh.Request) {
			for req := range in {
				if req.WantReply {
					req.Reply(true, nil)
				}
			}
		}(requests)

		go func() {
			defer channel.Close()
			spawnProcess(channel)
		}()
	}
}

func main() {
	listener, err := net.Listen("tcp", userconfig.Address)
	HandleError(err)

	defer listener.Close()
	for {
		nConn, err := listener.Accept()
		HandleError(err)
		go handleConnection(nConn)
	}
}
