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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var vgName string
		var lvName string
		var lvPath string
		var lvSize string
		var lvols = make([]*lvol, 0, 50)
		var pvDisks = make([]string, 0, 50)
		var vgNameRE, _ = regexp.Compile(`VG Name\s+(\w+)`)
		var lvNameRE, _ = regexp.Compile(`LV Name\s+(.+)`)
		var lvSizeRE, _ = regexp.Compile(`LV Size\s+(.+)`)
		var pvNameRE, _ = regexp.Compile(`PV Name\s+(.+)`)

		lines, err := ParseLines(file, func(s string) (string, bool) {
			return s, true
		})
		if err != nil {
			fmt.Println("Error while parsing file", err)
			return
		}

		for _, fileLine := range lines {
			//fmt.Println(fileLine)
			switch {
			case vgNameRE.MatchString(fileLine):
				vgName = vgNameRE.FindStringSubmatch(fileLine)[1]
			case lvNameRE.MatchString(fileLine):
				lvName = lvNameRE.FindStringSubmatch(fileLine)[1]
			case lvSizeRE.MatchString(fileLine):
				lvSize = lvSizeRE.FindStringSubmatch(fileLine)[1]
				cleanlvSize := strings.Replace(lvSize, " ", "", -1)
				if targetvgPtr != "" {
					lvPath = strings.Replace(lvPath, vgName, targetvgPtr, -1)
					vgName = targetvgPtr
				}
				newlvol := &lvol{lvName: lvName, vgName: vgName, lvPath: lvPath, lvSize: cleanlvSize}
				lvols = append(lvols, newlvol)
			case pvNameRE.MatchString(fileLine):
				pvName := pvNameRE.FindStringSubmatch(fileLine)[1]
				cleanDisk := strings.TrimSpace(pvName)
				pvDisks = append(pvDisks, cleanDisk)
			}
		}

		newvg := &vg{vgName: vgName, disks: pvDisks}
		newvg.vgCreate()
		for _, currLvol := range lvols {
			currLvol.lvCreate()
		}

		//for _, currLvol := range lvols {
		//	currLvol.lvExtend()
		//}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&file, "file", "f", "", "Help message for toggle")
	createCmd.Flags().StringVarP(&targetvgPtr, "target", "t", "", "Help message for toggle")
}
