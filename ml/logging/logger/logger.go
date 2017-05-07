package logger

import (
    "ml/logging"
    "ml/os2"
)

var (
    logger      = logging.NewLogger(os2.ExecutableName())
    Debug       = logger.Debug
    Info        = logger.Info
    Warning     = logger.Warning
    Error       = logger.Error
    Fatal       = logger.Fatal
    SetLevel    = logger.SetLevel
    Level       = logger.Level
    SetFormater = logger.SetFormater
    LogToFile   = logger.LogToFile
)

func init() {
    logger.SetSkip(1)
}
