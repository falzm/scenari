# scenari: Simple scenario execution

This package provides a simple library for writing arbitrary *scenarios*, which consist of steps executed sequentially. It can be useful for performing unit/functionnal testing.

## Example usage

This examples illustrates a *chaos engineering* scenario, where we describe different steps injecting various chaotic behavior into a web application using the [Chaos HTTP middleware](https://github.com/falzm/chaos):

```go
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/falzm/chaos"
	"github.com/falzm/scenari"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	chaoscli := chaos.NewClient("")

	if err := scenari.NewScenario("chaos").
		// Add a 1-second delay at a 0.5 probability to route POST /api/a
		Step(scenari.NewStep(func() error {
			return chaoscli.AddRouteChaos("POST", "/api/a", chaos.NewSpec().Delay(1000, 0.5))
		})).

		// Wait between 0 and 60 seconds
		Pause(time.Duration(rnd.Intn(60))*time.Second).

		// Add a 3-second delay at a 0.75 probability and return 504 status code error
		// at a 0.5 probability to route POST /api/a
		Step(scenari.NewStep(func() error {
			return chaoscli.AddRouteChaos("POST", "/api/a", chaos.NewSpec().
				Delay(3000, 0.75).
				Error(http.StatusGatewayTimeout, "Whoopsie!", 0.5))
		})).

		// Wait between 0 and 60 seconds
		Pause(time.Duration(rnd.Intn(60))*time.Second).

		// Repeat above scenario 2 more times, waiting 1 minute between each iteration
		Repeat(3, 1*time.Minute).
		Rollout(); err != nil {
		fmt.Printf("scenario rollout failed: %s\n", err)
	}

	// Once the scenario rollout is over, reset chaos effects to route POST /api/a
	chaoscli.DeleteRouteChaos("POST", "/api/a")
}

```
