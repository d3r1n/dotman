package main

import (
	"github.com/5elenay/ezcli"
)

func Handle() {
	var handler *ezcli.CommandHandler = ezcli.NewApp("dotman")

	handler.AddCommand(&ezcli.Command{
		Name:        "init",
		Description: "initialize the \"dotman.json\" file",
		Options: []*ezcli.CommandOption{
			{
				Name:        "description",
				Description: "Set the description of the dotfiles (optional)",
				Aliases:     []string{"d", "desc"},
			},
			{
				Name:        "install_path",
				Description: "Set the installation path of the dotfiles (optional)",
				Aliases:     []string{"i", "ip"},
			},
			{
				Name:        "remote_url",
				Description: "Set the remote git url of the dotfiles (optional)",
			},
			{
				Name:        "branch_name",
				Description: "Set the installation path of the dotfiles (optional)",
			},
			{
				Name:        "git",
				Description: "Initiliaze a git repository for the dotfiles (default: false)",
				Aliases:     []string{"g"},
			},
		},
		Execute: func(c *ezcli.Command) {
			Init(c)
		},
	})

	handler.AddCommand(&ezcli.Command{
		Name:        "add",
		Description: "Add a dotfile ",
		Options: []*ezcli.CommandOption{
			{
				Name:        "description",
				Description: "Set the description of the dotfile (optional)",
				Aliases:     []string{"d", "desc"},
			},
			{
				Name:        "priority",
				Description: "Set the priority of the dotfile (optional)",
				Aliases:     []string{"p", "prio"},
			},
		},
		Execute: func(c *ezcli.Command) {
			Add(c)
		},
	})

	handler.AddCommand(&ezcli.Command{
		Name:        "remove",
		Description: "Remove a dotfile ",
		Aliases:     []string{"rm"},
		Execute: func(c *ezcli.Command) {
			Remove(c)
		},
	})

	handler.AddCommand(&ezcli.Command{
		Name:        "update",
		Description: "Update a dotfile ",
		Aliases:     []string{"up"},
		Execute: func(c *ezcli.Command) {
			Update(c)
		},
	})

	handler.AddCommand(&ezcli.Command{
		Name:        "command",
		Description: "Command utilities",
		Aliases:     []string{"cmd"},
		Options: []*ezcli.CommandOption{
			{
				Name:        "sudo",
				Description: "determine the sudo usage for the command",
				Aliases:     []string{"s"},
			},
		},
		Execute: func(c *ezcli.Command) {
			CommandHandler(c)
		},
	})

	handler.AddCommand(&ezcli.Command{
		Name:        "generate_installer",
		Description: "Generates Installer Script",
		Aliases:     []string{"gen"},
		Execute: func(c *ezcli.Command) {
			Installer(c)
		},
	})

	handler.AddCommand(&ezcli.Command{
		Name:        "status",
		Description: "Shows the status of the dotfiles",
		Aliases:     []string{"ss"},
		Execute: func(c *ezcli.Command) {
			Status(c)
		},
	})

	handler.Handle()
}
