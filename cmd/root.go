package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	services "marco-souza/zup/services"
)

func Execute() {
	cli := initRootCmd()
	cli.Execute()
}

func initRootCmd() *cobra.Command {
	cli := &cobra.Command{
		Use:   "zup",
		Short: "Zsh setUP",
		Long:  "cli tool to setup your ZSH terminal.",
		Run:   rootHandler,
		Args:  valdiateFlags,
	}

	// flags
	setupRootFlags(cli)

	return cli
}

func rootHandler(cmd *cobra.Command, args []string) {
	system, _ := cmd.Flags().GetString("system")
	dest, _ := cmd.Flags().GetString("output")
	supportJava, _ := cmd.Flags().GetBool("java")

	services.CreateZshFiles(dest, system, supportJava)

	fmt.Println("Done, you are ready to go! 🐹")
	fmt.Println("You can start a new zsh session now :)")
}

func setupRootFlags(cmd *cobra.Command) {
	homeFolder := os.Getenv("HOME")
	persistentFlags := cmd.PersistentFlags()
	usage := fmt.Sprintf("Select operational system: [%s]", strings.Join(services.SystemOptions, "|"))

	persistentFlags.StringP("system", "s", "arch", usage) // TODO: add get default OS
	persistentFlags.StringP("output", "o", homeFolder, "Destination folder")
	persistentFlags.BoolP("java", "j", false, "Enable Java Support (with SDKMan)")
}

func valdiateFlags(cmd *cobra.Command, args []string) error {
	if err := validateSystem(cmd); err != nil {
		return err
	}
	if err := validateDestination(cmd); err != nil {
		return err
	}
	return nil
}

func validateSystem(cmd *cobra.Command) error {
	system, _ := cmd.Flags().GetString("system")
	if services.IsSystemOption(system) {
		return nil
	}
	return fmt.Errorf("%s is not a valid operational system", system)
}

func validateDestination(cmd *cobra.Command) error {
	dest, _ := cmd.Flags().GetString("output")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		if err := os.Mkdir(dest, 0755); err != nil {
			return err
		}
	}
	return nil
}
