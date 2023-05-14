package cmd

import (
	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

// NewServeCommand creates and returns new command responsible for
// starting the default PocketBase web server.
func NewSQLiteCommand(app core.App) *cobra.Command {

	command := &cobra.Command{
		Use:   "sqlite",
		Short: "Run arbitrary SQL commands on the database",
		Run: func(command *cobra.Command, args []string) {
			//Checks to make sure a query was provided
			if len(args) == 0 {
				color.Red("Error: No query provided")
				return
			}
			//Runs the SQLite command
			query := args[0]
			result, err := app.Dao().DB().NewQuery(query).Execute()
			if err != nil {
				color.Red("Error: %v", err)
				return
			}
			//Outputs how many rows were affected
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				color.Yellow("Error: %v", err)
				return
			}
			color.Green("Sucess! Affected %v row(s)", rowsAffected)
		},
	}

	return command
}
