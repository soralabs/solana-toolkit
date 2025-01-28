package gmgn

import tls_client "github.com/bogdanfinn/tls-client"

func New(tlsClient tls_client.HttpClient) *GMGN {
	return &GMGN{tlsClient: tlsClient}
}
