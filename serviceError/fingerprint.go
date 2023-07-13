package serviceError

import (
	"fmt"
	"hash/crc32"
)

func fingerprint(str string) string {
	hash := crc32.NewIEEE()

	fmt.Fprintf(hash, str)

	return fmt.Sprintf("%x", hash.Sum32())
}