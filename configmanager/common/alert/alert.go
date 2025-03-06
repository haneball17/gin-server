package alert

import (
	"fmt"
	"log"
	"time"

	"gin-server/config"
)

// AlertLevel 告警级别
type AlertLevel int

const (
	AlertLevelInfo    AlertLevel = 1 // 信息
	AlertLevelWarning AlertLevel = 2 // 警告
	AlertLevelError   AlertLevel = 3 // 错误
	AlertLevelFatal   AlertLevel = 4 // 致命
)

// AlertType 告警类型
type AlertType int

const (
	AlertTypeLogGenerate   AlertType = 1 // 日志生成
	AlertTypeLogEncrypt    AlertType = 2 // 日志加密
	AlertTypeLogUpload     AlertType = 3 // 日志上传
	AlertTypeStrategySync  AlertType = 4 // 策略同步
	AlertTypeStrategyApply AlertType = 5 // 策略应用
)

// Alert 告警信息
type Alert struct {
	Level      AlertLevel `json:"level"`      // 告警级别
	Type       AlertType  `json:"type"`       // 告警类型
	Message    string     `json:"message"`    // 告警消息
	Error      error      `json:"error"`      // 错误信息
	RetryCount int        `json:"retryCount"` // 重试次数
	Timestamp  time.Time  `json:"timestamp"`  // 告警时间
	Module     string     `json:"module"`     // 告警模块
}

// Alerter 告警接口
type Alerter interface {
	Alert(alert *Alert) error
}

// LogAlerter 日志告警器
type LogAlerter struct {
	logger *log.Logger
}

// NewLogAlerter 创建日志告警器
func NewLogAlerter() *LogAlerter {
	return &LogAlerter{
		logger: log.Default(),
	}
}

// Alert 实现告警接口
func (a *LogAlerter) Alert(alert *Alert) error {
	cfg := config.GetConfig()
	if cfg.DebugLevel == "true" {
		log.Printf("触发告警: %+v\n", alert)
	}

	// 格式化告警消息
	alertMsg := fmt.Sprintf("[%s][%s] %s - %v (重试次数: %d)",
		alert.Level.String(),
		alert.Type.String(),
		alert.Message,
		alert.Error,
		alert.RetryCount)

	// 记录告警日志
	a.logger.Println(alertMsg)
	return nil
}

// String 告警级别字符串表示
func (l AlertLevel) String() string {
	switch l {
	case AlertLevelInfo:
		return "INFO"
	case AlertLevelWarning:
		return "WARNING"
	case AlertLevelError:
		return "ERROR"
	case AlertLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// String 告警类型字符串表示
func (t AlertType) String() string {
	switch t {
	case AlertTypeLogGenerate:
		return "LOG_GENERATE"
	case AlertTypeLogEncrypt:
		return "LOG_ENCRYPT"
	case AlertTypeLogUpload:
		return "LOG_UPLOAD"
	case AlertTypeStrategySync:
		return "STRATEGY_SYNC"
	case AlertTypeStrategyApply:
		return "STRATEGY_APPLY"
	default:
		return "UNKNOWN"
	}
}
