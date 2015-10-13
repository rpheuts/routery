package authentication

type LDAPAuthProvider struct {
	config *AuthConfig
}

func (ap *LDAPAuthProvider) Initialize(config *AuthConfig) error {
	ap.config = config

	return nil
}

func (ap *LDAPAuthProvider) Authenticate(username string, password string) bool {
	return true
}
