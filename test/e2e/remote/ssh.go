package remote

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Connection is
type Connection struct {
	Host           string
	Port           string
	User           string
	PrivateKeyPath string
}

// NewConnection will build and return a new Connection object
func NewConnection(host, port, user, keyPath string) *Connection {
	return &Connection{
		Host:           host,
		Port:           port,
		User:           user,
		PrivateKeyPath: keyPath,
	}
}

// Execute will run cmd against a remote host
func (c *Connection) Execute(cmd string) ([]byte, error) {
	connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
	out, err := exec.Command("ssh", "-i", c.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, "-p", c.Port, cmd).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Connection) Write(data, path string) error {
	cmd := fmt.Sprintf("echo %s > %s", data, path)
	connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
	out, err := exec.Command("ssh", "-i", c.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, "-p", c.Port, cmd).CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
		return err
	}
	return nil
}

func (c *Connection) Read(path string) ([]byte, error) {
	cmd := fmt.Sprintf("cat %s", path)
	connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
	out, err := exec.Command("ssh", "-i", c.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, "-p", c.Port, cmd).CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
		return nil, err
	}
	return out, nil
}

// ExecuteWithRetries will keep retrying a command until it does not return an error or the duration is exceeded
func (c *Connection) ExecuteWithRetries(cmd string, sleep, duration time.Duration) ([]byte, error) {
	outCh := make(chan []byte, 1)
	errCh := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- fmt.Errorf("Timeout exceeded (%s) while waiting for command to not return an error: %s", duration.String(), cmd)
			default:
				out, err := c.Execute(cmd)
				if err == nil {
					outCh <- out
				}
				time.Sleep(sleep)
			}
		}
	}()
	for {
		select {
		case err := <-errCh:
			return nil, err
		case out := <-outCh:
			return out, nil
		}
	}
}
