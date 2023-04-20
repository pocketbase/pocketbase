package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"
)

// NewAdminCommand creates and returns new command for managing
// admin accounts (create, update, delete).
func NewAdminCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "admin",
		Short: "Manages admin accounts",
	}

	command.AddCommand(adminCreateCommand(app))
	command.AddCommand(adminUpdateCommand(app))
	command.AddCommand(adminDeleteCommand(app))

	return command
}

func adminCreateCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "create",
		Example: "admin create test@example.com 1234567890",
		Short:   "Creates a new admin account",
		Run: func(command *cobra.Command, args []string) {
			if len(args) != 2 {
				color.Red("Missing email and password arguments.")
				os.Exit(1)
			}

			if is.EmailFormat.Validate(args[0]) != nil {
				color.Red("Invalid email address.")
				os.Exit(1)
			}

			if len(args[1]) < 8 {
				color.Red("The password must be at least 8 chars long.")
				os.Exit(1)
			}

			admin := &models.Admin{}
			admin.Email = args[0]
			admin.SetPassword(args[1])

			if err := app.Dao().SaveAdmin(admin); err != nil {
				color.Red("Failed to create new admin account: %v", err)
				os.Exit(1)
			}

			color.Green("Successfully created new admin %s!", admin.Email)
		},
	}

	return command
}

func adminUpdateCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "update",
		Example: "admin update test@example.com 1234567890",
		Short:   "Changes the password of a single admin account",
		Run: func(command *cobra.Command, args []string) {
			if len(args) != 2 {
				color.Red("Missing email and password arguments.")
				os.Exit(1)
			}

			if is.EmailFormat.Validate(args[0]) != nil {
				color.Red("Invalid email address.")
				os.Exit(1)
			}

			if len(args[1]) < 8 {
				color.Red("The new password must be at least 8 chars long.")
				os.Exit(1)
			}

			admin, err := app.Dao().FindAdminByEmail(args[0])
			if err != nil {
				color.Red("Admin with email %s doesn't exist.", args[0])
				os.Exit(1)
			}

			admin.SetPassword(args[1])

			if err := app.Dao().SaveAdmin(admin); err != nil {
				color.Red("Failed to change admin %s password: %v", admin.Email, err)
				os.Exit(1)
			}

			color.Green("Successfully changed admin %s password!", admin.Email)
		},
	}

	return command
}

func adminDeleteCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "delete",
		Example: "admin delete test@example.com",
		Short:   "Deletes an existing admin account",
		Run: func(command *cobra.Command, args []string) {
			if len(args) == 0 || is.EmailFormat.Validate(args[0]) != nil {
				color.Red("Invalid or missing email address.")
				os.Exit(1)
			}

			admin, err := app.Dao().FindAdminByEmail(args[0])
			if err != nil {
				color.Yellow("Admin %s is already deleted.", args[0])
				return
			}

			if err := app.Dao().DeleteAdmin(admin); err != nil {
				color.Red("Failed to delete admin %s: %v", admin.Email, err)
				os.Exit(1)
			}

			color.Green("Successfully deleted admin %s!", admin.Email)
		},
	}

	return command
}
