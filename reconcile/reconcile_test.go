package reconcile

import (
	"crypto/sha512"
	"fmt"
	"testing"
)

type object struct {
	i interface{}
}

func (o object) Hash() []byte {
	h := sha512.New()
	fmt.Fprintf(h, "%s", o.i)
	return h.Sum(nil)
}

type testState struct {
	i map[string]interface{}
}

func strptr(s string) *string                          { return &s }
func (ts *testState) Add(key string, v interface{})    { ts.i[key] = v }
func (ts *testState) Update(key string, v interface{}) { ts.i[key] = v }
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
		updates          []update
	}{
		{
			name: "string add",
			updates: []update{
				{
					key:   "Hello",
					state: new,
				},
			},

			current: &testState{
				map[string]interface{}{},
			},
			desired: &testState{
				map[string]interface{}{
					"Hello": "World",
				},
			},
		},
		{
			name: "string update",
			updates: []update{
				{
					key:   "Hello",
					state: dirty,
				},
			},
			current: &testState{
				map[string]interface{}{
					"Hello": "Go",
				},
			},
			desired: &testState{
				map[string]interface{}{
					"Hello": "World",
				},
			},
		},
		{
			name: "string delete",
			updates: []update{
				{
					key:   "Hello",
					state: old,
				},
			},
			current: &testState{
				map[string]interface{}{
					"Hello": "Go",
				},
			},
			desired: &testState{
				map[string]interface{}{},
			},
		},
		{
			name:    "string noop",
			updates: []update{},
			current: &testState{
				map[string]interface{}{
					"Hello": "World",
				},
			},
			desired: &testState{
				map[string]interface{}{
					"Hello": "World",
				},
			},
		},
		{
			name:    "object noop",
			updates: []update{},
			current: &testState{
				map[string]interface{}{
					"Hello": object{true},
				},
			},
			desired: &testState{
				map[string]interface{}{
					"Hello": object{true},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updates := diff(test.current, test.desired)
			if len(test.updates) != len(updates) {
				t.Log("actual updates don't match the expected ones")
				return
			}
			for i, update := range test.updates {
				keyMatch := updates[i].key == update.key
				stateMatch := updates[i].state == update.state

				if !(keyMatch && stateMatch) {
					t.Fail()
				}
			}
		})
	}
}
