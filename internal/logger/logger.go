package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init cấu hình logger toàn cục của zerolog, ghi ra 2 nơi:
//   - Console: dạng dễ đọc khi dev
//   - File logPath: dạng JSON, mỗi dòng 1 log (dễ grep / đưa vào hệ thống log)
//
// Trả về file đã mở để main đóng lại khi tắt server.
// Sau khi gọi, dùng ở mọi nơi qua package "github.com/rs/zerolog/log":
//
//	log.Info().Str("user_id", id).Msg("get user")
//	log.Error().Err(err).Msg("something failed")
func Init(logPath string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = time.RFC3339

	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	output := zerolog.MultiLevelWriter(console, file)

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
	return file, nil
}
