package test

import (
	"errors"
	"testing"
	"time"

	"gin-server/configmanager/common/alert"
)

func TestLogAlerter(t *testing.T) {
	// 创建日志告警器
	alerter := alert.NewLogAlerter()

	// 测试不同级别的告警
	testCases := []struct {
		name    string
		alert   *alert.Alert
		wantErr bool
	}{
		{
			name: "信息级别告警",
			alert: &alert.Alert{
				Level:      alert.AlertLevelInfo,
				Type:       alert.AlertTypeLogGenerate,
				Message:    "测试信息告警",
				Error:      nil,
				RetryCount: 0,
				Timestamp:  time.Now(),
				Module:     "TestModule",
			},
			wantErr: false,
		},
		{
			name: "警告级别告警",
			alert: &alert.Alert{
				Level:      alert.AlertLevelWarning,
				Type:       alert.AlertTypeLogEncrypt,
				Message:    "测试警告告警",
				Error:      errors.New("警告错误"),
				RetryCount: 1,
				Timestamp:  time.Now(),
				Module:     "TestModule",
			},
			wantErr: false,
		},
		{
			name: "错误级别告警",
			alert: &alert.Alert{
				Level:      alert.AlertLevelError,
				Type:       alert.AlertTypeLogUpload,
				Message:    "测试错误告警",
				Error:      errors.New("严重错误"),
				RetryCount: 2,
				Timestamp:  time.Now(),
				Module:     "TestModule",
			},
			wantErr: false,
		},
		{
			name: "致命级别告警",
			alert: &alert.Alert{
				Level:      alert.AlertLevelFatal,
				Type:       alert.AlertTypeStrategySync,
				Message:    "测试致命告警",
				Error:      errors.New("致命错误"),
				RetryCount: 3,
				Timestamp:  time.Now(),
				Module:     "TestModule",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := alerter.Alert(tc.alert)
			if (err != nil) != tc.wantErr {
				t.Errorf("Alert() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestDefaultAlerter(t *testing.T) {
	// 获取默认告警器
	defaultAlerter := alert.GetDefaultAlerter()
	if defaultAlerter == nil {
		t.Error("GetDefaultAlerter() 返回了空的告警器")
	}

	// 测试设置新的默认告警器
	newAlerter := alert.NewLogAlerter()
	alert.SetDefaultAlerter(newAlerter)

	// 验证新的默认告警器是否生效
	currentDefaultAlerter := alert.GetDefaultAlerter()
	if currentDefaultAlerter != newAlerter {
		t.Error("SetDefaultAlerter() 未能正确设置新的默认告警器")
	}

	// 测试使用新的默认告警器发送告警
	testAlert := &alert.Alert{
		Level:      alert.AlertLevelInfo,
		Type:       alert.AlertTypeLogGenerate,
		Message:    "测试默认告警器",
		Error:      nil,
		RetryCount: 0,
		Timestamp:  time.Now(),
		Module:     "TestModule",
	}

	err := currentDefaultAlerter.Alert(testAlert)
	if err != nil {
		t.Errorf("使用默认告警器发送告警失败: %v", err)
	}
}

func TestAlertLevelString(t *testing.T) {
	tests := []struct {
		level    alert.AlertLevel
		expected string
	}{
		{alert.AlertLevelInfo, "INFO"},
		{alert.AlertLevelWarning, "WARNING"},
		{alert.AlertLevelError, "ERROR"},
		{alert.AlertLevelFatal, "FATAL"},
		{alert.AlertLevel(99), "UNKNOWN"},
	}

	for _, test := range tests {
		if got := test.level.String(); got != test.expected {
			t.Errorf("AlertLevel(%d).String() = %v, 期望 %v", test.level, got, test.expected)
		}
	}
}

func TestAlertTypeString(t *testing.T) {
	tests := []struct {
		alertType alert.AlertType
		expected  string
	}{
		{alert.AlertTypeLogGenerate, "LOG_GENERATE"},
		{alert.AlertTypeLogEncrypt, "LOG_ENCRYPT"},
		{alert.AlertTypeLogUpload, "LOG_UPLOAD"},
		{alert.AlertTypeStrategySync, "STRATEGY_SYNC"},
		{alert.AlertTypeStrategyApply, "STRATEGY_APPLY"},
		{alert.AlertType(99), "UNKNOWN"},
	}

	for _, test := range tests {
		if got := test.alertType.String(); got != test.expected {
			t.Errorf("AlertType(%d).String() = %v, 期望 %v", test.alertType, got, test.expected)
		}
	}
}
