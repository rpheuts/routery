package authentication

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"log"
)

type LDAPAuthProvider struct {
	config *AuthConfig
}

func (ap *LDAPAuthProvider) Initialize(config *AuthConfig) error {
	ap.config = config

	return nil
}

func (ap *LDAPAuthProvider) Authenticate(username string, password string) bool {
	return ap.LDAPAuthenticate(username, password)
}

func (ap *LDAPAuthProvider) LDAPAuthenticate(username string, password string) bool {
	bindusername := ap.config.Username
	bindpassword := ap.config.Password

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ap.config.Hostname, ap.config.Port))
	if err != nil {
		log.Println(err)
		return false
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Println(err)
		return false
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		ap.config.Domain,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(ap.config.Arguments, username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println(err)
		return false;
	}

	if len(sr.Entries) != 1 {
		return false
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
