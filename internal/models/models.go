package models

type SecretType string

const (
	SecretTypeLogin  SecretType = "login"
	SecretTypeText   SecretType = "text"
	SecretTypeBinary SecretType = "binary"
	SecretTypeCard   SecretType = "card"
)

func (st SecretType) Validate() bool {
	switch st {
	case SecretTypeLogin, SecretTypeText, SecretTypeBinary, SecretTypeCard:
		return true
	}
	return false
}
