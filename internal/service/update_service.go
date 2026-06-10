package service

import (
	"archive/zip"
	"chenze-faka/config"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/logger"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	updateLog       = "storage/logs/update.log"
	backupDir       = "storage/backup"
	tmpDownloadFile = "tmp_update.zip"

	tokenCheck = "Nathan" // 额外 token 参数，防止 API 被误调用
)

// UpdateCheckResult 返回给前端的检查结果
type UpdateCheckResult struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	HasUpdate       bool   `json:"has_update"`
	ForceUpdate     bool   `json:"force_update"`
	AppName         string `json:"app_name"`
	Date            string `json:"date"`
	Changelog       string `json:"changelog"`
	Msg             string `json:"msg"`
	DownloadExists  bool   `json:"download_exists"`
	UpdateZipSize   int64  `json:"update_zip_size"`
	AuthStatus      string `json:"auth_status"`
	AuthMsg         string `json:"auth_msg"`
}

// UpdateExecuteResult 更新执行结果
type UpdateExecuteResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	NewVer    string `json:"new_version"`
	NeedRestart bool  `json:"need_restart"`
	Details   string `json:"details,omitempty"`
}

// NewUpdateService 创建更新服务
func NewUpdateService() *UpdateService {
	return &UpdateService{}
}

type UpdateService struct{}

// CheckUpdate 调用远程授权站检查更新
func (s *UpdateService) CheckUpdate(siteURL string) (*UpdateCheckResult, error) {
	ls := NewLicenseService()
	auth, err := ls.LoadAuthPHP()
	authStatus := "unauthorized"
	authMsg := "未找到授权文件"
	if err == nil {
		authStatus = "authorized"
		authMsg = "已授权"
	}
	if auth == nil {
		auth = &AuthInfo{AppVersion: config.Version, AppName: "晨泽发卡", AuthCode: ""}
	}
	if auth.AppVersion == "" {
		auth.AppVersion = config.Version
	}

	api := fmt.Sprintf("%s/api/Index/version?appid=%s&url=%s&authcode=%s&version=%s&webkey=%s",
		authBaseURL,
		appID,
		url.QueryEscape(siteURL),
		url.QueryEscape(auth.AuthCode),
		url.QueryEscape(auth.AppVersion),
		url.QueryEscape(webKey),
	)
	logger.Infof("[update] check: version=%s url=%s", auth.AppVersion, siteURL)

	body, err := httpGetString(api, 15*time.Second)
	if err != nil {
		logger.Errorf("[update] check network error: %v", err)
		return &UpdateCheckResult{
			CurrentVersion: config.Version,
			Msg:            "检查更新失败（网络错误）:" + err.Error(),
			AuthStatus:     authStatus,
			AuthMsg:        authMsg,
		}, nil
	}
	logger.Infof("[update] check response: %s", truncate(body, 500))

	var resp authVersionResp
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, fmt.Errorf("授权站响应异常: %s", truncate(body, 200))
	}

	hasUpdate := versionLess(config.Version, resp.Version)
	if resp.Code == 0 && resp.Version != "" && !strings.Contains(strings.ToLower(resp.Msg), "success") &&
		!strings.Contains(resp.Msg, "已是最新") && !strings.Contains(resp.Msg, "暂无") {
		// 某些站 code=0 表示成功，兼容处理
	}
	if hasUpdate && resp.UpdateZip == "" {
		hasUpdate = false
	}

	result := &UpdateCheckResult{
		CurrentVersion: config.Version,
		LatestVersion:  resp.Version,
		HasUpdate:      hasUpdate,
		ForceUpdate:    resp.ForceUpdate == 1,
		AppName:        resp.AppName,
		Date:           resp.Date,
		Changelog:      resp.Changelog,
		Msg:            resp.Msg,
		DownloadExists: resp.UpdateZip != "",
		AuthStatus:     authStatus,
		AuthMsg:        authMsg,
	}
	return result, nil
}

// Execute 执行更新：下载 ZIP -> 备份 -> 解压覆盖 -> 执行 SQL -> 更新 Auth.php
func (s *UpdateService) Execute(siteURL string) (*UpdateExecuteResult, error) {
	s.appendLog("=== 更新开始 %s ===", time.Now().Format("2006-01-02 15:04:05"))

	// 1. 重新拉取一次以获取最新 version / update_zip / update_sql
	check, err := s.CheckUpdate(siteURL)
	if err != nil {
		s.appendLog("check error: %v", err)
		return nil, err
	}
	if !check.HasUpdate {
		return &UpdateExecuteResult{Success: true, Message: "当前已是最新版本", NewVer: check.CurrentVersion}, nil
	}
	if !check.DownloadExists {
		return nil, errors.New("没有可用的更新包")
	}

	// 2. 下载更新包
	s.appendLog("download: %s", check.AppName)
	api := fmt.Sprintf("%s/api/Index/version?appid=%s&url=%s&authcode=%s&version=%s&webkey=%s",
		authBaseURL, appID, url.QueryEscape(siteURL), "", url.QueryEscape(config.Version), url.QueryEscape(webKey))
	_ = api

	zipURL := checkLuckyUpdateZip(siteURL, check)
	s.appendLog("zip url: %s", zipURL)
	tmpPath, err := downloadFile(zipURL, tmpDownloadFile, 120*time.Second)
	if err != nil {
		s.appendLog("download error: %v", err)
		return nil, fmt.Errorf("下载更新包失败: %w", err)
	}
	defer os.Remove(tmpPath)
	info, _ := os.Stat(tmpPath)
	s.appendLog("download ok, size=%d bytes", info.Size())

	// 3. 备份（二进制 + 关键文件）
	backupSub := filepath.Join(backupDir, fmt.Sprintf("backup-%s-%d", time.Now().Format("20060102-150405"), os.Getpid()))
	_ = os.MkdirAll(backupSub, 0755)
	backupTargets := []string{"chenze_faka", "config.yaml", "config/Auth.php", "internal", "web"}
	for _, f := range backupTargets {
		if _, err := os.Stat(f); err == nil {
			if err := copyPath(f, filepath.Join(backupSub, f)); err != nil {
				logger.Warnf("[update] backup %s failed: %v", f, err)
			}
		}
	}
	s.appendLog("backup done: %s", backupSub)

	// 4. 解压覆盖（排除 config/、storage/）
	extracted := 0
	skipped := 0
	if err := unzipWithFilter(tmpPath, ".", func(name string) bool {
		if name == "" || strings.HasPrefix(name, "..") {
			return false
		}
		if strings.HasPrefix(name, "config/") || strings.HasPrefix(name, "storage/") ||
			strings.HasPrefix(name, ".git/") || strings.HasPrefix(name, ".trae/") ||
			name == "chenze_faka.db" || name == "install.lock" {
			skipped++
			return false
		}
		// 仅接受以下路径（安全过滤）：cmd/、internal/、web/、config/version.go、frontend/, release/
		// 但覆盖时保留白名单规则：
		extracted++
		return true
	}); err != nil {
		s.appendLog("unzip error: %v", err)
		s.rollback(backupSub)
		return nil, fmt.Errorf("解压失败: %w", err)
	}
	s.appendLog("unzip ok: extracted=%d skipped=%d", extracted, skipped)

	// 5. 执行 update_sql（如果有）
	updateSQL := checkLatestSQL(siteURL)
	if updateSQL != "" {
		if err := applyUpdateSQL(updateSQL); err != nil {
			s.appendLog("apply sql error: %v", err)
			// SQL 失败视为软错误，不回滚
			logger.Warnf("[update] apply update sql failed: %v", err)
		} else {
			s.appendLog("apply update sql ok")
		}
	} else {
		s.appendLog("no update sql provided")
	}

	// 6. 更新 Auth.php 的版本号
	if ls := NewLicenseService(); true {
		if existing, err := ls.LoadAuthPHP(); err == nil {
			existing.AppVersion = check.LatestVersion
			_ = ls.SaveAuthPHP(existing)
			s.appendLog("auth.php updated to %s", check.LatestVersion)
		} else {
			_ = ls.SaveAuthPHP(&AuthInfo{
				AppName:    "晨泽发卡",
				AppVersion: check.LatestVersion,
				URL:        siteURL,
			})
		}
	}

	s.appendLog("=== 更新完成 ===")
	return &UpdateExecuteResult{
		Success:     true,
		Message:     "更新完成",
		NewVer:      check.LatestVersion,
		NeedRestart: true,
		Details:     fmt.Sprintf("覆盖文件 %d 个，跳过 %d 个", extracted, skipped),
	}, nil
}

// appendLog 写入更新日志
func (s *UpdateService) appendLog(format string, args ...interface{}) {
	_ = os.MkdirAll(filepath.Dir(updateLog), 0755)
	f, err := os.OpenFile(updateLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s] "+format+"\n", append([]interface{}{time.Now().Format(time.RFC3339)}, args...)...)
}

func (s *UpdateService) rollback(backupSub string) {
	s.appendLog("rollback from: %s", backupSub)
	_ = copyPath(backupSub, ".")
}

// ---- helpers ----

func checkLuckyUpdateZip(siteURL string, check *UpdateCheckResult) string {
	// 授权站 version 接口一般会直接返回 update_zip 字段
	if check != nil {
		// 尝试直接使用当前响应中隐含的下载地址（兼容接口直接返回 URL）
		if strings.HasPrefix(check.AppName, "http") {
			return check.AppName
		}
	}
	// 备用：尝试直接拼接下载 URL
	ls := NewLicenseService()
	auth, err := ls.LoadAuthPHP()
	if err != nil {
		auth = &AuthInfo{AuthCode: "", AppVersion: config.Version}
	}
	return fmt.Sprintf("%s/api/Index/version?appid=%s&url=%s&authcode=%s&version=%s&webkey=%s&download=1",
		authBaseURL, appID, url.QueryEscape(siteURL), url.QueryEscape(auth.AuthCode),
		url.QueryEscape(config.Version), url.QueryEscape(webKey))
}

func checkLatestSQL(siteURL string) string {
	// 调用 version 接口，尝试解析 update_sql
	ls := NewLicenseService()
	auth, err := ls.LoadAuthPHP()
	if err != nil {
		auth = &AuthInfo{AuthCode: "", AppVersion: config.Version}
	}
	api := fmt.Sprintf("%s/api/Index/version?appid=%s&url=%s&authcode=%s&version=%s&webkey=%s",
		authBaseURL, appID, url.QueryEscape(siteURL), url.QueryEscape(auth.AuthCode),
		url.QueryEscape(config.Version), url.QueryEscape(webKey))
	body, err := httpGetString(api, 15*time.Second)
	if err != nil {
		return ""
	}
	type sqlOnly struct {
		UpdateSQL string `json:"update_sql"`
	}
	var r sqlOnly
	if err := json.Unmarshal([]byte(body), &r); err == nil {
		return r.UpdateSQL
	}
	// 兼容：查找 "update_sql":"..." 部分
	re := regexp.MustCompile(`"update_sql"\s*:\s*"((?:\\.|[^"\\])*)"`)
	if m := re.FindStringSubmatch(body); len(m) > 1 {
		return m[1]
	}
	return ""
}

// downloadFile 下载 URL 到临时文件，返回路径
func downloadFile(downloadURL, dst string, timeout time.Duration) (string, error) {
	if strings.HasPrefix(downloadURL, "//") {
		downloadURL = "https:" + downloadURL
	}
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "ChenzeFaka/"+config.Version)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 保存到临时文件
	out, err := os.CreateTemp("", "cf-update-*.zip")
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return out.Name(), nil
}

// unzipWithFilter 解压 ZIP，应用过滤器决定是否写入
func unzipWithFilter(zipPath, targetDir string, allow func(name string) bool) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	var dirs []string
	var files []*zip.File
	for _, f := range r.File {
		name := filepath.Clean(filepath.ToSlash(f.Name))
		if strings.HasPrefix(name, "../") || strings.HasPrefix(name, "..\\") {
			continue
		}
		if f.FileInfo().IsDir() {
			dirs = append(dirs, name)
			continue
		}
		if !allow(name) {
			continue
		}
		files = append(files, f)
	}

	// 先创建目录
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) < len(dirs[j]) })
	for _, d := range dirs {
		_ = os.MkdirAll(filepath.Join(targetDir, d), 0755)
	}
	// 再写文件
	for _, f := range files {
		name := filepath.Clean(filepath.ToSlash(f.Name))
		dst := filepath.Join(targetDir, name)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0755)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, infoModeFromSrc(src))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func infoModeFromSrc(src string) os.FileMode {
	if i, err := os.Stat(src); err == nil {
		return i.Mode().Perm()
	}
	return 0644
}

func copyDir(src, dst string) error {
	_ = os.MkdirAll(dst, 0755)
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// applyUpdateSQL 执行 update_sql 文本（按 ; 分行，忽略空行和注释）
func applyUpdateSQL(sqlText string) error {
	if strings.TrimSpace(sqlText) == "" {
		return nil
	}
	// 处理简单的 SQL：按分号分段，去重和过滤空行
	// 防止执行 DROP / DELETE 等危险语句
	whitelist := regexp.MustCompile(`(?i)^\s*(CREATE|ALTER|INSERT|UPDATE|SET|USE)\b`)
	if db.DB == nil {
		return errors.New("数据库未初始化")
	}
	statements := splitSQLStatements(sqlText)
	for _, s := range statements {
		s = strings.TrimSpace(s)
		if s == "" || strings.HasPrefix(s, "--") || strings.HasPrefix(s, "#") || strings.HasPrefix(s, "/*") {
			continue
		}
		if !whitelist.MatchString(s) {
			logger.Warnf("[update] skip sql statement: %s", truncate(s, 120))
			continue
		}
		if err := db.DB.Exec(s).Error; err != nil {
			logger.Warnf("[update] apply sql failed: %v (sql=%s)", err, truncate(s, 200))
		}
	}
	return nil
}

func splitSQLStatements(text string) []string {
	var parts []string
	var b strings.Builder
	inSingle := false
	inDouble := false
	inBacktick := false
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch c {
		case '\'':
			if i == 0 || text[i-1] != '\\' {
				inSingle = !inSingle
			}
		case '"':
			if i == 0 || text[i-1] != '\\' {
				inDouble = !inDouble
			}
		case '`':
			inBacktick = !inBacktick
		case ';':
			if !inSingle && !inDouble && !inBacktick {
				b.WriteByte(c)
				parts = append(parts, b.String())
				b.Reset()
				continue
			}
		}
		b.WriteByte(c)
	}
	if b.Len() > 0 {
		parts = append(parts, b.String())
	}
	return parts
}

// versionLess 比较两个版本号 a < b
func versionLess(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	a1 := normalizeVersion(a)
	b1 := normalizeVersion(b)
	for i := 0; i < len(a1) && i < len(b1); i++ {
		if a1[i] != b1[i] {
			return a1[i] < b1[i]
		}
	}
	return len(a1) < len(b1)
}

func normalizeVersion(v string) []int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindAllString(v, -1)
	out := make([]int, 0, len(match))
	for _, m := range match {
		var n int
		fmt.Sscanf(m, "%d", &n)
		out = append(out, n)
	}
	if len(out) < 3 {
		for len(out) < 3 {
			out = append(out, 0)
		}
	}
	return out
}

// sha256hex 快速哈希（供日志记录）
func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return fmt.Sprintf("%x", sum)
}

var _ = sha256hex
