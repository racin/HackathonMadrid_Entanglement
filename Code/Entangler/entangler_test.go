package Entangler

import (
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMemoryPosition(t *testing.T) {
	var r, h, l int
	// Datablock 1
	r, h, l = GetMemoryPosition(1)
	assert.Equal(t, 0, r, "Datablock 1: R position should be 0")
	assert.Equal(t, 5, h, "Datablock 1: H position should be 5")
	assert.Equal(t, 11, l, "Datablock 1: L position should be 11")

	// Datablock 2
	r, h, l = GetMemoryPosition(2)
	assert.Equal(t, 4, r, "Datablock 2: R position should be 4")
	assert.Equal(t, 6, h, "Datablock 2: H position should be 6")
	assert.Equal(t, 12, l, "Datablock 2: L position should be 12")

	// Datablock 3
	r, h, l = GetMemoryPosition(3)
	assert.Equal(t, 3, r, "Datablock 3: R position should be 3")
	assert.Equal(t, 7, h, "Datablock 3: H position should be 7")
	assert.Equal(t, 13, l, "Datablock 3: L position should be 13")

	// Datablock 4
	r, h, l = GetMemoryPosition(4)
	assert.Equal(t, 2, r, "Datablock 4: R position should be 2")
	assert.Equal(t, 8, h, "Datablock 4: H position should be 8")
	assert.Equal(t, 14, l, "Datablock 4: L position should be 14")

	// Datablock 5
	r, h, l = GetMemoryPosition(5)
	assert.Equal(t, 1, r, "Datablock 5: R position should be 1")
	assert.Equal(t, 9, h, "Datablock 5: H position should be 9")
	assert.Equal(t, 10, l, "Datablock 5: L position should be 10")

	// Datablock 21
	r, h, l = GetMemoryPosition(21)
	assert.Equal(t, 4, r, "Datablock 21: R position should be 4")
	assert.Equal(t, 5, h, "Datablock 21: H position should be 5")
	assert.Equal(t, 10, l, "Datablock 21: L position should be 10")

	// Datablock 22
	r, h, l = GetMemoryPosition(22)
	assert.Equal(t, 3, r, "Datablock 22: R position should be 3")
	assert.Equal(t, 6, h, "Datablock 22: H position should be 6")
	assert.Equal(t, 11, l, "Datablock 22: L position should be 11")

	// Datablock 23
	r, h, l = GetMemoryPosition(23)
	assert.Equal(t, 2, r, "Datablock 23: R position should be 2")
	assert.Equal(t, 7, h, "Datablock 23: H position should be 7")
	assert.Equal(t, 12, l, "Datablock 23: L position should be 12")

	// Datablock 24
	r, h, l = GetMemoryPosition(24)
	assert.Equal(t, 1, r, "Datablock 24: R position should be 1")
	assert.Equal(t, 8, h, "Datablock 24: H position should be 8")
	assert.Equal(t, 13, l, "Datablock 24: L position should be 13")

	// Datablock 25
	r, h, l = GetMemoryPosition(25)
	assert.Equal(t, 0, r, "Datablock 25: R position should be 0")
	assert.Equal(t, 9, h, "Datablock 25: H position should be 9")
	assert.Equal(t, 14, l, "Datablock 25: L position should be 14")
}

func TestEntangleFile(t *testing.T) {
	filePath := "../../resources/images/swarm_logo.jpeg"
	EntangleFile(filePath)

}

func TestEntangler(t *testing.T) {
	// check hash of data block
	data, _ := ReadChunk("../../testing/d8")
	p1, _ := ReadChunk("../../testing/p4_8")
	p2, _ := ReadChunk("../../testing/p8_12")

	// Xor two parities
	rData, _ := XORByteSlice(p1, p2)

	// Compare hash
	dataHash := sha256.Sum256(data)
	parityHash := sha256.Sum256(rData)
	dataHashStr := fmt.Sprintf("%x", dataHash)
	parityHashStr := fmt.Sprintf("%x", parityHash)

	fmt.Println(dataHashStr)
	fmt.Println(parityHashStr)

	assert.Equal(t, dataHashStr, parityHashStr, "Hash should be equal")
}

func TestXOR(t *testing.T) {
	a := []byte{0xb9, 0x63}
	b := []byte{0x66, 0xcc}

	// Xor two parities
	c, _ := XORByteSlice(a, b)

	// Compare hash
	cStr := fmt.Sprintf("%x", c)

	fmt.Println(c)

	assert.Equal(t, "dfaf", cStr, "Hash should be equal")
}
