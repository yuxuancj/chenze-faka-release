package service

import (
	"chenze-faka/config"
	"chenze-faka/internal/pkg/logger"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	authBaseURL = "https://auth.seanld.com"
	webKey      = "CENathan_Auth"
	appID       = "2"
	authFile    = "config/Auth.php"

	statusOK    = "authorized"
	statusFail  = "unauthorized"
	statusError = "network_error"
)

// AuthInfo 授权信息（从 Auth.php 解析或安装时构建）
type AuthInfo struct {
	AppName     string `json:"app_name"`
	AppVersion  string `json:"app_version"`
	AuthCode    string `json:"authcode"`
	URL         string `json:"url"`
	InstalledAt int64  `json:"installed_at"`
	// 运行时字段
	Status      string `json:"status,omitempty"`
	StatusMsg   string `json:"status_msg,omitempty"`
	LastCheckAt int64  `json:"last_check_at,omitempty"`
}

// authCheckResp 授权站 check_auth 响应
type authCheckResp struct {
	Code int    `json:"code,string"`
	Msg  string `json:"msg"`
	// 可能额外返回的字段
	Expire    string `json:"expire,omitempty"`
	AuthType  string `json:"auth_type,omitempty"`
}

// authVersionResp 授权站 version 响应（检查更新）
type authVersionResp struct {
	Code        int    `json:"code,string"`
	Msg         string `json:"msg"`
	ID          int    `json:"id"`
	AppName     string `json:"appname"`
	Version     string `json:"version"`
	UpdateZip   string `json:"update_zip"`
	UpdateSQL   string `json:"update_sql"`
	ForceUpdate int    `json:"force_update"`
	Date        string `json:"date"`
	Changelog   string `json:"changelog,omitempty"`
}

var (
	licenseOnce sync.Once
	licenseCache *AuthInfo
	licenseMu    sync.RWMutex
)

// NewLicenseService 创建授权服务
func NewLicenseService() *LicenseService {
	return &LicenseService{}
}

type LicenseService struct{}

// DetectURL 从请求上下文探测当前站点域名
func (s *LicenseService) DetectURL(req *http.Request) string {
	if req == nil {
		return ""
	}
	host := req.Host
	if host == "" {
		host = "localhost"
	}
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	if fw := req.Header.Get("X-Forwarded-Proto"); fw != "" {
		scheme = fw
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

// VerifyLicense 验证授权码（安装时调用，会请求授权站）
func (s *LicenseService) VerifyLicense(siteURL, authCode string) (string, error) {
	if strings.TrimSpace(authCode) == "" {
		return statusFail, errors.New("授权码不能为空")
	}
	if len(authCode) < 32 || len(authCode) > 128 {
		return statusFail, errors.New("授权码格式不正确")
	}

	api := fmt.Sprintf("%s/api/Index/check_auth?appid=%s&url=%s&authcode=%s",
		authBaseURL,
		appID,
		url.QueryEscape(siteURL),
		url.QueryEscape(authCode),
	)
	logger.Infof("[license] verify: url=%s", siteURL)
	body, err := httpGetString(api, 8*time.Second)
	if err != nil {
		logger.Errorf("[license] verify network error: %v", err)
		return statusError, fmt.Errorf("授权服务器不可达: %w", err)
	}
	logger.Infof("[license] verify response: %s", body)

	var resp authCheckResp
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		// 尝试兼容字符串形式的 code
		return statusError, fmt.Errorf("授权服务器响应异常: %s", truncate(body, 200))
	}
	if resp.Code == 1 || strings.Contains(resp.Msg, "正版授权") || strings.Contains(resp.Msg, "授权有效") {
		return statusOK, nil
	}
	if resp.Msg != "" {
		return statusFail, errors.New(resp.Msg)
	}
	return statusFail, errors.New("授权验证失败")
}

// SaveAuthPHP 写入 config/Auth.php
func (s *LicenseService) SaveAuthPHP(info *AuthInfo) error {
	dir := filepath.Dir(authFile)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	content := fmt.Sprintf(`<?php
// 晨泽发卡系统授权配置
// 本文件由安装向导自动生成，请勿手动修改。
return [
    'app_version'  => '%s',
    'app_name'     => '%s',
    'authcode'     => '%s',
    'url'          => '%s',
    'installed_at' => %d,
];
`,
		escapePHP(info.AppVersion),
		escapePHP(info.AppName),
		escapePHP(info.AuthCode),
		escapePHP(info.URL),
		info.InstalledAt,
	)
	return os.WriteFile(authFile, []byte(content), 0644)
}

// LoadAuthPHP 从 Auth.php 读取授权信息
func (s *LicenseService) LoadAuthPHP() (*AuthInfo, error) {
	data, err := os.ReadFile(authFile)
	if err != nil {
		return nil, err
	}
	info := parsePHPReturn(string(data))
	if info.AuthCode == "" {
		return nil, fmt.Errorf("authcode 未找到")
	}
	if info.AppVersion == "" {
		info.AppVersion = config.Version
	}
	if info.AppName == "" {
		info.AppName = "晨泽发卡"
	}
	return &info, nil
}

// GetInfo 读取当前授权信息（带缓存）
func (s *LicenseService) GetInfo() *AuthInfo {
	licenseMu.RLock()
	if licenseCache != nil {
		defer licenseMu.RUnlock()
		cp := *licenseCache
		return &cp
	}
	licenseMu.RUnlock()

	info, err := s.LoadAuthPHP()
	if err != nil {
		info = &AuthInfo{AppVersion: config.Version, AppName: "晨泽发卡", Status: statusFail, StatusMsg: "未找到授权文件"}
	} else {
		info.Status = statusOK
		info.StatusMsg = "已授权"
	}

	licenseMu.Lock()
	licenseCache = info
	licenseMu.Unlock()

	cp := *info
	return &cp
}

// RuntimeCheck 运行时检查（后台可调用），结果写入缓存
func (s *LicenseService) RuntimeCheck(siteURL string) *AuthInfo {
	info, err := s.LoadAuthPHP()
	if err != nil {
		return &AuthInfo{Status: statusFail, StatusMsg: "授权文件不存在"}
	}
	status, err := s.VerifyLicense(siteURL, info.AuthCode)
	info.Status = status
	info.LastCheckAt = time.Now().Unix()
	if err != nil {
		info.StatusMsg = err.Error()
	} else if status == statusOK {
		info.StatusMsg = "正版授权"
	} else {
		info.StatusMsg = "授权验证失败"
	}
	licenseMu.Lock()
	licenseCache = info
	licenseMu.Unlock()
	cp := *info
	return &cp
}

// ---- helpers ----

func httpGetString(u string, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "ChenzeFaka/"+config.Version)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf := make([]byte, 0, 1024)
	tmp := make([]byte, 4096)
	for {
		n, rerr := resp.Body.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if rerr != nil {
			break
		}
		if len(buf) > 1024*64 { // 限制 64KB
			break
		}
	}
	return string(buf), nil
}

// parsePHPReturn 以最小方式解析 PHP return 数组中的值
// 只提取 app_version / app_name / authcode / url / installed_at
func parsePHPReturn(content string) AuthInfo {
	info := AuthInfo{}
	fields := map[string]*string{
		"app_version":  &info.AppVersion,
		"app_name":     &info.AppName,
		"authcode":     &info.AuthCode,
		"url":          &info.URL,
	}
	intFields := map[string]*int64{
		"installed_at": &info.InstalledAt,
	}
	re := regexp.MustCompile(`'([a-zA-Z_]+)'\s*=>\s*'([^']*)'`)
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		if ptr, ok := fields[m[1]]; ok {
			*ptr = m[2]
		}
	}
	reInt := regexp.MustCompile(`'([a-zA-Z_]+)'\s*=>\s*(\d+)`)
	for _, m := range reInt.FindAllStringSubmatch(content, -1) {
		if ptr, ok := intFields[m[1]]; ok {
			var v int64
			fmt.Sscanf(m[2], "%d", &v)
			*ptr = v
		}
	}
	return info
}

func escapePHP(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
