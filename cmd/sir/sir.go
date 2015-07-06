// Main SOON Instance Registry Package Entrypoint

package main

import (
	"github.com/spf13/cobra"

	"github.com/thisissoon/sir"
)

// Long Description
var sirCobraCmsDesc = `SOON_ Instance Rregistry is a simple webservice
that allows instances to obtain a hostname and register that hostname
as well as deregistry hostnames when an instance terminates.`

// Main command entry point
var SirCobraCmd = &cobra.Command{
	Use:   "sir",
	Short: "SOON_ Instance Registry Service",
	Long:  sirCobraCmsDesc,
	Run: func(cmd *cobra.Command, args []string) {
		sir.Serve()
	},
}

// Main function
func main() {
	SirCobraCmd.Execute()
}
