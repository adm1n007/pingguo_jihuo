package inputRecord

import (
    "time"
    "ml/random"
)

func randomTime(min, max int) time.Duration {
    return time.Duration(random.IntRange(min, max)) * time.Millisecond
}
