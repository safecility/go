package lib

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/rs/zerolog/log"
)

// NewTlsConfig - we use the x509 certPool for a connection to TTN but one could be loaded if necessary:
// certPool := x509.NewCertPool()
// ca, err := os.ReadFile("ca.pem")
//
//	if err != nil {
//		log.Fatalln(err.Error())
//	}
//
// certPool.AppendCertsFromPEM(ca)
//
// clientKeyPair, err := tls.LoadX509KeyPair("client-crt.pem", "client-key.pem")
func NewTlsConfig() *tls.Config {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load system cert pool")
	}
	return &tls.Config{
		RootCAs:            certPool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		//Certificates:       []tls.Certificate{clientKeyPair},
	}
}
