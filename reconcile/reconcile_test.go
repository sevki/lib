package reconcile

import (
	"testing"
)

type testState struct {
	i map[string]string
}

func strptr(s string) *string                          { return &s }
func (ts *testState) Add(key string, v interface{})    { ts.i[key] = v.(string) }
func (ts *testState) Update(key string, v interface{}) { ts.i[key] = v.(string) }
func (ts *testState) Get(key string) interface{} {
	if val, ok := ts.i[key]; ok {
		return val
	}
	return nil
}
func (ts *testState) Delete(key string) { delete(ts.i, key) }
func (ts *testState) Walk(f StateWalkFunc) {
	for k, v := range ts.i {
		f(k, v)
	}
}

func TestReconcile(t *testing.T) {
	tests := []struct {
		name             string
		current, desired State
	}{
		{
			name: "map",
			current: &testState{
				map[string]string{},
			},
			desired: &testState{
				map[string]string{
					"Hello": "World",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Reconcile(test.current, test.desired)
			test.desired.Walk(func(key string, v interface{}) {
				if test.current.Get(key) != v {
					t.Fail()
					return
				}
			})
		})
	}
}
