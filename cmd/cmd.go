package cmd

import (
	_ "k8s.io/klog/v2"
)

// Execute executes the commands.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
