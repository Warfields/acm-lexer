/*
Copyright © 2024 Sam Warfield <swarfield@todyl.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	acm "github.com/Warfields/acm-lexer/parser"
	"github.com/antlr4-go/antlr/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "acm",
	Short: "Parses \"ACM Filter files\"",
	Long:  `Demonstrates a basic parser to parse "ACM Filter Files"`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: parseFile,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.acm-lexer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".acm-lexer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".acm-lexer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func parseFile(cmd *cobra.Command, args []string) {
	lines, err := readFileLines(args[0])
	if err != nil {
		log.Fatal(err, args[0])
	}
	for _, line := range lines {
		// Skip empty lines
		if line == "" {
			continue
		}
		lexer := acm.NewAcmLexer(antlr.NewInputStream(line))
		tokStream := antlr.NewCommonTokenStream(lexer, 0)
		parser := acm.NewAcmParser(tokStream)

		parser.AddErrorListener(antlr.NewConsoleErrorListener())

		// Convert input to parse tree
		parseTree := parser.Filter()

		// Walk the tree to get valuable information
		listener := acm.NewAcmFieldListener()
		antlr.ParseTreeWalkerDefault.Walk(listener, parseTree)

		fmt.Println("Filter:", line)
		fmt.Println("Fields:", listener.GetFields())
		fmt.Println("Values:", listener.GetValues())
		fmt.Println("---")
	}
}

func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
