package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func sign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
	}
	sb.WriteString("&key=")
	sb.WriteString(key)
	h := md5.Sum([]byte(sb.String()))
	return hex.EncodeToString(h[:])
}

func main() {
	params := map[string]string{
		"money":       "0.01",
		"name":        "压测商品",
		"notify_url":  "http://localhost:8080/api/v1/pay/epay/notify",
		"out_trade_no": "995790054E310906",
		"pid":         "test_mch_id",
		"return_url":  "http://localhost:8080/order/995790054E310906",
		"type":        "1",
	}
	key := "test_key_for_signature"

	result := sign(params, key)
	fmt.Println("Go signature:", result)

	// 模拟 form Values
	form := make(url.Values)
	for k, v := range params {
		form.Set(k, v)
	}

	// 模拟 VerifyNotify（去除 sign/sign_type）
	verifyParams := make(map[string]string)
	for k, v := range form {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if len(v) > 0 {
			verifyParams[k] = v[0]
		}
	}

	s2 := sign(verifyParams, key)
	fmt.Println("Verify signature:", s2)
	fmt.Println("Match:", result == s2)
}
