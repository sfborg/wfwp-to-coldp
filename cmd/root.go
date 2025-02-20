/*
Copyright © 2024 Geoff Ower <gdower@illinois.edu>

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
	"log/slog"
	"os"
	"path/filepath"

	//"github.com/gnames/coldp/ent/coldp"
	"wfwp-to-coldp/internal/ent/wfwp"
	"wfwp-to-coldp/internal/io/wfwparcio"
	"wfwp-to-coldp/internal/io/sysio"
	fwfwp "wfwp-to-coldp/pkg"
	"wfwp-to-coldp/pkg/config"
	"github.com/sfborg/sflib/io/dbio"
	"wfwp-to-coldp/internal/io/schemaio"
	"github.com/spf13/cobra"
)

var opts []config.Option

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wfwp-to-coldp",
	Short: "Converts WF/WP archive file to COLDP archive",
	Long:  `Converts World Ferns/World Plants archive file to COLDP archive`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		versionFlag(cmd)
		flags := []flagFunc{
			debugFlag, cacheDirFlag, jobsNumFlag, binFlag, zipFlag, quotesFlag,
			fieldsNumFlag,
		}
		for _, v := range flags {
			v(cmd)
		}

		if len(args) != 2 {
			cmd.Help()
			os.Exit(0)
		}

		slog.Info("Converting WF/WP to COLDP")
		wfwpPath := args[0]
		outputPath := args[1]

		ext := filepath.Ext(outputPath)
		if ext == ".sqlite" {
			opts = append(opts, config.OptWithBinOutput(true))
		}

		cfg := config.New(opts...)
		err = sysio.New(cfg).ResetCache()
		if err != nil {
			slog.Error("Cannot initialize file system", "error", err)
			os.Exit(1)
		}

		wfwpSchema := schemaio.New(cfg.SchemaPath)
		wfwpDB := dbio.New(cfg.CacheWfwpDir)

		wfwparc := wfwparcio.New(cfg, wfwpSchema, wfwpDB)
		err = wfwparc.Connect()
		if err != nil {
			slog.Error("Cannot initialize storage", "error", err)
			os.Exit(1)
		}

		fc := fwfwp.New(cfg, wfwparc)
		var arc wfwp.Archive

		slog.Info("Importing WF/WP data", "file", wfwpPath)
		arc, err = fc.GetWFWP(wfwpPath)
		if err != nil {
			slog.Error("Cannot get WF/WP Archive", "error", err)
			os.Exit(1)
		}
		_ = arc

		slog.Info("Exporting data to SQLite")
		err = fc.ImportWFWP(arc)
		if err != nil {
			slog.Error("Cannot export data", "error", err)
			os.Exit(1)
		}
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
	rootCmd.Flags().BoolP("debug", "d", false, "set debug mode")
	rootCmd.Flags().StringP("cache-dir", "c", "", "cache directory for temporary files")
	rootCmd.Flags().StringP("wrong-fields-num", "w", "",
		`how to process rows with wrong fields number
     choices: 'stop', 'skip', 'process'
     default: 'stop'`)
	rootCmd.Flags().IntP("jobs-number", "j", 0, "number of concurrent jobs")
	rootCmd.Flags().BoolP("binary-output", "b", false, "return binary SQLite database")
	rootCmd.Flags().BoolP("zip-output", "z", false, "compress output with zip")
	rootCmd.Flags().BoolP("version", "V", false, "shows app's version")
}
