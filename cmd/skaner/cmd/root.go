package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/marcelo-rocha/skaner/domain/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	options    runner.Options
	cfgFile    string
	mainLogger *zap.Logger

	rootCmd = &cobra.Command{
		Use:   "skaner file1 file2 ...",
		Short: "A security source scanner",
		Long: `Skaner is a CLI application that scan source code, and text files, looking for vulnerabilities
		and sensitive data exposures`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a file argument")
			}
			for _, n := range args {
				if ok, _ := fileExists(n); !ok {
					return errors.New("file not found: " + n)
				}
			}
			if len(options.SensitiveText) == 0 && !options.DisableExposureCheck {
				return errors.New("sensitive strings must be defined")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			runner.Run(context.Background(), os.Stdout, args, options, mainLogger)
			return nil
		},
	}
)

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func Execute(logger *zap.Logger) {
	mainLogger = logger
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&options.JsonOutput, "json", "j", false, "output vulnerabilities in JSON format")
	rootCmd.PersistentFlags().IntVar(&options.WorkersQty, "workers", 3, "number of concurrent checkers")

	rootCmd.Flags().BoolVar(&options.DisableExposureCheck, "no-exposure-checker", false, "disable checking of sensitive data exposure")
	rootCmd.Flags().BoolVar(&options.DisableSQLCheck, "no-sql-injection-checker", false, "disable checking of SQL injection")
	rootCmd.Flags().BoolVar(&options.DisableXSSCheck, "no-xss-checker", false, "disable checking of cross site scripting")
	rootCmd.Flags().StringSliceVarP(&options.SensitiveText, "sensitive-text", "s", []string{}, "inform a comma separate list of strings")

	//viper.BindPFlag("json-format", rootCmd.PersistentFlags().Lookup("json"))
	//viper.BindPFlag("sensitive-text", rootCmd.Flags().Lookup("sensitive-text"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	if cfgFile != "" {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintln(os.Stderr, "config file error:", err)
			os.Exit(1)
		}
	}
	if ss := viper.GetStringSlice("sensitive-text"); len(ss) > 0 {
		if len(options.SensitiveText) == 0 {
			options.SensitiveText = ss
		}
	}

	if viper.GetBool("json-format") {
		options.JsonOutput = true
	}
}
