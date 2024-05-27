package tool

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	url2 "net/url"
	"sort"
	"strings"
)

// DEC-CBC加密
func ApiEncrypt(data, key, iv string) (string, error) {
	if len(key) != 8 {
		return "", errors.New("密钥必须8位长度")
	}
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	data = string(PKCS5Padding([]byte(data), bs))
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	out := make([]byte, len([]byte(data)))
	blockMode.CryptBlocks(out, []byte(data))

	return string(out), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte(fmt.Sprintf("%c", padding)), padding)
	ciphertext = append(ciphertext, padText...)
	return ciphertext
}

// DEC-CBC解密
func ApiDecrypt(data, key, iv string) (string, error) {
	byteStr, _ := base64.StdEncoding.DecodeString(data)
	data = string(byteStr)
	if len(key) != 8 {
		return "", errors.New("密钥必须8位长度")
	}
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	out := make([]byte, len([]byte(data)))
	blockMode.CryptBlocks(out, []byte(data))
	out = PKCS5UnPadding(out)
	return string(out), nil
}

func PKCS5UnPadding(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}

func MakeSign(data map[string]interface{}, key string) string {
	var keys []string
	for k, _ := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortData := make(map[string]interface{})
	for _, k := range keys {
		sortData[k] = data[k]
	}
	url := ToUrlParams(sortData, keys)
	h := md5.New()
	str := url + "&key=" + key
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifySign(data url2.Values, key string) bool {
	var keys []string
	//判断是否有sign字段
	_, ok := data["sign"]
	if !ok {
		return false
	}
	sign := cast.ToString(data["sign"][0])
	delete(data, "sign")
	for k, _ := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sortData map[string]interface{}
	sortData = make(map[string]interface{})
	for _, k := range keys {
		sortData[k] = data[k][0]
	}
	url := ToUrlParams(sortData, keys)
	h := md5.New()
	str := url + "&key=" + key
	h.Write([]byte(str))
	makeSign := hex.EncodeToString(h.Sum(nil))
	return makeSign == sign
}

func base64Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func base64Decode(src string) string {
	reader := strings.NewReader(src)
	decoder := base64.NewDecoder(base64.RawStdEncoding, reader)
	buf := make([]byte, 1024)
	dst := ""
	for {
		n, err := decoder.Read(buf)
		dst += string(buf[:n])
		if n == 0 || err != nil {
			break
		}
	}

	return dst
}
