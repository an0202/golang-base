/**
 * @Author: jie.an
 * @Description:
 * @File:  hash.go
 * @Version: 1.0.0
 * @Date: 2020/3/22 1:23 下午
 */

package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(data string) (HashedString string) {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
