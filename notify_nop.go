//go:build !systemd

package aged

func NotifyReady() error {
	return nil
}

func NotifyShutdown() error {
	return nil
}
