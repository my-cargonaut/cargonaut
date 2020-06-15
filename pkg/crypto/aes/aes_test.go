package aes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/my-cargonaut/cargonaut/pkg/crypto/aes"
)

var key = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567" // 33 byte

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{"AES-128", []byte(key[:16])},
		{"AES-192", []byte(key[:24])},
		{"AES-256", []byte(key[:32])},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Encrypt(tc.key, []byte("Secret Message"))
			require.NoError(t, err)

			got, err = Decrypt(tc.key, got)
			require.NoError(t, err)

			assert.Equal(t, "Secret Message", got)
		})
	}
}
