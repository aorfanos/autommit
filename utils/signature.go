package utils

import (
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

func (a *Autommit) GetOpenPGPKeyring() (err error) {
	// open the keyring
	keyringFile, err := os.Open(a.PgpKeyPath)
	if err != nil {
		return err
	}
	defer keyringFile.Close()

	// Read the armored keyring
	block, err := armor.Decode(keyringFile)
	if err != nil {
		return err
	}

	// get the keyring
	keyring, err := openpgp.ReadKeyRing(block.Body)
	if err != nil {
		return err
	}

	a.GitConfig.PGPKeyRing = keyring
	return err
}
