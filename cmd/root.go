package cmd

import (
	"github.com/edgehook/ithings/webserver"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:     "broadcast",
	Long:    `iot broadcast manager.. `,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: To help debugging, immediately log version
		klog.Infof("###########  Start the broadcast...! ###########")
		klog.Infof("args: %v", args)
		registerModules()
		// start all modules
		core.Run()
	},
}

func init() {

	persistent := rootCmd.PersistentFlags()
	persistent.StringVarP(&webserver.BindAddress, "address", "", ":9001", "web address")
}

// register all module into beehive.
func registerModules() {
	webserver.Register()
}
