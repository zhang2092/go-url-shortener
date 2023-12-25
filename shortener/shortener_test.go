package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const userId = "7c729139-7ff4-445f-976b-2d842f55cb0e"

func TestGenerateShortLink(t *testing.T) {
	link1 := "https://www.baidu.com/"
	short1, err := GenerateShortLink(link1, userId)
	assert.NoError(t, err)
	assert.Equal(t, short1, "egtq236P")

	link2 := "https://www.163.com/"
	short2, err := GenerateShortLink(link2, userId)
	assert.NoError(t, err)
	assert.Equal(t, short2, "DiCqg9Yp")

	link3 := "https://www.qq.com/"
	short3, err := GenerateShortLink(link3, userId)
	assert.NoError(t, err)
	assert.Equal(t, short3, "4QhQ62cZ")
}
