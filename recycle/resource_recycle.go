package recycle

import (
	"context"

	"prometheus-test/lib/logger"
)

var (
	recycles []ResourceRecyclable
)

type ResourceRecyclable func() bool

func RegisterRecycles(recycle ResourceRecyclable) {
	recycles = append(recycles, recycle)
}

func ReleaseResources() {
	for _, r := range recycles {
		r()
	}
	logger.Warn(context.Background(), "release all closable resources ...")
}
