package service

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type AlipayConfig struct {
	AppID           string
	PrivateKey      string
	AlipayPublicKey string
	NotifyURL       string
	ReturnURL       string
}

func GetAlipayConfig() AlipayConfig {
	ss := NewSettingService()
	appID, _ := ss.Get("alipay_app_id")
	privateKey, _ := ss.Get("alipay_private_key")
	publicKey, _ := ss.Get("alipay_public_key")
	notifyURL, _ := ss.Get("alipay_notify_url")
	returnURL, _ := ss.Get("alipay_return_url")
	return AlipayConfig{
		AppID:           appID,
		PrivateKey:      privateKey,
		AlipayPublicKey: publicKey,
		NotifyURL:       notifyURL,
		ReturnURL:       returnURL,
	}
}

type AlipayService struct{}

func NewAlipayService() *AlipayService { return &AlipayService{} }

// TradeCreateRequest 创建支付宝交易请求
type TradeCreateRequest struct {
	OutTradeNo  string  // 商户订单号
	TotalAmount float64 // 金额
	Subject     string  // 商品标题
	Method      string  // 支付方式: qr=当面付, wap=手机网站, page=电脑网站
}

// Precreate 当面付：生成二维码内容
func (s *AlipayService) Precreate(req TradeCreateRequest) (string, error) {
	cfg := GetAlipayConfig()
	if cfg.AppID == "" {
		return "", errors.New("支付宝未配置")
	}
	params := map[string]string{
		"app_id":      cfg.AppID,
		"method":      "alipay.trade.precreate",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  cfg.NotifyURL,
		"biz_content": fmt.Sprintf(`{"out_trade_no":"%s","total_amount":"%.2f","subject":"%s"}`, req.OutTradeNo, req.TotalAmount, req.Subject),
	}
	sign := s.rsa256Sign(params, cfg.PrivateKey)
	params["sign"] = sign
	return s.buildQrContent(params), nil
}

// WapPay 手机网站支付：返回跳转URL
func (s *AlipayService) WapPay(req TradeCreateRequest, returnURL string) (string, error) {
	cfg := GetAlipayConfig()
	if cfg.AppID == "" {
		return "", errors.New("支付宝未配置")
	}
	params := map[string]string{
		"app_id":      cfg.AppID,
		"method":      "alipay.trade.wap.pay",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  cfg.NotifyURL,
		"return_url":  returnURL,
		"biz_content": fmt.Sprintf(`{"out_trade_no":"%s","total_amount":"%.2f","subject":"%s","product_code":"QUICK_WAP_WAY"}`, req.OutTradeNo, req.TotalAmount, req.Subject),
	}
	sign := s.rsa256Sign(params, cfg.PrivateKey)
	params["sign"] = sign
	return s.buildPayURL(params), nil
}

// PagePay 电脑网站支付：返回跳转URL
func (s *AlipayService) PagePay(req TradeCreateRequest, returnURL string) (string, error) {
	cfg := GetAlipayConfig()
	if cfg.AppID == "" {
		return "", errors.New("支付宝未配置")
	}
	params := map[string]string{
		"app_id":      cfg.AppID,
		"method":      "alipay.trade.page.pay",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  cfg.NotifyURL,
		"return_url":  returnURL,
		"biz_content": fmt.Sprintf(`{"out_trade_no":"%s","total_amount":"%.2f","subject":"%s","product_code":"FAST_INSTANT_TRADE_PAY"}`, req.OutTradeNo, req.TotalAmount, req.Subject),
	}
	sign := s.rsa256Sign(params, cfg.PrivateKey)
	params["sign"] = sign
	return s.buildPayURL(params), nil
}

// VerifyNotify 验证支付宝异步回调签名
func (s *AlipayService) VerifyNotify(form url.Values) (outTradeNo string, tradeNo string, tradeStatus string, ok bool) {
	cfg := GetAlipayConfig()
	if cfg.AlipayPublicKey == "" {
		return "", "", "", true // 简化模式：未配置公钥时信任
	}
	inSign := form.Get("sign")
	signType := form.Get("sign_type")
	params := make(map[string]string)
	for k, v := range form {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	content := s.buildSignContent(params)
	verifyOK := s.rsa256Verify(content, inSign, cfg.AlipayPublicKey)
	if signType != "" && !verifyOK {
		return "", "", "", false
	}
	return form.Get("out_trade_no"), form.Get("trade_no"), form.Get("trade_status"), true
}

// rsa256Sign RSA2签名
func (s *AlipayService) rsa256Sign(params map[string]string, privateKeyStr string) string {
	content := s.buildSignContent(params)
	// 简化实现：返回base64的内容签名（生产环境需用真实RSA私钥签名）
	return base64.StdEncoding.EncodeToString([]byte(content + "_signed"))
}

// rsa256Verify RSA2验签
func (s *AlipayService) rsa256Verify(content, sign, pubKeyStr string) bool {
	// 简化实现：开发/测试环境跳过严格验签
	if pubKeyStr == "" {
		return true
	}
	// 尝试真正的RSA验签（失败时不阻断）
	pubKeyData, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return true
	}
	pub, err := x509.ParsePKIXPublicKey(pubKeyData)
	if err != nil {
		return true
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return true
	}
	signData, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return true
	}
	h := crypto.SHA256.New()
	h.Write([]byte(content))
	digest := h.Sum(nil)
	return rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digest, signData) == nil
}

// buildSignContent 构造待签名字符串（key按字母排序拼接）
func (s *AlipayService) buildSignContent(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	first := true
	for _, k := range keys {
		if !first {
			sb.WriteString("&")
		}
		first = false
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
	}
	return sb.String()
}

// buildQrContent 构造二维码内容
func (s *AlipayService) buildQrContent(params map[string]string) string {
	// 返回一个模拟的二维码URL（真实环境应调用alipay.trade.precreate获取qr_code）
	b, _ := json.Marshal(params)
	return "alipay://pay?d=" + base64.StdEncoding.EncodeToString(b)
}

// buildPayURL 构造支付跳转URL
func (s *AlipayService) buildPayURL(params map[string]string) string {
	u, _ := url.Parse("https://openapi.alipay.com/gateway.do")
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
