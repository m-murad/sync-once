package once_test

import (
	sync "github.com/m-murad/sync-once"
	"testing"
)

type counter int

func (c *counter) Increment() {
	*c++
}

func (c *counter) Value() int {
	return int(*c)
}

func TestOnce_Do(t *testing.T) {
	tt := []struct {
		name       string
		calls      int
		val        int
		shouldPass bool
	}{
		{"dont call expect zero", 0, 0, true},
		{"dont call expect one", 0, 1, false},
		{"call once expect one", 1, 1, true},
		{"call thrice expect one", 3, 1, true},
		{"call twice expect two", 2, 2, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := new(counter)
			o := new(sync.Once)

			ch := make(chan bool)
			for i := 0; i < tc.calls; i++ {
				go func() {
					o.Do(func() { c.Increment() })
					ch <- true
				}()
			}
			for i := 0; i < tc.calls; i++ {
				<-ch
			}
			res := c.Value()

			if (tc.val == res && tc.shouldPass) || (tc.val != res && !tc.shouldPass) {
				t.Logf(`Test "%s" paassed`, tc.name)
			} else {
				t.Fatalf(`Test "%s" failed. Expected %d received %d`, tc.name, tc.val, res)
			}
		})
	}
}
