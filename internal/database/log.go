package database

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm/logger"
)

func dbLogPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "jat", "db.log"), nil
}

// In TUI mode we default to Warn (migrations/queries don't spam the file).
// Set JAT_DEV=1 to bump it to Info so you can see them while developing.
func fileLogLevel() logger.LogLevel {
	if os.Getenv("JAT_DEV") != "" {
		return logger.Info
	}
	return logger.Warn
}

func newLogger(w io.Writer, level logger.LogLevel, colorful bool) logger.Interface {
	return logger.New(
		log.New(w, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  colorful,
		},
	)
}

// loggerFor returns a gorm logger appropriate to the run mode:
//   - server: stdout at Info (server has no TUI to corrupt);
//   - standalone/client: a db.log file next to the app config, at Warn
//     (Info if JAT_DEV is set), so nothing hits the terminal the TUI uses.
func loggerFor(server bool) (logger.Interface, error) {
	if server {
		return newLogger(os.Stdout, logger.Info, true), nil
	}
	path, err := dbLogPath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return newLogger(f, fileLogLevel(), false), nil
}
