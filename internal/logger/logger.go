package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init cấu hình logger toàn cục của zerolog.
// Sau khi gọi, dùng ở mọi nơi qua package "github.com/rs/zerolog/log":
//
//	log.Info().Str("user_id", id).Msg("get user")
//	log.Error().Err(err).Msg("something failed")
func Init() {
	zerolog.TimeFieldFormat = time.RFC3339

	// ConsoleWriter in log dạng dễ đọc khi dev.
	// Muốn log JSON (cho production) thì thay bằng: zerolog.New(os.Stdout)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
