package main

import (
	"fmt"
	"os"
	"sort"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()

	app.Version = "greet 1.0.0"

	flagMap := make(map[string]cli.Flag)
	flagMap["lang"] = cli.StringFlag{Name: "lang"}
	flagMap["config"] = cli.StringFlag{Name: "config"}
	flagMap["networkid"] = cli.Int64Flag{Name: "networkid", Value: 123}

	for _, flagValue := range flagMap {
		app.Flags = append(app.Flags, flagValue)
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("==== Action ====")
		fmt.Println("NumFlags :", c.NumFlags())
		fmt.Println("FlagNames :", len(c.FlagNames()))
		for _, v := range c.FlagNames() {
			fmt.Println(v)
		}

		fmt.Println(">> ", len(c.GlobalFlagNames()))
		for _, name := range c.GlobalFlagNames() {
			//if c.IsSet(name)
			{
				fmt.Println(name, "is set:", c.String(name))
			}
		}

		fmt.Println("NArg :", c.NArg())
		if c.NArg() > 0 {
			for _, arg := range c.Args() {
				fmt.Println(arg)
			}
		}

		cmd := c.Command
		if cmd.Name != "" {
			fmt.Printf("Command:%#v\n", cmd)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: MigrateFlags(func(c *cli.Context) error {
				fmt.Println("--complete")
				return nil
			}),
			Flags: []cli.Flag{
				flagMap["config"],
				flagMap["networkid"],
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task on the list",
			Action: MigrateFlags(func(c *cli.Context) error {
				fmt.Println("--add")
				return nil
			}),
			Flags: []cli.Flag{
				flagMap["lang"],
				flagMap["networkid"],
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		fmt.Println("==== Before ====")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Println("==== After ====")
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {

	return func(ctx *cli.Context) error {
		fmt.Println("[ MigrateFlags ]")

		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))

				fmt.Println(name, "is set:", ctx.String(name))
			}
		}

		fmt.Println("GlobalFlagNames>> ", len(ctx.GlobalFlagNames()))
		for _, name := range ctx.GlobalFlagNames() {
			if ctx.GlobalIsSet(name) {
				fmt.Println(name, "is set:", ctx.GlobalString(name))
			}
		}
		return action(ctx)
	}
}
