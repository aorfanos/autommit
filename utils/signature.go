package utils

// import (
// 	"os"

// 	"github.com/ProtonMail/go-crypto/openpgp"
// 	"github.com/ProtonMail/go-crypto/openpgp/armor"
// 	"github.com/ProtonMail/go-crypto/openpgp/packet"
// )

// func GetOpenPGPKeyring(path string) (keyring openpgp.KeyRing, err error) {
// 	// open the keyring
// 	keyringFile, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer keyringFile.Close()

// 	// Read the armored keyring
// 	block, err := armor.Decode(keyringFile)
// 	ErrCheck(err)

// 	// get the keyring
// 	keyring, err = openpgp.ReadKeyRing(packet.NewReader(block.Body))
// 	ErrCheck(err)

// 	return keyring, nil
// }