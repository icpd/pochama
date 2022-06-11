package sshtunnel

import (
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Local  string
	Server string
	Remote string
	Config *ssh.ClientConfig
}

func (t *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", t.Local)
	if err != nil {
		return err
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("error closing listener: %v", err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go t.forward(conn)
	}
}

func (t *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", t.Server, t.Config)
	if err != nil {
		log.Printf("server dial error: %s", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", t.Remote)
	if err != nil {
		log.Printf("remote dial error: %s", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Printf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func NewSSHTunnel(options ...OptionFunc) *SSHTunnel {
	sshTunnel := &SSHTunnel{
		Config: &ssh.ClientConfig{
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
	}

	for _, option := range options {
		option(sshTunnel)
	}

	return sshTunnel
}

type OptionFunc func(*SSHTunnel)

func (f OptionFunc) apply(st *SSHTunnel) {
	f(st)
}

func WithLocal(endpoint string) OptionFunc {
	return func(st *SSHTunnel) {
		st.Local = endpoint
	}
}

func WithServer(endpoint string) OptionFunc {
	return func(st *SSHTunnel) {
		st.Server = endpoint
	}
}

func WithAuth(user string, auth ssh.AuthMethod) OptionFunc {
	return func(st *SSHTunnel) {
		if st.Config == nil {
			st.Config = &ssh.ClientConfig{}
		}

		st.Config.User = user
		st.Config.Auth = append(st.Config.Auth, auth)
	}
}

func WithRemote(endpoint string) OptionFunc {
	return func(st *SSHTunnel) {
		st.Remote = endpoint
	}
}

func WithHostKeyCallback(HostKeyCallback ssh.HostKeyCallback) OptionFunc {
	return func(st *SSHTunnel) {
		if st.Config == nil {
			st.Config = &ssh.ClientConfig{}
		}

		st.Config.HostKeyCallback = HostKeyCallback
	}
}
