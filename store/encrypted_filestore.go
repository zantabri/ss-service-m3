package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

type enryptedFileStore struct {
	fstore SecretStore
	gcm    cipher.AEAD
	nonce  []byte
}

func encrypt(plain_text string, gcm cipher.AEAD, nonce []byte) (cypher string, err error) {

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	data := gcm.Seal(nonce, nonce, []byte(plain_text), nil)

	cypher = string(data)
	return

}
func decrypt(cyper_text string, gcm cipher.AEAD) (plain_text string, err error) {

	encData := []byte(cyper_text)
	nonce := encData[:gcm.NonceSize()]
	encData = encData[gcm.NonceSize():]
	data, err := gcm.Open(nil, nonce, encData, nil)

	if err != nil {
		return
	}

	plain_text = string(data)

	return

}

func (store *enryptedFileStore) StoreSecret(secret string) string {

	cypher_text, err := encrypt(secret, store.gcm, store.nonce)

	if err != nil {
		panic(err.Error())
	}

	return store.fstore.StoreSecret(cypher_text)
}

func (store *enryptedFileStore) RetriveSecret(id string) string {

	cypher_text := store.fstore.RetriveSecret(id)

	plain_text, err := decrypt(cypher_text, store.gcm)

	if err != nil {
		panic(err.Error())
	}
	
	return plain_text

}

func NewEncryptedFileStore(secret_store SecretStore, salt string, password string) (store SecretStore, err error) {

	key, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())

	store = &enryptedFileStore{fstore: secret_store, gcm: gcm, nonce: nonce}
	return

}
