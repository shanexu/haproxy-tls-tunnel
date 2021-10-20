package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	var conn io.ReadWriteCloser

	if enableTls {
		cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		tlsConn, err := tls.Dial("tcp", addr, &config)
		if err != nil {
			log.Fatalf("client: dial: %s", err)
		}
		conn = tlsConn
		log.Println("client: connected to: ", tlsConn.RemoteAddr())
		state := tlsConn.ConnectionState()
		for _, v := range state.PeerCertificates {
			fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
			fmt.Println(v.Subject)
		}
		log.Println("client: handshake: ", state.HandshakeComplete)
		log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	} else {
		var err error
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("client: dial: %s", err)
		}
	}

	defer conn.Close()
	message := "Hello\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Println("client: exiting")
}

var addr string
var enableTls bool

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "server addr")
	flag.BoolVar(&enableTls, "tls", false, "enable tls")

	flag.Parse()
}
