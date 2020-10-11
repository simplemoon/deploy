package utils

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// 文件锁
type FileLock struct {
	p string   // 路径
	f *os.File // 对应的文件
}

// 创建一个文件锁
func NewFileLock(p string) *FileLock {
	return &FileLock{
		p: p,
	}
}

// 加锁
func (l *FileLock) Lock() error {
	f, err := os.Open(l.p)
	if err != nil {
		return err
	}
	l.f = f
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return err
	}
	return nil
}

// 加锁但有时间限制
func (l *FileLock) LockWithTime(t time.Duration) error {
	var cnt time.Duration = 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for ; cnt*time.Second < t; cnt++ {
		select {
		case <-ticker.C:
			if err := l.Lock(); err == nil {
				return err
			}
		}
	}
	return fmt.Errorf("lock %s timeout", l.p)
}

// 释放锁
func (l *FileLock) UnLock() error {
	defer func() {
		if err := l.f.Close(); err != nil {
			return
		}
		// 删除文件
		os.Remove(l.p)
	}()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
