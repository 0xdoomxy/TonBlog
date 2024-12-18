package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcrecover(t *testing.T) {
	address, err := Ecrecover(
		"{\"types\":{\"EIP712Domain\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"version\",\"type\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\"}],\"Verify\":[{\"name\":\"content\",\"type\":\"string\"},{\"name\":\"date\",\"type\":\"uint256\"}]},\"domain\":{\"chainId\":\"0x1\",\"name\":\"0xdoomxy blog\",\"version\":\"1\"},\"primaryType\":\"Verify\",\"message\":{\"content\":\"Welcome to 0xdoomxy blog\",\"date\":1734495476259}}",
		"0x2ff2700432235806f53a25b12b7593810a040d5d31fe52be720593b3563828536972b99e7f6c971eaf8ca7d8f5edf4965ca904b9565742c5be3d5d592f946b1a1b")
	assert.Equal(t, "0x7BceBBF3E62dcFfEda814866e5A9088E0423F1d3", address, err)
}
