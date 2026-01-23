package provider

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSopsEncryptDataFromAgeKeys(t *testing.T) {
	testCases := []struct {
		name                    string
		data                    string
		format                  string
		agePublicKeys           []string
		wantErr                 bool
		encryptionConfig        *EncryptionConfig
		expectedUnEncryptedKeys []string
	}{
		{
			name:                    "valid json encryption",
			data:                    `{"test": "value", "test_unencrypted": "value"}`,
			format:                  "json",
			agePublicKeys:           []string{agePubkey},
			wantErr:                 false,
			expectedUnEncryptedKeys: []string{"test_unencrypted"},
		},
		{
			name:                    "valid json encryption with unencrypted_regex",
			data:                    `{"test": "value", "_foo": "bar"}`,
			format:                  "json",
			agePublicKeys:           []string{agePubkey},
			wantErr:                 false,
			encryptionConfig:        &EncryptionConfig{UnencryptedRegex: "^_.*"},
			expectedUnEncryptedKeys: []string{"_foo"},
		},
		{
			name:          "invalid age key",
			data:          `{"test": "value"}`,
			format:        "json",
			agePublicKeys: []string{"invalid-key"},
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SopsEncryptDataFromAgeKeys(tc.data, tc.format, tc.agePublicKeys, tc.encryptionConfig)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "sops")

			parsedInput := make(map[string]interface{})
			err = json.Unmarshal([]byte(tc.data), &parsedInput)
			if err != nil {
				t.Log("could not parse json input")
				t.Fatal(err)
			}

			parsedOutput := make(map[string]interface{})
			err = json.Unmarshal([]byte(result), &parsedOutput)
			if err != nil {
				t.Log("could not parse json output")
				t.Fatal(err)
			}

			for _, key := range tc.expectedUnEncryptedKeys {
				assert.Contains(t, parsedInput, key)
				assert.Contains(t, parsedOutput, key)
				assert.Equal(t, parsedInput[key], parsedOutput[key])
			}
		})
	}
}

func TestSopsDecryptDataFromAgeKey(t *testing.T) {
	testCases := []struct {
		name          string
		encrypted     string
		format        string
		agePrivateKey string
		want          string
		wantErr       bool
	}{
		{
			name:          "valid json decryption",
			format:        "json",
			agePrivateKey: agePrivkey,
			wantErr:       false,
		},
		{
			name:          "invalid private key",
			format:        "json",
			agePrivateKey: "invalid-key",
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := `{"test": "value"}`
			encrypted, err := SopsEncryptDataFromAgeKeys(data, "json", []string{agePubkey}, nil)
			assert.NoError(t, err)
			result, err := SopsDecryptDataFromAgeKey(encrypted, tc.format, tc.agePrivateKey)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.JSONEq(t, data, result)
		})
	}
}
