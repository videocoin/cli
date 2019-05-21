package key

import (
	"io/ioutil"
	"testing"
)

const testKey = `{"version":3,"id":"fd60ff00-d82f-4853-86de-e713ecc70dcc","address":"b4644d5723368ca1f20613935663e50f683efb12","crypto":{"ciphertext":"6675cdbe2001adbd156810cd23c052c8b5bce19d4f2e5837df17950990beef12","cipherparams":{"iv":"8b08399bfea4e646dfe4e91b76f4a2b1"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"6dbfe1dcefe615293afd1b1fdeab2311021c9cf80bd4056ab43b8458a89d7cd0","n":131072,"r":8,"p":1},"mac":"5f48e9f6aa8040dab77680627006d42776957641db6e7769fad3940073ad8fd5"}}`
const testKeyPwd = "EjbY9HBqeTXKeKe4FBJWZQLNfuKp3j"

func TestImportKey(t *testing.T) {
	tables := []struct {
		key      string
		password string
		result   bool
	}{
		{"", "", false},
		{testKey, "", false},
		{testKey, testKeyPwd, true},
	}

	filepath := "/tmp/key"

	ks := new(keyStore)
	for i, table := range tables {
		err := ioutil.WriteFile(filepath, []byte(table.key), 0644)
		if err != nil {
			t.Errorf("Test %d ImportKey failed with err: %s", i, err)
			continue
		}

		key, err := ks.ImportKey(filepath, table.password)
		if err != nil && table.result == true {
			t.Errorf("Test %d ImportKey failed with err: %s", i, err)
			continue
		}

		if key == nil && table.result == true {
			t.Errorf("Test %d ImportKey failed", i)
			continue
		}
	}
}

func TestGetKey(t *testing.T) {
	tables := []struct {
		key      string
		password string
		before   bool
	}{
		{testKey, testKeyPwd, false},
		{testKey, testKeyPwd, true},
	}

	filepath := "/tmp/key"

	ks := new(keyStore)
	for i, table := range tables {
		err := ioutil.WriteFile(filepath, []byte(table.key), 0644)
		if err != nil {
			t.Errorf("Test %d ImportKey failed with err: %s", i, err)
			continue
		}

		if table.before == false {
			key := ks.GetKey()
			if key != nil {
				t.Errorf("Test %d GetKey failed", i)
			}
		} else {
			_, err := ks.ImportKey(filepath, table.password)
			if err != nil {
				t.Errorf("Test %d GetKey failed with err: %s", i, err)
				continue
			}

			key := ks.GetKey()
			if key == nil {
				t.Errorf("Test %d GetKey failed", i)
			}
		}
	}
}
