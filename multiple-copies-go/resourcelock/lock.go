package resourcelock

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type FileLock struct {
	path     string
	holderID string
	file     *os.File
}

func New(path, holderID string) *FileLock {
	return &FileLock{
		path:     path,
		holderID: holderID,
	}
}

// TryLock 尝试以非阻塞方式获取文件排他锁。如果锁已被获得，则返回 true，如果已经被其他进程持有，则返回 false。
func (fl *FileLock) TryLock() (bool, error) {
	f, err := os.OpenFile(fl.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false, fmt.Errorf("open lock file: %w", err)
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == syscall.EWOULDBLOCK {
		f.Close()
		return false, nil
	}
	if err != nil {
		f.Close()
		return false, fmt.Errorf("flock: %w", err)
	}

	// Truncate and write holder info
	if err := f.Truncate(0); err != nil {
		f.Close()
		return false, fmt.Errorf("truncate: %w", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		f.Close()
		return false, fmt.Errorf("seek: %w", err)
	}
	if _, err := f.WriteString(fl.holderID); err != nil {
		f.Close()
		return false, fmt.Errorf("write holder: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return false, fmt.Errorf("sync: %w", err)
	}

	fl.file = f
	return true, nil
}

// Heartbeat rewrites the holder info to the lock file, proving liveness.
func (fl *FileLock) Heartbeat() error {
	if fl.file == nil {
		return fmt.Errorf("not holding lock")
	}
	if err := fl.file.Truncate(0); err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	if _, err := fl.file.Seek(0, 0); err != nil {
		return fmt.Errorf("seek: %w", err)
	}
	line := fmt.Sprintf("%s\n%s", fl.holderID, time.Now().Format(time.RFC3339))
	if _, err := fl.file.WriteString(line); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return fl.file.Sync()
}

// Unlock releases the lock.
func (fl *FileLock) Unlock() error {
	if fl.file == nil {
		return nil
	}

	if err := fl.file.Truncate(0); err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	if _, err := fl.file.Seek(0, 0); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	if err := syscall.Flock(int(fl.file.Fd()), syscall.LOCK_UN); err != nil {
		return fmt.Errorf("flock unlock: %w", err)
	}

	if err := fl.file.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	fl.file = nil
	return nil
}

type LeaseInfo struct {
	HolderID  string
	UpdatedAt time.Time
}

// GetLease parses the lock file and returns the holder ID and last heartbeat time.
// Returns nil if the file is empty or unparseable.
func GetLease(path string) *LeaseInfo {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if content == "" {
		return nil
	}

	// Format: "holderID\ntimestamp"
	lines := strings.SplitN(strings.TrimSpace(content), "\n", 2)
	if len(lines) < 2 {
		// Fallback: just holder ID, no timestamp — set current time to avoid false expiry
		return &LeaseInfo{HolderID: strings.TrimSpace(lines[0]), UpdatedAt: time.Now()}
	}

	ts, err := time.Parse(time.RFC3339, strings.TrimSpace(lines[1]))
	if err != nil {
		return &LeaseInfo{HolderID: strings.TrimSpace(lines[0]), UpdatedAt: time.Now()}
	}

	return &LeaseInfo{
		HolderID:  strings.TrimSpace(lines[0]),
		UpdatedAt: ts,
	}
}

// GetHolder returns the current lock holder by reading the lock file.
func GetHolder(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
