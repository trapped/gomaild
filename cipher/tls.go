package cipher

import (
	"crypto/tls"
	"github.com/trapped/gomaild/config"
	"log"
	"net"
)

var TLSAvailable bool
var TLSConfig *tls.Config

func TLSLoadCertificate() {
	cert, err := tls.LoadX509KeyPair(config.Configuration.TLS.CertificateFile,
		config.Configuration.TLS.CertificateKeyFile)
	if err != nil {
		log.Println("TLS:", "Failed loading SSL certificate:", err)
		return
	}
	TLSAvailable = true
	TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
	}
}

func TLSTransmuteConn(c net.Conn) net.Conn {
	tc := tls.Server(c, TLSConfig)
	tc.Handshake()
	return net.Conn(tc)
}
