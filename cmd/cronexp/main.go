// Copyright 2019 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cronexp draws a very nice table of what your cron job will do
// like so
//
//	 ┌───────────────────minutes           [0s 15m0s 30m0s 45m0s]
//	 │    ┌──────────────hours             [0s]
//	 │    │ ┌────────────days of the month [1 15]
//	 │    │ │    ┌───────months            [January February March April May June July August September October November December]
//	 │    │ │    │ ┌─────days of the week  [Monday Tuesday Wednesday Thursday Friday]
//	 │    │ │    │ │   ┌─args              [/usr/bin/find]
//	 */15 0 1,15 * 1-5 /usr/bin/find
//
// Important '*' is a wildcard character in most shells.
// either escape it with \* or run with
//
//	set -f
//	cronexp */15 0 1,15 * 1-5 /usr/bin/find
//
package main

import (
	"flag"
	"fmt"

	"sevki.org/x/cron"
)

func main() {
	fset := cron.NewFlagSet("test", flag.PanicOnError)
	fset.Parse()
	fmt.Println(fset.CronTab())
}
