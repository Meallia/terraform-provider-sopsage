package provider

import (
	"fmt"
	"github.com/getsops/sops/v3/config"
	"os"
	"strings"
	"sync"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	keysource "github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/cmd/sops/formats"
	"github.com/getsops/sops/v3/decrypt"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/version"
	sf "github.com/sa-/slicefunk"
)

func SopsEncryptDataFromAgeKeys(data string, format string, agePublicKeys []string) (string, error) {
	store := common.StoreForFormat(
		formats.FormatFromString(format),
		config.NewStoresConfig(),
	)

	branches, err := store.LoadPlainFile([]byte(data))
	if err != nil {
		return "", err
	}

	masterKeys, err := keysource.MasterKeysFromRecipients(strings.Join(agePublicKeys, ","))
	if err != nil {
		return "", err
	}

	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			KeyGroups:         []sops.KeyGroup{sf.Map(masterKeys, func(key *keysource.MasterKey) keys.MasterKey { return keys.MasterKey(key) })},
			UnencryptedSuffix: "_unencrypted",
			Version:           version.Version,
		},
	}

	dataKey, errs := tree.GenerateDataKeyWithKeyServices(
		[]keyservice.KeyServiceClient{keyservice.NewLocalClient()},
	)
	if len(errs) > 0 {
		return "", fmt.Errorf("could not use provider age keys for encryption: %s", errs)
	}

	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  aes.NewCipher(),
	})

	if err != nil {
		return "", err
	}
	result, err := store.EmitEncryptedFile(tree)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

var decryptMutex = sync.Mutex{}

func SopsDecryptDataFromAgeKey(data string, format string, agePrivateKey string) (string, error) {
	decryptMutex.Lock()
	defer decryptMutex.Unlock()

	originalKey := os.Getenv("SOPS_AGE_KEY")
	err := os.Setenv("SOPS_AGE_KEY", agePrivateKey)
	if err != nil {
		return "", err
	}
	defer os.Setenv("SOPS_AGE_KEY", originalKey)

	decrypted, err := decrypt.Data([]byte(data), format)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
