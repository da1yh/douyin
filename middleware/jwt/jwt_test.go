package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token := GenerateToken(3, "aaa", "qqq")
	assert.True(t, len(token) > 0)
}

func TestParseToken(t *testing.T) {
	token := GenerateToken(4, "qqq", "ppp")
	assert.True(t, len(token) > 0)
	claims, err := ParseToken(token)
	assert.True(t, claims.Name == "qqq" && claims.Password == "ppp" && err == nil)
}
