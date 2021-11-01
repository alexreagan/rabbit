package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
)

func NewTLSConfig(ldapCAchain string) *tls.Config {
	// Load client cert and key
	//cert, err := tls.LoadX509KeyPair(ldapCert, ldapKey)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Load CA chain
	caCert, err := ioutil.ReadFile(ldapCAchain)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup TLS with ldap client cert
	tlsConfig := &tls.Config{
		//Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	return tlsConfig
}

func Authenticate(username string, password string) error {
	// ldap
	tlsConfig := NewTLSConfig(viper.GetString("ldap.ad_certfile"))
	conn, err := ldap.DialTLS("tcp", viper.GetString("ldap.ad_url"), tlsConfig)
	if err != nil {
		return err
	}
	err = conn.Bind(viper.GetString("ldap.ad_dn"), viper.GetString("ldap.ad_pw"))
	if err != nil {
		return err
	}
	defer conn.Unbind()

	searchRequest := ldap.NewSearchRequest(
		viper.GetString("ldap.ad_domain"), // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=user)(objectCategory=person)(sAMAccountName=%s))", username), // 查询所有
		[]string{"cn", "sAMAccountName"}, // A list attributes to retrieve
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		log.Error(err)
		return err
	}
	userDn := sr.Entries[0].DN
	log.Debugf("ldap response: %#v", sr)
	log.Debugf("userDn: %s", userDn)
	log.Debugf("password: %s", password)

	err = conn.Bind(userDn, password)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
