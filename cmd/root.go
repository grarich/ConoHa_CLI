// Package cmd
/*
Copyright © 2023 grarich <grarich@grawlily.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"conoha_cli/cmd/identity"
)

// rootCmd represents the base command when called without any subcommands
var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "ConoHa_CLI",
		Short: "ConoHa API Client",
		Long: `ConoHa API をCLIから操作できます。
主にダッシュボードから操作できない設定をするために開発されました。`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		rootCmd.SetErr(os.Stderr)
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 設定ファイルを作成
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.conoha_cli.yaml)")

	initCmd()
}

// 設定ファイルの読み込み
func initConfig() {
	// 設定ファイルのパスが渡された場合
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".conoha_cli")
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		errMsg := fmt.Errorf("設定ファイルの読み込みに失敗しました。")
		fmt.Printf(errMsg.Error())
	}
}

// 各コマンドのロード
func initCmd() {
	rootCmd.AddCommand(identity.Cmd)
}
