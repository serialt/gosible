/*
Copyright © 2021 windvalley

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
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cimau/gossh/internal/pkg/configflags"
	"github.com/cimau/gossh/internal/pkg/sshtask"
	"github.com/cimau/gossh/pkg/util"
)

var (
	files          []string
	fileDstPath    string
	allowOverwrite bool
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push [HOST...]",
	Short: "Copy local files and dirs to target hosts",
	Long: `
Copy local files and dirs to target hosts.`,
	Example: `
  # Copy a local file or dir to /tmp/ of the target hosts by default.
  $ gossh push host[1-2] -f /path/foo -k

  # Copy local files and dirs to /home/user/ of the target hosts. 
  $ gossh push host[1-2] -f /path/foo.txt,/path/bar/ -d /home/user -k

  Find more examples at: https://github.com/cimau/gossh/blob/main/docs/push.md`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if errs := configflags.Config.Validate(); len(errs) != 0 {
			util.CheckErr(errs)
		}

		if len(files) != 0 {
			for _, f := range files {
				_, err := os.Stat(f)
				if err != nil {
					util.CheckErr(err)
				}
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		task := sshtask.NewTask(sshtask.PushTask, configflags.Config)

		var zipFiles []string

		workDir, err := os.Getwd()
		if err != nil {
			util.CheckErr(err)
		}

		for _, f := range files {
			fileName := filepath.Base(f)
			zipName := "." + fileName + "." + fmt.Sprintf("%d", time.Now().UnixMicro())
			zipFile := path.Join(workDir, zipName)

			if err := util.Zip(strings.TrimSuffix(f, string(os.PathSeparator)), zipFile); err != nil {
				util.CheckErr(err)
			}

			zipFiles = append(zipFiles, zipFile)
		}

		defer func() {
			for _, f := range zipFiles {
				if err := os.Remove(f); err != nil {
					fmt.Printf("Warning: %v\n", err)
				}
			}
		}()

		task.SetTargetHosts(args)
		task.SetPushfiles(files, zipFiles)
		task.SetPushOptions(fileDstPath, allowOverwrite)

		task.Start()

		util.CobraCheckErrWithHelp(cmd, task.CheckErr())
	},
}

func init() {
	pushCmd.Flags().StringSliceVarP(&files, "files", "f", nil,
		"local files/dirs to be copied to target hosts",
	)

	pushCmd.Flags().StringVarP(&fileDstPath, "dest-path", "d", "/tmp",
		"path of target hosts where files/dirs will be copied to",
	)

	pushCmd.Flags().BoolVarP(
		&allowOverwrite,
		"force",
		"F",
		false,
		"allow overwrite files/dirs if they already exist on target hosts",
	)

	pushCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		util.CobraMarkHiddenGlobalFlags(
			command,
			"run.sudo",
			"run.as-user",
			"run.lang",
		)

		command.Parent().HelpFunc()(command, strings)
	})
}
