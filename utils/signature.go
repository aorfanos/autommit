package utils

import (
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

func (a *Autommit) GetOpenPGPKeyring() (err error) {
	// open the keyring
	keyringFile, err := os.Open(a.PgpKeyPath)
	ErrCheck(err)
	defer keyringFile.Close()

	// Read the armored keyring
	block, err := armor.Decode(keyringFile)
	ErrCheck(err)

	// get the keyring
	keyring, err := openpgp.ReadKeyRing(block.Body)
	ErrCheck(err)

	a.GitConfig.PGPKeyRing = keyring
	return err
}
