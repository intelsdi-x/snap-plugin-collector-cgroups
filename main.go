// +build linux

package main

import (
	"fmt"
	"os"

	"github.com/intelsdi-x/snap-plugin-collector-cgroups/cgroups"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {

	if cgplugin, err := cgroups.NewCgroups(false); err != nil {
		fmt.Printf("Error: %s \n", err)
		return
	} else {
		plugin.Start(plugin.NewPluginMeta(
			cgroups.NS_PLUGIN,
			cgroups.VERSION,
			plugin.CollectorPluginType,
			[]string{},
			[]string{plugin.SnapGOBContentType},
			plugin.ConcurrencyCount(1)),
			cgplugin,
			os.Args[1])
	}

}
