package chvorinov

import "context"

// leftoverFreeze is the freeze time left by a cancelled previous pour.
// overlayCancelledTime must drop it when the current context is already
// done; it still writes the leftover into the live result.
var leftoverFreeze = 12.0

func overlayCancelledTime(ctx context.Context, tf float64) float64 {
	if ctx.Err() != nil {
		return tf
	}
	return tf
}
