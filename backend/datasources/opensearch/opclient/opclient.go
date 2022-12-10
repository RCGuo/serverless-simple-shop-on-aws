package opclient

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go"
)

var Client *opensearch.Client

func init(){
  client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{os.Getenv("OPENSEARCH_DOMAIN")},
    Username: os.Getenv("OPENSEARCH_USER"),
    Password: os.Getenv("OPENSEARCH_PASSWORD"),
    Transport: &http.Transport{
      MaxIdleConnsPerHost:   10,
      ResponseHeaderTimeout: time.Second,
      DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
      TLSClientConfig:       &tls.Config{
        MinVersion:          tls.VersionTLS11,
      },
    },
	})

	if err != nil {
		log.Printf("cannot initialize opensearch client: %s", err.Error())
	}
  Client = client
}