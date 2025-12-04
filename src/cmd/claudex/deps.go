package main

import (
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

// Commander abstracts process execution for testability
type Commander interface {
	// Run executes command and returns combined output
	Run(name string, args ...string) ([]byte, error)
	// Start launches interactive command with stdio attached
	Start(name string, stdin io.Reader, stdout, stderr io.Writer, args ...string) error
}

// Environment abstracts environment variable access
type Environment interface {
	Get(key string) string
	Set(key, value string)
}

// Clock abstracts time for testability
type Clock interface {
	Now() time.Time
}

// UUIDGenerator abstracts UUID generation for testability
type UUIDGenerator interface {
	New() string
}

// OsCommander is the production implementation of Commander
type OsCommander struct{}

func (c *OsCommander) Run(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

func (c *OsCommander) Start(name string, stdin io.Reader, stdout, stderr io.Writer, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

// OsEnv is the production implementation of Environment
type OsEnv struct{}

func (e *OsEnv) Get(key string) string {
	return os.Getenv(key)
}

func (e *OsEnv) Set(key, value string) {
	os.Setenv(key, value)
}

// SystemClock is the production implementation of Clock
type SystemClock struct{}

func (c *SystemClock) Now() time.Time {
	return time.Now()
}

// SystemUUID is the production implementation of UUIDGenerator
type SystemUUID struct{}

func (u *SystemUUID) New() string {
	return uuid.New().String()
}

// Package-level default instances for production use
var AppCmd Commander = &OsCommander{}
var AppClock Clock = &SystemClock{}
var AppUUID UUIDGenerator = &SystemUUID{}
var AppEnv Environment = &OsEnv{}
