package client

// NoOpLogger is a logger that does nothing
type NoOpLogger struct{}

// Errorf does nothing
func (NoOpLogger) Errorf(format string, v ...interface{}) {}

// Warnf does nothing
func (NoOpLogger) Warnf(format string, v ...interface{}) {}

// Debugf does nothing
func (NoOpLogger) Debugf(format string, v ...interface{}) {}

// Infof does nothing
func (NoOpLogger) Infof(format string, v ...interface{}) {}
