//go:build systemd

package aged

import (
	"log"

	"github.com/mdlayher/sdnotify"
)

var notify *sdnotify.Notifier

func init() {
	var err error
	notify, err = sdnotify.New()
	if err != nil {
		log.Printf("Failed to open sd_notify socket: %v", err.Error())
	}
}

func NotifyReady() error {
	if notify == nil {
		return nil
	}
	return notify.Notify(sdnotify.Ready)
}

func NotifyShutdown() error {
	if notify == nil {
		return nil
	}
	return notify.Notify(sdnotify.Stopping)
}
