// Copyright 2015 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time // import "sevki.org/lib/time"

import (
	"testing"
	"time"
)

func TestThreeMins(t *testing.T) {
	r := NsReadable(int64(time.Minute.Nanoseconds() * 3))
	if r != "3.00m" {
		t.Log(r)
		t.Fail()
	}
}
func TestOneTenthOfMin(t *testing.T) {
	r := NsReadable(int64(time.Minute.Nanoseconds() / 10))
	if r != "0.10m" {
		t.Log(r)
		t.Fail()
	}
}

func TestThreeHundrentsOfMin(t *testing.T) {
	r := NsReadable(int64(time.Minute.Nanoseconds() / 100 * 3))
	if r != "0.03m" {
		t.Log(r)
		t.Fail()
	}
}
func TestOneSecond(t *testing.T) {
	r := NsReadable(int64(time.Second.Nanoseconds()))
	if r != "0.02m" { // rounded up from 0.0166667
		t.Log(r)
		t.Fail()
	}
}
func TestOneTenthsOfSecond(t *testing.T) {
	r := NsReadable(int64(time.Second.Nanoseconds() / 10))
	if r != "0.10s" {
		t.Log(r)
		t.Fail()
	}
}
func TestNineTenthsOfSecond(t *testing.T) {
	r := NsReadable(int64(time.Second.Nanoseconds()/10) * 9)
	if r != "0.01m" {
		t.Log(r)
		t.Fail()
	}
}
func TestNineOneHundrenthOfSecond(t *testing.T) {
	r := NsReadable(int64(time.Second.Nanoseconds()/100) * 9)
	if r != "0.09s" {
		t.Log(r)
		t.Fail()
	}
}
