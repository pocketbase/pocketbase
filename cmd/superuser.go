package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

// NewSuperuserCommand creates and returns new command for managing
// superuser accounts (create, update, upsert, delete).
func NewSuperuserCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "superuser",
		Short: "Manages superuser accounts",
	}

	command.AddCommand(superuserUpsertCommand(app))
	command.AddCommand(superuserCreateCommand(app))
	command.AddCommand(superuserUpdateCommand(app))
	command.AddCommand(superuserDeleteCommand(app))

	return command
}

func superuserUpsertCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:          "upsert",
		Example:      "superuser upsert test@example.com 1234567890",
		Short:        "Creates, or updates if email exists, a single superuser account",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Missing email and password arguments.")
			}

			if args[0] == "" || is.EmailFormat.Validate(args[0]) != nil {
				return errors.New("Missing or invalid email address.")
			}

			superusersCol, err := app.FindCachedCollectionByNameOrId(core.CollectionNameSuperusers)
			if err != nil {
				return fmt.Errorf("Failed to fetch %q collection: %w.", core.CollectionNameSuperusers, err)
			}

			superuser, err := app.FindAuthRecordByEmail(superusersCol, args[0])
			if err != nil {
				superuser = core.NewRecord(superusersCol)
			}

			superuser.SetEmail(args[0])
			superuser.SetPassword(args[1])

			if err := app.Save(superuser); err != nil {
				return fmt.Errorf("Failed to upsert superuser account: %w.", err)
			}

			color.Green("Successfully saved superuser %q!", superuser.Email())
			return nil
		},
	}

	return command
}

func superuserCreateCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Example:      "superuser create test@example.com 1234567890",
		Short:        "Creates a new superuser account",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Missing email and password arguments.")
			}

			if args[0] == "" || is.EmailFormat.Validate(args[0]) != nil {
				return errors.New("Missing or invalid email address.")
			}

			superusersCol, err := app.FindCachedCollectionByNameOrId(core.CollectionNameSuperusers)
			if err != nil {
				return fmt.Errorf("Failed to fetch %q collection: %w.", core.CollectionNameSuperusers, err)
			}

			superuser := core.NewRecord(superusersCol)
			superuser.SetEmail(args[0])
			superuser.SetPassword(args[1])

			if err := app.Save(superuser); err != nil {
				return fmt.Errorf("Failed to create new superuser account: %w.", err)
			}

			color.Green("Successfully created new superuser %q!", superuser.Email())
			return nil
		},
	}

	return command
}

func superuserUpdateCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:          "update",
		Example:      "superuser update test@example.com 1234567890",
		Short:        "Changes the password of a single superuser account",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Missing email and password arguments.")
			}

			if args[0] == "" || is.EmailFormat.Validate(args[0]) != nil {
				return errors.New("Missing or invalid email address.")
			}

			superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, args[0])
			if err != nil {
				return fmt.Errorf("Superuser with email %q doesn't exist.", args[0])
			}

			superuser.SetPassword(args[1])

			if err := app.Save(superuser); err != nil {
				return fmt.Errorf("Failed to change superuser %q password: %w.", superuser.Email(), err)
			}

			color.Green("Successfully changed superuser %q password!", superuser.Email())
			return nil
		},
	}

	return command
}

func superuserDeleteCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:          "delete",
		Example:      "superuser delete test@example.com",
		Short:        "Deletes an existing superuser account",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 0 || args[0] == "" || is.EmailFormat.Validate(args[0]) != nil {
				return errors.New("Invalid or missing email address.")
			}

			superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, args[0])
			if err != nil {
				color.Yellow("Superuser %q is missing or already deleted.", args[0])
				return nil
			}

			if err := app.Delete(superuser); err != nil {
				return fmt.Errorf("Failed to delete superuser %q: %w.", superuser.Email(), err)
			}

			color.Green("Successfully deleted superuser %q!", superuser.Email())
			return nil
		},
	}

	return command
}
