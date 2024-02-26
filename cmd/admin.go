package cmd

import (
	"errors"
	"fmt"

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
		Use:          "create",
		Example:      "admin create test@example.com 1234567890",
		Short:        "Creates a new admin account",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Missing email and password arguments.")
			}

			if args[0] == "" || is.EmailFormat.Validate(args[0]) != nil {
				return errors.New("Missing or invalid email address.")
			}

			if len(args[1]) < 8 {
				return errors.New("The password must be at least 8 chars long.")
			}

			admin := &models.Admin{}
			admin.Email = args[0]
			admin.SetPassword(args[1])

			if !app.Dao().HasTable(admin.TableName()) {
				return errors.New("Migration are not initialized yet. Please run 'migrate up' and try again.")
			}

			if err := app.Dao().SaveAdmin(admin); err != nil {
				return fmt.Errorf("Failed to create new admin account: %v", err)
			}

			color.Green("Successfully created new admin %s!", admin.Email)
			return nil
		},
	}

	return command
}

func adminUpdateCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:          "update",
		Example:      "admin update test@example.com 1234567890",
		Short:        "Changes the password of a single admin account",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Missing email and password arguments.")
			}

			if args[0] == "" || is.EmailFormat.Validate(args[0]) != nil {
				return errors.New("Missing or invalid email address.")
			}

			if len(args[1]) < 8 {
				return errors.New("The new password must be at least 8 chars long.")
			}

			if !app.Dao().HasTable((&models.Admin{}).TableName()) {
				return errors.New("Migration are not initialized yet. Please run 'migrate up' and try again.")
			}

			admin, err := app.Dao().FindAdminByEmail(args[0])
			if err != nil {
				return fmt.Errorf("Admin with email %s doesn't exist.", args[0])
			}

			admin.SetPassword(args[1])

			if err := app.Dao().SaveAdmin(admin); err != nil {
				return fmt.Errorf("Failed to change admin %s password: %v", admin.Email, err)
			}

			color.Green("Successfully changed admin %s password!", admin.Email)
			return nil
		},
	}

	return command
}

func adminDeleteCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:          "delete",
		Example:      "admin delete test@example.com",
		Short:        "Deletes an existing admin account",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 0 || args[0] == "" || is.EmailFormat.Validate(args[0]) != nil {
				return errors.New("Invalid or missing email address.")
			}

			if !app.Dao().HasTable((&models.Admin{}).TableName()) {
				return errors.New("Migration are not initialized yet. Please run 'migrate up' and try again.")
			}

			admin, err := app.Dao().FindAdminByEmail(args[0])
			if err != nil {
				color.Yellow("Admin %s is already deleted.", args[0])
				return nil
			}

			if err := app.Dao().DeleteAdmin(admin); err != nil {
				return fmt.Errorf("Failed to delete admin %s: %v", admin.Email, err)
			}

			color.Green("Successfully deleted admin %s!", admin.Email)
			return nil
		},
	}

	return command
}
