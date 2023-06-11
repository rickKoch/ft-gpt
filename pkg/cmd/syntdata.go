package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/rickKoch/ftgpt/pkg/openai"
	"github.com/spf13/cobra"
)

var API_TOKEN = os.Getenv("OPENAI_API_KEY")

var (
	templatePath   string
	outputDir      string
	model          string
	maxCompletions int
	maxTokens      int
	temperature    float32
)

type additionalArgs map[string]string

type dataParams struct {
	name   string
	values []string
}

type completionData struct {
	prompt map[string]interface{}
	text   string
}

var GenerateSyntDataCMD = &cobra.Command{
	Use: "generate-synt-data",
	Run: func(_ *cobra.Command, args []string) {
		parsedArgs := parseAdditionalArgs(args)
		params := readFilesFromArgs(parsedArgs)

		paramCombo := []map[string]interface{}{}
		combineParams(params, &paramCombo, nil)

		templateContent := strings.Join(readFile(templatePath), "\n")
		template, err := template.New("template").Parse(templateContent)
		if err != nil {
			panic(err)
		}

		var completions []completionData
		var buff bytes.Buffer
		for _, v := range paramCombo {
			template.Execute(&buff, v)
			completions = append(completions, completionData{
				prompt: v,
				text:   buff.String(),
			})
			buff.Reset()
		}

		wg := sync.WaitGroup{}
		concurrently := maxCompletions
		if maxCompletions > len(completions) {
			concurrently = len(completions)
		}

		wg.Add(concurrently)

		client := openai.OpenAI{}
		for _, completion := range completions {
			concurrently--

			go func(wg *sync.WaitGroup, client *openai.OpenAI, cmp completionData) {
				defer wg.Done()
				res, err := client.CreateCompletion(&openai.CompletionRequest{
					Model:       model,
					Prompt:      cmp.text,
					MaxTokens:   maxTokens,
					Temperature: float32(temperature),
				})
				if err != nil {
					fmt.Printf("error sending http request: %s\n", err)
					panic(err)
				}

				fmt.Println("Completion:", res)
				filename := uuid.New().String()
				prompt := ""
				for k, v := range cmp.prompt {
					prompt += fmt.Sprintf("%s: %s\n", strings.ToUpper(k[:1])+k[1:], v)
				}
				prompt += "\n\nOUTLINE:\n"
				writeFile("prompts", filename, prompt)
				writeFile("completions", filename, res)
			}(&wg, &client, completion)

			if concurrently <= 0 {
				break
			}
		}
		wg.Wait()
	},
}

func combineParams(params []dataParams, results *[]map[string]interface{}, p map[string]interface{}) {
	if len(params) == 0 {
		return
	}

	if p == nil {
		p = make(map[string]interface{})
	}

	for _, value := range params[0].values {
		p[params[0].name] = value
		if len(params) == 1 {
			copyP := make(map[string]interface{})
			for k, v := range p {
				copyP[k] = v
			}
			*results = append(*results, copyP)
		} else {
			combineParams(params[1:], results, p)
		}
	}
}

// parseAdditionalArgs parses the additional arguments passed to the command.
// additional arguments are used for the template params.
func parseAdditionalArgs(args []string) additionalArgs {
	parsedAdditionalArgs := make(additionalArgs)
	for _, arg := range args {
		parsedArg := strings.Split(arg, "=")
		parsedAdditionalArgs[parsedArg[0]] = parsedArg[1]
	}
	return parsedAdditionalArgs
}

// readFilesFromArgs reads the files passed as additional arguments.
// files contains all the template params.
func readFilesFromArgs(args additionalArgs) []dataParams {
	var data []dataParams
	for name, path := range args {
		fmt.Println("reading file: ", path)
		values := readFile(path)
		data = append(data, dataParams{name: name, values: values})
	}

	return data
}

// readFile reads the file from the given path line by line.
// each line is a different param.
func readFile(path string) []string {
	var result []string
	readFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if len(line) == 0 {
			continue
		}
		result = append(result, line)
	}

	return result
}

func writeFile(dirPath, filename, data string) error {
	dirPath = filepath.Join(outputDir, dirPath)
	if dirPath != "" {
		err := os.MkdirAll(dirPath, 0o755)
		if err != nil {
			return err
		}
	}
	path := filepath.Join(dirPath, filename)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(data)
	if err != nil {
		f.Close()
		panic(err)
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func init() {
	if API_TOKEN == "" {
		panic("OPENAI_API_TOKEN environment variable not set")
	}

	GenerateSyntDataCMD.Flags().StringVarP(&templatePath, "template", "t", "", "template file")
	GenerateSyntDataCMD.Flags().StringVarP(&outputDir, "output-dir", "o", "data", "output directory")
	GenerateSyntDataCMD.Flags().StringVarP(&model, "model", "m", "text-davinci-003", "output directory")
	GenerateSyntDataCMD.Flags().IntVarP(&maxCompletions, "max-completions", "c", 10, "max number of completions")
	GenerateSyntDataCMD.Flags().IntVarP(&maxTokens, "max-tokens", "a", 1000, "Max tokens per completion")
	GenerateSyntDataCMD.Flags().Float32VarP(&temperature, "temperature", "e", 1.0, "completion temperature")
	RootCMD.AddCommand(GenerateSyntDataCMD)
}
