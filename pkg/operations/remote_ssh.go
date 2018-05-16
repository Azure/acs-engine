package operations

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

// RemoteRun executes remote command
func RemoteRun(user string, addr string, privateKey string, cmd string) (string, error) {
	// Create the Signer for this private key.
	key, err := ioutil.ReadFile("/home/dima/.ssh/id_rsa_test")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// Authentication
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error { return nil },
	}
	// Connect
	client, err := ssh.Dial("tcp", addr+":22", config)
	if err != nil {
		return "", err
	}
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b // get output

	err = session.Run(cmd)
	return b.String(), err
}
