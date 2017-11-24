// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// extendCmd represents the extend command
var extendCmd = &cobra.Command{
	Use:   "extend",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		lines, err := ParseLines(file, func(s string) (string, bool) {
			return s, true
		})
		if err != nil {
			fmt.Println("Error while parsing file", err)
			return
		}

		_, lvols, _ := matchLines(lines)

		for _, currLvol := range lvols {
			currLvol.lvExtend()
		}

	},
}

func init() {
	RootCmd.AddCommand(extendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//flag := Command.Flags()
	//flag.StringVarP(&file, "file", "", file,
	//	"file")
	//flag.StringVarP(&target, "target", "", target,
	//	"target")

	extendCmd.Flags().StringVarP(&file, "file", "f", "", "Help message for toggle")
	extendCmd.Flags().StringVarP(&targetvgPtr, "target", "t", "", "Help message for toggle")
}
