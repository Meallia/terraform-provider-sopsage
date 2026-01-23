package provider

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/getsops/sops/v3/config"

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

type EncryptionConfig struct {
	UnencryptedSuffix       string
	EncryptedSuffix         string
	UnencryptedRegex        string
	EncryptedRegex          string
	UnencryptedCommentRegex string
	EncryptedCommentRegex   string
}

func DefaultEncryptionConfig() *EncryptionConfig {
	return &EncryptionConfig{
		UnencryptedSuffix: "_unencrypted",
	}
}

func SopsEncryptDataFromAgeKeys(data string, format string, agePublicKeys []string, encryptionConfig *EncryptionConfig) (string, error) {
	if encryptionConfig == nil {
		encryptionConfig = DefaultEncryptionConfig()
	}
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
			KeyGroups:               []sops.KeyGroup{sf.Map(masterKeys, func(key *keysource.MasterKey) keys.MasterKey { return keys.MasterKey(key) })},
			Version:                 version.Version,
			UnencryptedSuffix:       encryptionConfig.UnencryptedSuffix,
			EncryptedSuffix:         encryptionConfig.EncryptedSuffix,
			UnencryptedRegex:        encryptionConfig.UnencryptedRegex,
			EncryptedRegex:          encryptionConfig.EncryptedRegex,
			UnencryptedCommentRegex: encryptionConfig.UnencryptedCommentRegex,
			EncryptedCommentRegex:   encryptionConfig.EncryptedCommentRegex,
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

	storedSOPSEnv := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "SOPS_") {
			parts := strings.SplitN(env, "=", 2)
			storedSOPSEnv[parts[0]] = parts[1]
			os.Unsetenv(parts[0])
		}
	}

	os.Setenv("SOPS_AGE_KEY", agePrivateKey)
	defer os.Unsetenv("SOPS_AGE_KEY")

	defer func() {
		for key, value := range storedSOPSEnv {
			os.Setenv(key, value)
		}
	}()

	decrypted, err := decrypt.Data([]byte(data), format)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
