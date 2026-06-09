package service

import (
	"chenze-faka/internal/pkg/config"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type EpayService struct{}

func NewEpayService() *EpayService { return &EpayService{} }

// BuildRedirectURL 构造易支付跳转支付URL
func (s *EpayService) BuildRedirectURL(orderNo string, amount float64, productName, returnURL, notifyURL string) string {
	cfg := config.AppConfig.Pay.Epay
	params := map[string]string{
		"pid":          cfg.MchID,
		"out_trade_no": orderNo,
		"name":         productName,
		"money":        fmt.Sprintf("%.2f", amount),
		"notify_url":   notifyURL,
		"return_url":   returnURL,
	}
	sign := s.sign(params, cfg.Key)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	u, _ := url.Parse(cfg.APIURL + "/submit.php")
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// sign 生成 MD5 签名（参数按字母排序后拼接 &key=KEY）
func (s *EpayService) sign(params map[string]string, key string) string {
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

// VerifyNotify 验证回调签名
func (s *EpayService) VerifyNotify(form url.Values) (orderNo string, tradeNo string, ok bool) {
	key := config.AppConfig.Pay.Epay.Key
	params := make(map[string]string)
	for k, v := range form {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	inSign := form.Get("sign")
	outSign := s.sign(params, key)
	if inSign != outSign {
		return "", "", false
	}
	if form.Get("trade_status") != "TRADE_SUCCESS" {
		return "", "", false
	}
	return form.Get("out_trade_no"), form.Get("trade_no"), true
}

// NotifyOK 返回给易支付的成功响应
func (s *EpayService) NotifyOK(w io.Writer) {
	fmt.Fprint(w, "success")
}

func HTTPGet(u string) (string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
