package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// InitLogger инициализирует Zerolog с выводом в файл
func InitLogger() zerolog.Logger {
	// Открываем файл для записи логов (создается, если нет)
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open log file")
	}

	// Настраиваем zerolog на вывод в файл с human-friendly форматом
	logger := zerolog.New(file).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	return logger
}
