package periodic

import "time"

func RunPeriodically(f func(), period time.Duration) {
	go func() {
		for {
			f()
			time.Sleep(period)
		}
	}()
}
