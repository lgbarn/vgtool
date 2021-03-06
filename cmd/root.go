// Copyright © 2017 Luther Barnum
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
	"os"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vgFile struct {
	file        string
	targetvgPtr string
}

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

// Creater interface used to create Lvols and VGs
type Creater interface {
	Create()
}

// Extender interface used to  Lextendvols and VGs
type Extender interface {
	Extend()
}

func (Lvol *lvol) Extend() {
	fmt.Printf("lvextend -r -L %s -n %s\n", Lvol.lvSize, Lvol.lvName)
}
func (VG *vg) Create() {
	disk := strings.Join(VG.disks[:], " ")
	fmt.Printf("vgcreate %s %s\n", VG.vgName, disk)
}
func (Lvol *lvol) Create() {
	fmt.Printf("lvcreate -L %s -n %s %s\n", Lvol.lvSize, Lvol.lvName, Lvol.vgName)
}

// ParseLines is used by multiple commands to parse files to lines of data
func ParseLines(filePath string, parse func(string) (string, bool)) ([]string, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var results []string
	for scanner.Scan() {
		if output, add := parse(scanner.Text()); add {
			results = append(results, output)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// MatchLines used to match lines in file and populate variables
func matchLines(lines []string) (string, []lvol, []string) {

	var vgName string
	var lvName string
	var lvPath string
	var lvSize string
	var lvols []lvol
	var pvDisks = make([]string, 0, 50)
	var vgNameRE, _ = regexp.Compile(`VG Name\s+(\w+)`)
	var lvNameRE, _ = regexp.Compile(`LV Name\s+(.+)`)
	var lvSizeRE, _ = regexp.Compile(`LV Size\s+(.+)`)
	var pvNameRE, _ = regexp.Compile(`PV Name\s+(.+)`)

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
			if vgFile.targetvgPtr != "" {
				lvPath = strings.Replace(lvPath, vgName, vgFile.targetvgPtr, -1)
				vgName = vgFile.targetvgPtr
			}
			newlvol := lvol{lvName: lvName, vgName: vgName, lvPath: lvPath, lvSize: cleanlvSize}
			lvols = append(lvols, newlvol)
		case pvNameRE.MatchString(fileLine):
			pvName := pvNameRE.FindStringSubmatch(fileLine)[1]
			cleanDisk := strings.TrimSpace(pvName)
			pvDisks = append(pvDisks, cleanDisk)
		}
	}
	return vgName, lvols, pvDisks
}

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vgtool",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vgtool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vgtool" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".vgtool")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
