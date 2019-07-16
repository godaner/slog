package slog

import (
	"log"
	"os"
)

// is default new function
// writer are configured by default
func DefaultNew(f func() SLogConfig) (*LoggerS, error) {

	cfg := f()
	logger := new(LoggerS)
	logger.cfg = &cfg
	logger.SetSliceType(cfg.SplitType)

	logger.SetDebug(cfg.Debug)

	writer := new(logWriter)

	if cfg.FileNameHandler == nil {
		cfg.FileNameHandler = cfg.name_handler
	}
	filename := cfg.FileNameHandler(0)

	file := &os.File{}
	file_info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
		} else {
			return nil, err
		}
	} else {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, err
		}
	}

	switch cfg.SplitType {
	case SPLIT_TYPE_FILE_SIZE:
		logger.SetMaxSize(cfg.Condition)
		if file_info != nil {
			logger.size = file_info.Size()
		}
	case SPLIT_TYPE_TIME_CYCLE:
		logger.SetIntervalsTime(cfg.Condition)
	}

	if err != nil {
		return nil, err
	}

	writer.file = file
	if cfg.Debug {
		writer.stdout = os.Stdout
	}
	logger.writer = writer
	logger.Logger = log.New(logger.writer, cfg.Prefix, cfg.LogFlag)

	return logger, nil
}
