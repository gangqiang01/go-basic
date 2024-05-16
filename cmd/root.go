package cmd

import (
	"os"
	"path/filepath"
	"strings"

	_ "github.com/edgehook/ithings/common/dbm"
	"github.com/edgehook/ithings/webserver"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:     "ithings",
	Long:    `iot things manager.. `,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: To help debugging, immediately log version
		klog.Infof("###########  Start the ithings...! ###########")
		registerModules()
		// start all modules
		core.Run()
	},
}

func getRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}

	return strings.Replace(dir, "\\", "/", -1)
}

func addServerFlags(flags *pflag.FlagSet) {
	flags.StringP("usr", "u", "", "mqtt user name.")
	flags.StringP("passwd", "p", "", "mqtt passwd.")
}

func init() {
	flags := rootCmd.Flags()
	persistent := rootCmd.PersistentFlags()

	persistent.StringP("database", "d", filepath.Join(getRootPath(), "ithings.db"), "database path")
	addServerFlags(flags)
}

// register all module into beehive.
func registerModules() {
	webserver.Register()
}
