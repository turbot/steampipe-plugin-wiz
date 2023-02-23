package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/wiz"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: wiz.Plugin})
}
