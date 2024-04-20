package validator

import (
	"testing"

	"chatapp.kyriakidis.net/internal/assert"
)

func TestIsEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "Email",
			email: "kws@gmail.com",
			want:  true,
		}, {
			name:  "NotEmail",
			email: "ononafioa.com",
			want:  false,
		},
	}
	for _, tt := range tests {
		isEmail := IsEmail(tt.email)
		assert.Equal(t, isEmail, tt.want)
	}
}
