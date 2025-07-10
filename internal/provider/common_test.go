package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSopsEncryptDataFromAgeKeys(t *testing.T) {
	testCases := []struct {
		name          string
		data          string
		format        string
		agePublicKeys []string
		wantErr       bool
	}{
		{
			name:          "valid json encryption",
			data:          `{"test": "value"}`,
			format:        "json",
			agePublicKeys: []string{agePubkey},
			wantErr:       false,
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
			result, err := SopsEncryptDataFromAgeKeys(tc.data, tc.format, tc.agePublicKeys)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "sops")
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
			encrypted, err := SopsEncryptDataFromAgeKeys(data, "json", []string{agePubkey})
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
