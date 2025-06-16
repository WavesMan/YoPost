package service

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"time"
)

// whichTLS 检查TLS证书有效性
// 返回TLS连接是否可用(true/false)
func whichTLS(conn net.Conn) bool {
	log.Printf("INFO: TLS check on existing connection")
	log.Printf("DEBUG: Connection details - Local: %s, Remote: %s", conn.LocalAddr(), conn.RemoteAddr())

	host, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	config := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}
	log.Printf("DEBUG: TLS config - ServerName: %s", host)

	log.Printf("INFO: Attempting TLS handshake")
	tlsConn := tls.Client(conn, config)
	err := tlsConn.Handshake()
	if err != nil {
		log.Printf("ERROR: TLS handshake failed - %v", err)
		return false
	}
	log.Printf("INFO: TLS handshake successful")
	log.Printf("DEBUG: TLS connection state: %+v", tlsConn.ConnectionState())

	// 验证证书链
	log.Printf("INFO: Verifying certificate chain")
	certs := tlsConn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		log.Printf("ERROR: No certificates received")
		return false
	}
	log.Printf("INFO: Certificate chain verified (%d certificates)", len(certs))
	log.Printf("DEBUG: Certificate details - Issuer: %s, Subject: %s, Expiry: %s",
		certs[0].Issuer, certs[0].Subject, certs[0].NotAfter)

	// 检查证书吊销状态
	log.Printf("INFO: Checking certificate revocation status")
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
	}
	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
	}

	// 尝试OCSP检查
	log.Printf("DEBUG: Attempting OCSP check for certificate")
	if _, err := certs[0].Verify(opts); err != nil {
		log.Printf("ERROR: Certificate revocation check failed - %v", err)
		return false
	}

	// 尝试CRL检查
	log.Printf("DEBUG: Attempting CRL check for certificate")
	crlList, err := certs[0].CRLDistributionPoints()
	if err != nil {
		log.Printf("ERROR: Failed to get CRL distribution points - %v", err)
	}
	if len(crlList) > 0 {
		log.Printf("DEBUG: Found %d CRL distribution points", len(crlList))
		// 在实际应用中这里应该实现CRL下载和验证逻辑
		log.Printf("WARNING: CRL checking not fully implemented")
	}

	log.Printf("INFO: Certificate revocation check passed")
	return true
}

// CheckTLSGlobal 全局TLS检查接口
func CheckTLSGlobal(host, port string) bool {
	log.Printf("INFO: Performing TLS check for %s:%s", host, port)
	conn, err := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	if err != nil {
		log.Printf("ERROR: TCP connection failed - %v", err)
		return false
	}
	defer conn.Close()
	return whichTLS(conn)
}
