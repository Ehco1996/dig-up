package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	curlFilePath string
	startPage    int
	upUID        int
	favID        int
)

func init() {
	rootCmd.PersistentFlags().StringVar(&curlFilePath, "curl-path", ".curl",
		"保存从浏览器里复制的 curl 内容的文件地址")

	rootCmd.PersistentFlags().IntVar(&upUID, "up-uid", 697166795,
		"喜欢的 up 主的 uid，默认 up 是【徐云流浪中国】")

	rootCmd.PersistentFlags().IntVar(&favID, "favorite-id", 0,
		"收藏夹的 ID")

	rootCmd.PersistentFlags().IntVar(&startPage, "start-page", 1,
		"开始检查的页数")

}

var rootCmd = &cobra.Command{
	Use:               "dig-up",
	Short:             "发现了宝藏 up 主？快来考古吧！",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},

	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(curlFilePath)
		if err != nil {
			return err
		}
		if favID == 0 {
			return fmt.Errorf("收藏夹 ID 不能为空，请加上参数 --favorite-id={你的收藏夹 ID}\n")
		}
		return runTUI(false, string(content))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
