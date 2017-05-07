package entry

import (
	. "ml/strings"
	"time"

	"activator"
)

type entryCallbacks struct {
	findVerifyUrls func(*activator.AppleIdActivator) (activator.ActivateResult, []String)
}

var callbacks = entryCallbacks{}


var (
	successTotal int32
	failureTotal int32
	existsTotal  int32
	networkError int32
	startTime    time.Time
)
