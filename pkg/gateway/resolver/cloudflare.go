package resolver

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"strconv"
	"tls-dns-proxy/pkg/domain/proxy"

	"golang.org/x/net/dns/dnsmessage"
)

const cflRootCert = `-----BEGIN CERTIFICATE-----
MIIEQzCCAyugAwIBAgIQCidf5wTW7ssj1c1bSxpOBDANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD
QTAeFw0yMDA5MjMwMDAwMDBaFw0zMDA5MjIyMzU5NTlaMFYxCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxMDAuBgNVBAMTJ0RpZ2lDZXJ0IFRMUyBI
eWJyaWQgRUNDIFNIQTM4NCAyMDIwIENBMTB2MBAGByqGSM49AgEGBSuBBAAiA2IA
BMEbxppbmNmkKaDp1AS12+umsmxVwP/tmMZJLwYnUcu/cMEFesOxnYeJuq20ExfJ
qLSDyLiQ0cx0NTY8g3KwtdD3ImnI8YDEe0CPz2iHJlw5ifFNkU3aiYvkA8ND5b8v
c6OCAa4wggGqMB0GA1UdDgQWBBQKvAgpF4ylOW16Ds4zxy6z7fvDejAfBgNVHSME
GDAWgBQD3lA1VtFMu2bwo+IbG8OXsj3RVTAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0l
BBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBIGA1UdEwEB/wQIMAYBAf8CAQAwdgYI
KwYBBQUHAQEEajBoMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2VydC5j
b20wQAYIKwYBBQUHMAKGNGh0dHA6Ly9jYWNlcnRzLmRpZ2ljZXJ0LmNvbS9EaWdp
Q2VydEdsb2JhbFJvb3RDQS5jcnQwewYDVR0fBHQwcjA3oDWgM4YxaHR0cDovL2Ny
bDMuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNybDA3oDWgM4Yx
aHR0cDovL2NybDQuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNy
bDAwBgNVHSAEKTAnMAcGBWeBDAEBMAgGBmeBDAECATAIBgZngQwBAgIwCAYGZ4EM
AQIDMA0GCSqGSIb3DQEBDAUAA4IBAQDeOpcbhb17jApY4+PwCwYAeq9EYyp/3YFt
ERim+vc4YLGwOWK9uHsu8AjJkltz32WQt960V6zALxyZZ02LXvIBoa33llPN1d9R
JzcGRvJvPDGJLEoWKRGC5+23QhST4Nlg+j8cZMsywzEXJNmvPlVv/w+AbxsBCMqk
BGPI2lNM8hkmxPad31z6n58SXqJdH/bYF462YvgdgbYKOytobPAyTgr3mYI5sUje
CzqJx1+NLyc8nAK8Ib2HxnC+IrrWzfRLvVNve8KaN9EtBH7TuMwNW4SpDCmGr6fY
1h3tDjHhkTb9PA36zoaJzu0cIw265vZt6hCmYWJC+/j+fgZwcPwL
-----END CERTIFICATE-----
`

type cloudFlare struct {
	ip       string
	port     int
	rootCert string
}

func NewCloudFlareResolver(ip string, port int) proxy.Resolver {
	return &cloudFlare{ip, port, cflRootCert}
}

func (cfl *cloudFlare) GetTLSConnection() (*tls.Conn, error) {
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(cfl.rootCert)) {
		log.Println("Fail to parse rootCert")
		return nil, errors.New("Fail to parse rootCert")
	}
	dnsCloudFlareConn, err := tls.Dial("tcp", cfl.ip+":"+strconv.Itoa(cfl.port), &tls.Config{
		RootCAs: roots,
	})
	if err != nil {
		log.Println("Error connecting to CloudFlare")
		return nil, err
	}
	return dnsCloudFlareConn, nil
}

func (cfl *cloudFlare) Solve(dnsm dnsmessage.Message) error {
	var b []byte
	b, uerr := dnsm.Pack()
	if uerr != nil {
		log.Println(uerr.Error())
		return errors.New("Could not unpack DNS Message.")
	}
	conn, err := cfl.GetTLSConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	n, rerr := conn.Write(b)
	if rerr != nil {
		log.Println("Error on cloudFlare response")
	}
	dnsm.Unpack(b[:n])
	return nil
}
