package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:   "ftgpt",
	Short: "Fine tune GPT",
	Long: `
 _______  _______  _______  _______  _______
|       ||       ||       ||       ||       |
|    ___||_     _||    ___||    _  ||_     _|
|   |___   |   |  |   | __ |   |_| |  |   |
|    ___|  |   |  |   ||  ||    ___|  |   |
|   |      |   |  |   |_| ||   |      |   |
|___|      |___|  |_______||___|      |___|  `,
}

func Execute(version string) {
	RootCMD.Version = version

	if err := RootCMD.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
}
