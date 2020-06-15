package password_test

import (
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	. "github.com/my-cargonaut/cargonaut/pkg/password"
)

const (
	key      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567" // 33 byte
	password = "mysecretpassword"
)

func TestGenerateCompare(t *testing.T) {
	type args struct {
		pepper   string
		password string
		cost     int
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			"simple password",
			args{key[:16], "1234567890", bcrypt.MinCost},
			nil,
		},
		{
			"complex password",
			args{key[:16], "987nq459ad8sfmMK()=§?4le9w.Ö0243CL..ÖWL)$§(CM;:::XW§X;∞˘›‹^◊ﬁ", bcrypt.MinCost},
			nil,
		},
		{
			"short password",
			args{key[:16], "0", bcrypt.MinCost},
			nil,
		},
		{
			"long password",
			args{key[:16], "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", bcrypt.MinCost},
			nil,
		},
		{
			"no password",
			args{key[:16], "", bcrypt.MinCost},
			nil,
		},
		{
			"pepper to short",
			args{key[:15], "", bcrypt.MinCost},
			aes.KeySizeError(15),
		},
		{
			"pepper to long",
			args{key[:17], "", bcrypt.MinCost},
			aes.KeySizeError(17),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Generate([]byte(tc.args.pepper), tc.args.password, tc.args.cost)
			require.Equal(t, tc.err, err)

			err = Compare([]byte(tc.args.pepper), tc.args.password, got)
			require.Equal(t, tc.err, err)
		})
	}
}

// TestGenerate_Unique validates that even the same plaintext passwords won't
// result in the same hash. This verifies the bcrypts unique per-user hash works
// as intended. This test uses AES-128.
func TestGenerate_Unique(t *testing.T) {
	hashedPassword1, err := Generate([]byte(key[:16]), password, bcrypt.MinCost)
	require.NoError(t, err)

	hashedPassword2, err := Generate([]byte(key[:16]), password, bcrypt.MinCost)
	require.NoError(t, err)

	assert.NotEqual(t, hashedPassword1, hashedPassword2)
}

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Generate([]byte(key[:16]), password, bcrypt.MinCost)
		require.NoError(b, err)
	}
}

func BenchmarkCompare(b *testing.B) {
	hash, err := Generate([]byte(key[:16]), password, bcrypt.MinCost)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = Compare([]byte(key[:16]), password, hash)
		require.NoError(b, err)
	}
}
