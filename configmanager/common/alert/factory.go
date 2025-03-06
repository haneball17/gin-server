package alert

import "sync"

var (
	defaultAlerter Alerter
	once           sync.Once
)

// GetDefaultAlerter 获取默认告警器
func GetDefaultAlerter() Alerter {
	once.Do(func() {
		defaultAlerter = NewLogAlerter()
	})
	return defaultAlerter
}

// SetDefaultAlerter 设置默认告警器
func SetDefaultAlerter(alerter Alerter) {
	defaultAlerter = alerter
}
