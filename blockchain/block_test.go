package blockchain

import (
	"bytes"
	"testing"
)

func TestComputeHashForBlock(t *testing.T) {
	// genesis block
	block := Block{
		Height:       0,
		Timestamp:    0,
		Emails:       nil,
		PreviousHash: "0",
	}
	expectedHash := "d78eb8f6c461dd6d740095e76fb3c6b4ed7ecbfd830126bc14837b80871399f0"
	actualHash := ComputeHashForBlock(block)
	if expectedHash != actualHash {
		t.Errorf("Expected %s but got %s", expectedHash, actualHash)
	}
}
func TestBlock_Serialize(t *testing.T) {
	block := Block{
		Height:       0,
		Timestamp:    0,
		Emails:       nil,
		PreviousHash: "0",
	}
	expectedValue := []byte{123, 34, 104, 101, 105, 103, 104, 116, 34, 58, 48, 44, 34, 116, 105, 109, 101, 115, 116, 97, 109, 112, 34, 58, 48, 44, 34, 101, 109, 97, 105, 108, 115, 34, 58, 110, 117, 108, 108, 44, 34, 112, 114, 101, 118, 105, 111, 117, 115, 95, 104, 97, 115, 104, 34, 58, 34, 48, 34, 44, 34, 110, 111, 110, 99, 101, 34, 58, 48, 125}
	actualValue := block.Serialize()
	if bytes.Compare(expectedValue, actualValue) != 0 {
		t.Fail()
	}
}
func TestDeserializeBlock(t *testing.T) {
	input := []byte{123, 34, 104, 101, 105, 103, 104, 116, 34, 58, 48, 44, 34, 116, 105, 109, 101, 115, 116, 97, 109, 112, 34, 58, 48, 44, 34, 101, 109, 97, 105, 108, 115, 34, 58, 110, 117, 108, 108, 44, 34, 112, 114, 101, 118, 105, 111, 117, 115, 95, 104, 97, 115, 104, 34, 58, 34, 48, 34, 44, 34, 110, 111, 110, 99, 101, 34, 58, 48, 125}
	expectedOutput := Block{
		Height:       0,
		Timestamp:    0,
		Emails:       nil,
		PreviousHash: "0",
	}
	actualOutput := *DeserializeBlock(input)
	if actualOutput.Height != expectedOutput.Height || actualOutput.Timestamp != expectedOutput.Timestamp || actualOutput.PreviousHash != expectedOutput.PreviousHash || len(actualOutput.Emails) != len(expectedOutput.Emails) {
		t.Fail()
	}
	for i, v := range expectedOutput.Emails {
		if actualOutput.Emails[i] != v {
			t.Fail()
		}
	}
	for k, v := range actualOutput.Emails {
		if expectedOutput.Emails[k] != v {
			t.Fail()
		}
	}
}
