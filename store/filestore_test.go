package store

import (
	"os"
	"testing"

)

var sd string

func setup() {

	sd = os.TempDir()

}

func tearDown() {

	os.Remove(sd + "/" + DEFAULT_FILE_NAME)


}

func TestMain(m *testing.M) {

	setup()
	exitCode := m.Run()
	tearDown()

	os.Exit(exitCode)

}

func TestNewFileStore(t *testing.T) {

	store, err := NewFileStore(&sd)

	if err != nil {
		t.Fatal(err.Error())
	}

	if store == nil {
		t.Fatal("store is nil")
	}
	
}

func TestStoreAndRetrieveSecret(t *testing.T) {


	secret := "lubalu2323232balu"
	store, _ := NewFileStore(&sd)

	id := store.StoreSecret(secret)
	
	if len(id) == 0 {
		t.Error("id is nil")
	}

	if id != "1f2b78f8b8067dfd47df852f12697c69" {
		t.Errorf("id is %s", id)
	}

	returnedSecret := store.RetriveSecret(id)
	
	if secret != returnedSecret {
		t.Errorf("returned secert %s", returnedSecret)
	}

}

func TestStoreSecretWith2XRetrieve(t *testing.T) {

	secret := "dxcuPj99gKMLzRZz"
	store, _ := NewFileStore(&sd)

	id := store.StoreSecret(secret)

	store.RetriveSecret(id)
	second := store.RetriveSecret(id)


	if len(second) != 0 {
		t.Errorf("second retrieval %s", second)
	}

}
