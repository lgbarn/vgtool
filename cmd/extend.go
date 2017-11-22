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
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var file string
var filePtr *string
var targetvgPtr string
var lvextendPtr string

type lvol struct {
	lvPath     string
	lvName     string
	vgName     string
	lvSize     string
	lvSizeUnit string
}
type vg struct {
	vgName     string
	vgSize     float64
	vgSizeUnit string
	disks      []string
}

func (Lvol *lvol) lvExtend() {
	fmt.Printf("lvextend -L %s -n %s\n", Lvol.lvSize, Lvol.lvPath)
	fmt.Printf("fsadm resize %s\n", Lvol.lvPath)
}
func (VG *vg) vgCreate() {
	disk := strings.Join(VG.disks[:], " ")
	fmt.Printf("vgcreate %s %s\n", VG.vgName, disk)
}
func (Lvol *lvol) lvCreate() {
	fmt.Printf("lvcreate -L %s -n %s %s\n", Lvol.lvSize, Lvol.lvName, Lvol.vgName)
}

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
		fmt.Println("extend called")
		fmt.Println(file)

		var fileLine string
		var vgName string
		var lvName string
		var lvPath string
		var lvSize string
		lvols := make([]*lvol, 0, 50)
		pvDisks := make([]string, 0, 50)

		vgNameRE, _ := regexp.Compile(`VG Name\s+(\w+)`)
		lvNameRE, _ := regexp.Compile(`LV Name\s+(\w+)`)
		lvSizeRE, _ := regexp.Compile(`LV Size\s+(.+)`)
		lvPathRE, _ := regexp.Compile(`LV Path\s+(.+)`)
		pvNameRE, _ := regexp.Compile(`PV Name\s+(.+)`)
		filePtr, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		fileScanner := bufio.NewScanner(filePtr)
		for fileScanner.Scan() {
			fileLine = fileScanner.Text()
			if err := fileScanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading file:", err)
			} else {
				switch {
				case vgNameRE.MatchString(fileLine):
					vgName = vgNameRE.FindStringSubmatch(fileLine)[1]
				case lvNameRE.MatchString(fileLine):
					lvName = lvNameRE.FindStringSubmatch(fileLine)[1]
				case lvPathRE.MatchString(fileLine):
					lvPath = lvPathRE.FindStringSubmatch(fileLine)[1]
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
		}

		newvg := &vg{vgName: vgName, disks: pvDisks}
		newvg.vgCreate()
		for _, currLvol := range lvols {
			currLvol.lvCreate()
		}

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