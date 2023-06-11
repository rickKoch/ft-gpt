package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	prompts     string
	completions string
	jsonDir     string
	jsonFile    string
)

type output struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

var GenerateJsonCMD = &cobra.Command{
	Use: "generate-json",
	Run: func(_ *cobra.Command, _ []string) {
		files, err := os.ReadDir(completions)
		if err != nil {
			panic(err)
		}
		var outputs []output
		for _, file := range files {
			completionData, err := os.ReadFile(completions + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			promptData, err := os.ReadFile(prompts + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			o := output{
				Completion: string(completionData),
				Prompt:     string(promptData),
			}
			outputs = append(outputs, o)
		}

		path := filepath.Join(jsonDir, jsonFile)
		file, _ := os.OpenFile(path, os.O_CREATE, os.ModePerm)
		defer file.Close()

		json, err := json.Marshal(outputs)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(path, json, os.ModePerm)
	},
}

func init() {
	GenerateJsonCMD.Flags().StringVarP(&prompts, "prompts", "p", "", "Path to prompt files")
	GenerateJsonCMD.Flags().StringVarP(&completions, "completions", "c", "", "Path to completion files")
	GenerateJsonCMD.Flags().StringVarP(&jsonDir, "output", "o", "data", "Directory to output json files")
	GenerateJsonCMD.Flags().StringVarP(&jsonFile, "file", "f", "data.json", "Name of json file")
	RootCMD.AddCommand(GenerateJsonCMD)
}
