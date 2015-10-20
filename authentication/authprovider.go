package authentication
import (
	"github.com/rpheuts/routery/config"
	"log"
)

type AuthConfig struct {
	Enabled bool
	Hostname string
	Port int
	Arguments string
	Username string
	Password string
}

type Provider interface {
	Initialize(config *AuthConfig) error
	Authenticate(username string, password string) error
}

func Authenticate(routeryConfig *config.RouteryConfig, username string, password string) bool {

	retval := false

	// Iterate over Auth providers
	for _, authConfig := range routeryConfig.Auth {
		if authConfig.Type == "LDAP" {
			ldap  := LDAPAuthProvider{}
			ldap.Initialize(&AuthConfig{true, authConfig.Hostname, authConfig.Port, authConfig.Arguments, authConfig.Username, authConfig.Password})

			log.Println("Trying to provide auth using: LDAP")
			retval = ldap.Authenticate(username, password)
		}
	}

	return retval
}