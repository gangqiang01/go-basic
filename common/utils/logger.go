package utils

import (
	"k8s.io/klog/v2"
)

/*
* logger wrapper for klog.
 */
type LoggerWrapper struct {
	Level int
}

func (l *LoggerWrapper) Info(args ...interface{}) {
	klog.Info(args...)
}

func (l *LoggerWrapper) Infoln(args ...interface{}) {
	klog.Infoln(args...)
}

func (l *LoggerWrapper) Infof(format string, args ...interface{}) {
	klog.Infof(format, args...)
}

func (l *LoggerWrapper) Warning(args ...interface{}) {
	klog.Warning(args...)
}

func (l *LoggerWrapper) Warningln(args ...interface{}) {
	klog.Warningln(args...)
}

func (l *LoggerWrapper) Warningf(format string, args ...interface{}) {
	klog.Warningf(format, args...)
}

func (l *LoggerWrapper) Error(args ...interface{}) {
	klog.Error(args...)
}

func (l *LoggerWrapper) Errorln(args ...interface{}) {
	klog.Errorln(args...)
}

func (l *LoggerWrapper) Errorf(format string, args ...interface{}) {
	klog.Errorf(format, args...)
}

func (l *LoggerWrapper) Fatal(args ...interface{}) {
	klog.Fatal(args...)
}

func (l *LoggerWrapper) Fatalln(args ...interface{}) {
	klog.Fatalln(args...)
}

func (l *LoggerWrapper) Fatalf(format string, args ...interface{}) {
	klog.Fatalf(format, args...)
}

func (l *LoggerWrapper) V(level int) bool {
	return false
}
