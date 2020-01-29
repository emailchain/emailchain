package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
)
func ComputeHashSha256(bytes []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(bytes))
}

func GenerateNodeID(host ,port string ) string {
	return fmt.Sprintf("%s-%s",host,port)
}
func DBExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}