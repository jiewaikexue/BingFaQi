package Time

import (
	"time"
)

/*
格林尼治标准时间（GMT）：表示为 "GMT" 或 "UTC"。
世界协调时（UTC）：表示为 "UTC"。
美国东部时间（EST）：表示为 "America/New_York"。
美国中部时间（CST）：表示为 "America/Chicago"。
美国山区时间（MST）：表示为 "America/Denver"。
美国太平洋时间（PST）：表示为 "America/Los_Angeles"。
东京时间（JST）：表示为 "Asia/Tokyo"。
伦敦时间（GMT+0）：表示为 "Europe/London"。
中欧时间（CET）：表示为 "CET"。
北京时间（CST）：表示为 "Asia/Shanghai"。
悉尼时间（AEST）：表示为 "Australia/Sydney"。
新加坡时间（SGT）：表示为 "Asia/Singapore"。
孟买时间（IST）：表示为 "Asia/Kolkata"。
*/

// NewScheduledTimerFromString 接受一个时间字符串和一个时区字符串作为参数，
// 返回一个在指定时区的指定时间点触发的定时器。
// 入参统一 不要日期少些0
func NewScheduledTimerFromString(timeStr string, timeZoneStr string) (*time.Timer, error) {
	// 解析时区字符串
	location, err := time.LoadLocation(timeZoneStr)
	if err != nil {
		return nil, err
	}

	// 解析时间字符串为时间对象
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, location)
	if err != nil {
		return nil, err
	}
	// 计算倒计时
	duration := t.Sub(time.Now())
	// 创建定时器
	timer := time.NewTimer(duration)
	return timer, nil
}
