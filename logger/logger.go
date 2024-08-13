package logger

import (
	"fmt"
)

func Log(pref string, message ...interface{}) {
	str := fmt.Sprintf("[%s] ", pref)
	for _, m := range message {
		str += fmt.Sprintf("%v ", m)
	}

	fmt.Println(str)
}

func Logf(pref string, format string, message ...interface{}) {
	str := fmt.Sprintf("[%s] %s", pref, format)
	fmt.Printf(str, message...)
}

func Tent(message ...interface{})    { Log("🎪", message...) }
func Info(message ...interface{})    { Log("📝", message...) }
func Warn(message ...interface{})    { Log("☢", message...) }
func Error(message ...interface{})   { Log("❌", message...) }
func Success(message ...interface{}) { Log("✅", message...) }

func Tentf(format string, message ...interface{})    { Logf("🎪", format, message...) }
func Infof(format string, message ...interface{})    { Logf("📝", format, message...) }
func Warnf(format string, message ...interface{})    { Logf("☢", format, message...) }
func Errorf(format string, message ...interface{})   { Logf("❌", format, message...) }
func Successf(format string, message ...interface{}) { Logf("✅", format, message...) }
