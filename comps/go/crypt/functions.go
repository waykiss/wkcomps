package crypt

import "encoding/base64"

func StringToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
