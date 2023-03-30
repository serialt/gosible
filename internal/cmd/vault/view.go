/*
Copyright © 2022 windvalley

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

package vault

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/serialt/gosible/pkg/util"
)

// viewCmd represents the vault view command
var viewCmd = &cobra.Command{
	Use:   "view FILENAME",
	Short: "View vault encrypted file",
	Long: `
View vault encrypted file.`,
	Example: `
  # View a vault encrypted file by asking for vault password.
  $ gossh vault view /path/auth.txt

  # View a vault encrypted file by vault password file or script.
  $ gossh vault view /path/auth.txt -V /path/vault-password-file-or-script`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			util.CobraCheckErrWithHelp(cmd, "requires one arg to represent the vault encrypted file")
		}

		if len(args) > 1 {
			util.CobraCheckErrWithHelp(cmd, "to many args, only need one")
		}

		if !util.FileExists(args[0]) {
			util.CheckErr(fmt.Sprintf("file '%s' not found", args[0]))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		vaultPass := GetVaultPassword()

		file := args[0]

		decryptContent, err := decryptFile(file, vaultPass)
		util.CheckErr(err)

		err = util.LessContent(decryptContent)
		util.CheckErr(err)
	},
}
