/*
Copyright © 2023 Yusuke Yagyu <yu@hoaxster.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/csv"
	"os"
	"strings"

	"github.com/gyugyu/pasps/cmd/helper"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/spf13/cobra"
)

var padding int

var rootCmd = &cobra.Command{
	Use:   "pasps",
	Args:  cobra.ExactArgs(3),
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := helper.GetClient()
		if err != nil {
			return err
		}

		ctx := context.Background()
		srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
		spreadsheetId := args[0]
		appendRange := args[1]

		dataReader := strings.NewReader(args[2])
		data, err := csv.NewReader(dataReader).ReadAll()
		if err != nil {
			return err
		}

		_, err = srv.Spreadsheets.Values.Append(spreadsheetId, appendRange, &sheets.ValueRange{
			Values: helper.CreateValues(data, padding),
		}).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
		if err != nil {
			return err
		}

		return nil
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pasps.yaml)")

	rootCmd.Flags().BoolP("horizontal", "", false, "horizontal")
	rootCmd.Flags().IntVarP(&padding, "padding", "p", 1, "number of padding")
}
