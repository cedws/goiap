package proxy

import (
	"context"
	"io"
	"net"

	"github.com/cedws/iapc/iap"
	"github.com/charmbracelet/log"
	"golang.org/x/oauth2/google"
)

func Start(listen string, opts []iap.DialOption) {
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConn(opts, conn)
	}
}

func handleConn(opts []iap.DialOption, conn net.Conn) {
	log.Info("Client connected", "client", conn.RemoteAddr())

	opts = append(opts, iap.WithToken(getToken()))
	tun, err := iap.Dial(context.Background(), opts...)
	if err != nil {
		log.Error(err)
		return
	}
	defer tun.Close()
	log.Info("Established connection with proxy", "client", conn.RemoteAddr(), "sid", tun.SessionID())

	go io.Copy(conn, tun)
	io.Copy(tun, conn)

	log.Info("Client disconnected", "client", conn.RemoteAddr())
}

func getToken() string {
	credentials, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tok, err := credentials.TokenSource.Token()
	if err != nil {
		log.Fatal(err)
	}
	return tok.AccessToken
}