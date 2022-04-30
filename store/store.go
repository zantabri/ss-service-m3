package store

type SecretStore interface {

	StoreSecret(string) string

	RetriveSecret(string) string

}
