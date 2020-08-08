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

func TestOnce_DoForce(t *testing.T) {
	const (
		Do = iota
		DoForce
	)
	tt := []struct {
		name       string
		funCalls   []int
		val        int
		shouldPass bool
	}{
		{"Func calls - [DoForce]. Expect one", []int{DoForce}, 1, true},
		{"Func calls - [DoForce]. Expect two", []int{DoForce}, 2, false},

		{"Func calls - [Do, DoForce]. Expect one", []int{Do, DoForce}, 1, false},
		{"Func calls - [Do, DoForce]. Expect two", []int{Do, DoForce}, 2, true},
		{"Func calls - [DoForce, Do]. Expect one", []int{DoForce, Do}, 1, true},
		{"Func calls - [DoForce, Do]. Expect two", []int{DoForce, Do}, 2, false},
		{"Func calls - [DoForce, DoForce]. Expect one", []int{DoForce, DoForce}, 1, false},
		{"Func calls - [DoForce, DoForce]. Expect two", []int{DoForce, DoForce}, 2, true},

		{"Func calls - [DoForce, Do, Do]. Expect one", []int{DoForce, Do, Do}, 1, true},
		{"Func calls - [DoForce, Do, Do]. Expect three", []int{DoForce, Do, Do}, 3, false},

		{"Func calls - [Do, DoForce, Do]. Expect one", []int{Do, DoForce, Do}, 1, false},
		{"Func calls - [Do, DoForce, Do]. Expect two", []int{Do, DoForce, Do}, 2, true},
		{"Func calls - [Do, DoForce, Do]. Expect two", []int{Do, DoForce, Do}, 3, false},

		{"Func calls - [Do, DoForce, DoForce]. Expect three", []int{Do, DoForce, DoForce}, 3, true},
		{"Func calls - [Do, DoForce, DoForce]. Expect two", []int{Do, DoForce, DoForce}, 2, false},

		{"Func calls - [DoForce, DoForce, DoForce]. Expect three", []int{DoForce, DoForce, DoForce}, 3, true},
		{"Func calls - [DoForce, DoForce, DoForce]. Expect one", []int{DoForce, DoForce, DoForce}, 1, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := new(counter)
			o := new(sync.Once)

			for i := 0; i < len(tc.funCalls); i++ {
				switch tc.funCalls[i] {
				case Do:
					o.Do(func() { c.Increment() })
				case DoForce:
					o.DoForce(func() { c.Increment() })
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

func TestOnce_Reset(t *testing.T) {
	const (
		Do = iota
		Reset
	)
	tt := []struct {
		name       string
		funCalls   []int
		val        int
		shouldPass bool
	}{
		{"Func calls - [Do, Do]. Expect one", []int{Do, Do}, 1, true},
		{"Func calls - [Do, Do]. Expect two", []int{Do, Do}, 2, false},
		{"Func calls - [Do, Reset, Do]. Expect two", []int{Do, Reset, Do}, 2, true},
		{"Func calls - [Do, Reset, Do]. Expect one", []int{Do, Reset, Do}, 1, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := new(counter)
			o := new(sync.Once)

			for i := 0; i < len(tc.funCalls); i++ {
				switch tc.funCalls[i] {
				case Do:
					o.Do(func() { c.Increment() })
				case Reset:
					o.Reset()
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
