// Copyright 2018 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reconcile

// StateWalkFunc walks a state. Walk should be hierarchical to ensure
// no cascading updates occur.
type StateWalkFunc func(key string, v interface{})

// State represents the interface which Reconciler Accepts
type State interface {
	Add(key string, v interface{})
	Update(key string, v interface{})
	Get(key string) interface{}
	Delete(key string)
	Walk(f StateWalkFunc)
}

type state int
type mode int

const (
	new state = iota
	old
	dirty

	first mode = iota
	middle
	last
)

// Reconcile takes two states and applies updates to them until they are the same
func Reconcile(current, desired State) {
	fix(current, diff(current, desired))
}

type update struct {
	key   string
	state state
	v     interface{}
}

func diff(current, desired State) []*update {
	var updates []*update
	desired.Walk(func(key string, v interface{}) {
		currentValue := current.Get(key)
		n := update{
			key: key,
			v:   v,
		}
		switch currentValue {
		case nil:
			n.state = new
		case v:
			return
		default:
			n.state = dirty
		}

		updates = append(updates, &n)
	})

	current.Walk(func(key string, v interface{}) {
		if desired.Get(key) == nil {
			n := update{
				key: key,
				v:   nil,
			}
			n.state = old
			updates = append(updates, &n)
		}
	})

	return updates
}

func fix(current State, updates []*update) {
	for _, update := range updates {
		switch update.state {
		case new:
			current.Add(update.key, update.v)
		case old:
			current.Delete(update.key)
		case dirty:
			current.Update(update.key, update.v)
		}
	}
}
