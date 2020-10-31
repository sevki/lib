package cron

import (
	"flag"
	"fmt"
	"strings"
)

func ExampleFlagSet_CronTab() {
	fset := NewFlagSet("test", flag.PanicOnError)
	fset.Parse(strings.Fields("*/15 0 1,15 * 1-5 /usr/bin/find")...)
	fmt.Println(fset.CronTab())

	// Output:
	// ┌───────────────────minutes           [0s 15m0s 30m0s 45m0s]
	// │    ┌──────────────hours             [0s]
	// │    │ ┌────────────days of the month [1 15]
	// │    │ │    ┌───────months            [January February March April May June July August September October November December]
	// │    │ │    │ ┌─────days of the week  [Monday Tuesday Wednesday Thursday Friday]
	// │    │ │    │ │   ┌─args              [/usr/bin/find]
	// */15 0 1,15 * 1-5 /usr/bin/find

}
