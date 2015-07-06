// Main SOON Instance Registry Package Entrypoint
// This provides the CLI interface to sir.

package main

import (
	"github.com/spf13/cobra"

	"github.com/thisissoon/sir"
)

var (
	RedisAddr      string
	ImportFilePath string
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

// Import sub command
var SirImportCobraCommand = &cobra.Command{
	Use:   "import",
	Short: "Import hostnames from txt file",
	Long:  "Reads a line delimited text file of host names into Redis set",
	Run: func(cmd *cobra.Command, args []string) {
		sir.Import(&ImportFilePath, &RedisAddr)
	},
}

// Main function
func main() {
	SirCobraCmd.PersistentFlags().StringVarP(&RedisAddr, "redis", "r", "127.0.0.1:6379", "Redis Server Address (ip:port)")
	SirImportCobraCommand.Flags().StringVarP(&ImportFilePath, "path", "p", "", "Path to text file")

	SirCobraCmd.AddCommand(SirImportCobraCommand)
	SirCobraCmd.Execute()
}
