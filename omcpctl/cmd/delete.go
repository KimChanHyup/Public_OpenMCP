/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"openmcp/openmcp/omcpctl/apiServerMethod"
	"os"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	//"os"
	"path/filepath"
	"strings"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

openmcpctl delete -f <FILENAME>
openmcpctl delete -f <FILENAME> --context <CLUSTERNAME>

openmcpctl delete -f <FILEDIRECTORY>/<FILENAME>
openmcpctl delete -f <FILEDIRECTORY>/<FILENAME> --context <CLUSTERNAME>

openmcpctl delete <RESOURCE> <NAME>
openmcpctl delete <RESOURCE> <NAME> --context <CLUSTERNAME>`,

	Run: func(cmd *cobra.Command, args []string) {
		Delete(args)
		/*cmdStr := "kubectl delete"

		for i := 0; i < len(args); i++ {
			cmdStr = cmdStr + " " + args[i]
		}
		if cobrautil.Option_file != "" {
			cmdStr = cmdStr + " -f " + cobrautil.Option_file
		}
		if cobrautil.Option_namespace != "" {
			cmdStr = cmdStr + " -f " + cobrautil.Option_namespace
		}
		if cobrautil.Option_context != ""{
			cmdStr = cmdStr + " --context " + cobrautil.Option_context
		}
		util.CmdExec2(cmdStr)*/

	},
}

func Delete(args []string){

	if len(args) >= 2 {
		var metainfo cobrautil.MetaInfo

		resourceKind := args[0]
		resourceName := args[1]

		fmt.Println(resourceKind," / ",resourceName)

		metainfo.Kind = resourceKind
		metainfo.Metadata.Name = resourceName

		if cobrautil.Option_namespace == "" {
			metainfo.Metadata.Namespace = ""
		}else {
			metainfo.Metadata.Namespace = cobrautil.Option_namespace
		}

		fmt.Printf("Value: %#v\n", metainfo.Kind)
		fmt.Printf("Value: %#v\n", metainfo.Metadata.Name)
		fmt.Printf("Value: %#v\n", metainfo.Metadata.Namespace)

		SendToAPIServer(metainfo, nil, "resource")
	}else {
		fileOrDirname, _ := filepath.Abs(cobrautil.Option_file)
		filenameList := []string{}

		fi, err := os.Stat(fileOrDirname)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			// do directory stuff
			fmt.Println("directory")
			files, err := ioutil.ReadDir(fileOrDirname)

			fmt.Println(fileOrDirname)
			fmt.Println(files)

			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				if err != nil {
					fmt.Println(err)
				}
				if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml"{
					filenameList = append(filenameList, fileOrDirname+"/"+f.Name())
				}
			}
		case mode.IsRegular():
			// do file stuff
			fmt.Println("file")
			filenameList = append(filenameList, fileOrDirname)
		}


		for _, filename := range filenameList {
			var metainfo cobrautil.MetaInfo

			yamlFile, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}

			err = yaml.Unmarshal(yamlFile, &metainfo)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Value: %#v\n", metainfo.Kind)
			fmt.Printf("Value: %#v\n", metainfo.Metadata.Name)
			fmt.Printf("Value: %#v\n", metainfo.Metadata.Namespace)

			body := strings.NewReader(string(yamlFile))

			SendToAPIServer(metainfo, body, "kind")
		}
	}

}

func SendToAPIServer(metainfo cobrautil.MetaInfo, body io.Reader, metainfoKindType string){
	LINK := cobrautil.DeleteLinkParser(&metainfo, metainfoKindType)
	fmt.Println(LINK)

	msg, err := apiServerMethod.DeleteAPIServer(LINK, body)
	if err != nil {
		return
	}

	metainfo2, err := getMetaInfo(msg)

	if err != nil {
		return
	}
	if metainfo2.Message != "" {
		fmt.Println(metainfo2.Message)
	} else {
		fmt.Println(cobrautil.KindMap[metainfo.Kind] + " \""+metainfo.Metadata.Name+"\" deleted")
	}


}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")
	deleteCmd.Flags().StringVarP(&cobrautil.Option_file, "file","f", "", "input a option")
	deleteCmd.Flags().StringVarP(&cobrautil.Option_context, "context","c", "", "input a option")
	deleteCmd.Flags().StringVarP(&cobrautil.Option_namespace, "namespace","n", "", "input a option")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
