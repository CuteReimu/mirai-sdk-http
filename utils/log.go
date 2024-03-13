package utils

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"log/slog"
	"path"
	"time"
)

// InitLogger 将日志转储至文件。
//
// 用法: InitLogger("./logs", LogInfoLevel, LogWithStack)
func InitLogger(p string, logLevel slog.Level) {
	writer, err := rotatelogs.New(
		path.Join(p, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		slog.Error("unable to write logs", "error", err)
		return
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("15:04:05.000"))
				}
			default:
				if e, ok := a.Value.Any().(error); ok {
					a.Value = slog.StringValue(fmt.Sprintf("%+v", e))
				}
			}
			return a
		}})))
}
