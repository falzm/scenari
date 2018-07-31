package scenari

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func testPrintFunc(buf *bytes.Buffer, s string) func() error {
	return func() error {
		buf.WriteString(s)
		return nil
	}
}

func Test_Scenario(t *testing.T) {
	var (
		buf bytes.Buffer

		expectedResult      = bytes.NewBufferString("[...][...][...]")
		expectedMinDuration = 10 * time.Second
	)

	startTime := time.Now()

	err := NewScenario("").
		Step(NewStep(testPrintFunc(&buf, "...")).
			PreExec(testPrintFunc(&buf, "[")).
			PostExec(testPrintFunc(&buf, "]"))).
		Step(NewStep(func() error { return fmt.Errorf("oh noes :(") })).
		Repeat(3, 5*time.Second).
		CarryOn().
		Rollout()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	duration := time.Now().Sub(startTime)
	if duration < expectedMinDuration {
		t.Errorf("%s: expected scenario duration to be at least %s, took %s", t.Name(), expectedMinDuration, duration)
		t.FailNow()
	}

	if !bytes.Equal(buf.Bytes(), expectedResult.Bytes()) {
		t.Errorf("%s: expected %s, got %s", t.Name(), expectedResult.String(), buf.String())
		t.FailNow()
	}
}
