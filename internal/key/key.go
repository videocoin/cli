package key

import (
	"fmt"

	"github.com/VideoCoin/common/bcops"
	"github.com/VideoCoin/go-videocoin/accounts/keystore"
)

type keyStore struct {
	key *keystore.Key
}

func NewKeyStore() *keyStore {
	return &keyStore{}
}

func (k *keyStore) ImportKey(filepath, password string) (*keystore.Key, error) {
	key, err := bcops.LoadBcPrivKeys(filepath, password)
	if err != nil {
		return nil, fmt.Errorf("failed to load blockchain private keys: %s", err.Error())
	}

	k.key = key
	return key, nil
}

func (k *keyStore) GetKey() *keystore.Key {
	return k.key
}
