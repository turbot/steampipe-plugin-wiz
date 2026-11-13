package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-wiz/wiz"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: wiz.Plugin})
}
