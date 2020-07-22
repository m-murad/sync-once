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
				<- ch
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

func TestOnce_DoAgain(t *testing.T) {
	const (
		Do = iota
		DoAgain
	)
	tt := []struct {
		name       string
		funCalls   []int
		val        int
		shouldPass bool
	}{
		{"Func calls - [DoAgain]. Expect one", []int{DoAgain}, 1, true},
		{"Func calls - [DoAgain]. Expect two", []int{DoAgain}, 2, false},

		{"Func calls - [Do, DoAgain]. Expect one", []int{Do, DoAgain}, 1, false},
		{"Func calls - [Do, DoAgain]. Expect two", []int{Do, DoAgain}, 2, true},
		{"Func calls - [DoAgain, Do]. Expect one", []int{DoAgain, Do}, 1, true},
		{"Func calls - [DoAgain, Do]. Expect two", []int{DoAgain, Do}, 2, false},
		{"Func calls - [DoAgain, DoAgain]. Expect one", []int{DoAgain, DoAgain}, 1, false},
		{"Func calls - [DoAgain, DoAgain]. Expect two", []int{DoAgain, DoAgain}, 2, true},

		{"Func calls - [DoAgain, Do, Do]. Expect one", []int{DoAgain, Do, Do}, 1, true},
		{"Func calls - [DoAgain, Do, Do]. Expect three", []int{DoAgain, Do, Do}, 3, false},

		{"Func calls - [Do, DoAgain, Do]. Expect one", []int{Do, DoAgain, Do}, 1, false},
		{"Func calls - [Do, DoAgain, Do]. Expect two", []int{Do, DoAgain, Do}, 2, true},
		{"Func calls - [Do, DoAgain, Do]. Expect two", []int{Do, DoAgain, Do}, 3, false},

		{"Func calls - [Do, DoAgain, DoAgain]. Expect three", []int{Do, DoAgain, DoAgain}, 3, true},
		{"Func calls - [Do, DoAgain, DoAgain]. Expect two", []int{Do, DoAgain, DoAgain}, 2, false},

		{"Func calls - [DoAgain, DoAgain, DoAgain]. Expect three", []int{DoAgain, DoAgain, DoAgain}, 3, true},
		{"Func calls - [DoAgain, DoAgain, DoAgain]. Expect one", []int{DoAgain, DoAgain, DoAgain}, 1, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := new(counter)
			o := new(sync.Once)

			for i := 0; i < len(tc.funCalls); i++ {
				switch tc.funCalls[i] {
				case Do:
					o.Do(func() { c.Increment() })
				case DoAgain:
					o.DoAgain(func() { c.Increment() })
				}
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
