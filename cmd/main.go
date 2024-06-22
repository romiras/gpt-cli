package main

import (
	"fmt"
	"os"

	"github.com/romiras/gpt-cli/internal/api"
	"github.com/romiras/gpt-cli/internal/models"
	"github.com/romiras/gpt-cli/internal/services"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gpt-cli",
		Short: "CLI for GPT via API",
	}

	cfg := models.BuildDefaultConfig()

	flagSet := rootCmd.PersistentFlags()
	flagSet.StringVarP(&cfg.Model, "model", "m", cfg.Model, "Model to use")
	flagSet.IntVarP(&cfg.MaxTokens, "max-tokens", "t", cfg.MaxTokens, "Maximum tokens in response")
	flagSet.Float64VarP(&cfg.Temperature, "temperature", "p", cfg.Temperature, "Response temperature")
	flagSet.Float64VarP(&cfg.TopP, "top-p", "P", cfg.TopP, "Top P response probability")

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		provider, err := api.NewHyperbolicProvider(&cfg)
		if err != nil {
			return err
		}

		return services.Run(provider)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
