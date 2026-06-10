package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/logger"
	"fmt"

	"gorm.io/gorm"
)

type SettingService struct{}

func NewSettingService() *SettingService { return &SettingService{} }

func (s *SettingService) Get(key string) (string, error) {
	var setting model.Setting
	err := db.DB.Where("`key` = ?", key).First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		logger.Errorf("SettingService.Get error: %v", err)
		return "", err
	}
	return setting.Value, nil
}

func (s *SettingService) Set(key string, value string) error {
	var setting model.Setting
	err := db.DB.Where("`key` = ?", key).First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		setting.Key = key
		setting.Value = value
		err = db.DB.Create(&setting).Error
	} else if err == nil {
		setting.Value = value
		err = db.DB.Save(&setting).Error
	}
	if err != nil {
		logger.Errorf("SettingService.Set error: %v", err)
	}
	return err
}

func (s *SettingService) GetFloat(key string, defaultValue float64) float64 {
	val, err := s.Get(key)
	if err != nil || val == "" {
		return defaultValue
	}
	var result float64
	_, err = fmt.Sscanf(val, "%f", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

func (s *SettingService) GetInt(key string, defaultValue int) int {
	val, err := s.Get(key)
	if err != nil || val == "" {
		return defaultValue
	}
	var result int
	_, err = fmt.Sscanf(val, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

func (s *SettingService) GetBool(key string, defaultValue bool) bool {
	val, err := s.Get(key)
	if err != nil || val == "" {
		return defaultValue
	}
	return val == "true" || val == "1"
}

func (s *SettingService) InitDefaults() {
	defaults := map[string]string{
		"site_name":              "晨泽发卡系统",
		"site_desc":              "专业的数字商品发卡平台",
		"epay_enabled":           "false",
		"epay_url":               "",
		"epay_pid":               "",
		"epay_key":               "",
		"alipay_enabled":         "false",
		"alipay_app_id":          "",
		"alipay_private_key":     "",
		"alipay_public_key":      "",
		"alipay_notify_url":      "",
		"alipay_return_url":      "",
		"distrib_level1_rate":    "10",
		"distrib_level2_rate":    "5",
		"distrib_level3_rate":    "2",
		"distrib_min_withdraw":   "10",
		"distrib_enabled":        "true",
		"points_rate":            "1",
		"points_discount_rate":   "1",
		"points_signin_reward":   "10",
		"points_continuous_reward": "5",
		"points_max_discount_percent": "50",
		"points_enabled":          "true",
		"order_timeout_minutes":   "30",
	}

	for key, value := range defaults {
		existing, _ := s.Get(key)
		if existing == "" {
			s.Set(key, value)
			logger.Infof("Setting default value for %s", key)
		}
	}
}

func (s *SettingService) GetAll() ([]model.Setting, error) {
	var settings []model.Setting
	err := db.DB.Find(&settings).Error
	if err != nil {
		logger.Errorf("SettingService.GetAll error: %v", err)
	}
	return settings, err
}
