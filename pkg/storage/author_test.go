package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_SignatureString(t *testing.T) {
	s := Signature{
		Name:  "Foo Bar",
		Email: "foo@bar.com",
		Date:  time.Unix(1505935797, 0).In(time.FixedZone("", -25200)),
	}

	assert.Equal(t, s.String(), "Foo Bar <foo@bar.com> 1505935797 -0700")
}
func Test_parseSignature(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    Signature
		wantErr bool
	}{
		{"empty", "", Signature{}, true},
		{"no TZ", "First Lastname <first.lastname@example.com> 1505935797", Signature{}, true},
		{"no email", "Foo Bar <> 1505925797", Signature{}, true},
		{"no name", "<foo@bar.com> 1505925797 -0700", Signature{}, true},
		{"valid", "First Lastname <first.lastname@example.com> 1505935797 -0700", Signature{
			Name:  "First Lastname",
			Email: "first.lastname@example.com",
			Date:  time.Unix(1505935797, 0).In(time.FixedZone("", -25200)),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSignature(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, got.Name, tt.want.Name)
			assert.Equal(t, got.Email, tt.want.Email)
			assert.Equal(t, got.Date, tt.want.Date)
		})
	}
}
