package printer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Ehco1996/dig-up/pkg/bc"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	PageSize = 10
)

func AddNotSeenVideo(ctx context.Context, upUID, favID int, c *bc.Client, res *bc.GetUPVideoListRes) error {
	for _, v := range res.Data.List.Vlist {
		seen, err := c.AlreadySeen(ctx, upUID, v.Title)
		if err != nil {
			return err
		}
		if !seen {
			favErr := c.AddToFavorite(ctx, v.Aid, favID)
			if favErr != nil {
				return err
			}
			fmt.Printf("add title=%s to fav=%d success \n", v.Title, favID)
		}
		time.Sleep(time.Second)
	}
	return nil

}

func PrintUpVideos(curlString string, upUID int) error {
	c, err := bc.NewClient(curlString)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	videoInfo, err := c.GetUPVideoList(ctx, upUID, 1, PageSize)
	if err != nil {
		return err
	}

	totalCount := videoInfo.Data.Page.Count
	// totalPage := totalCount/PageSize + 1

	for i := 0; i < 10; i++ {
		fmt.Print("\033[u\033[K")
		fmt.Println(i)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"视频 AVID", "标题", "播放时间", "已观看"})
		for _, v := range videoInfo.Data.List.Vlist {
			t.AppendRow([]interface{}{v.Aid, v.Title, v.Created, i})
		}

		t.AppendFooter(table.Row{"up", "徐云浏览中国", "视频总量", totalCount})
		t.Render()

		for j := 0; j < 10; j++ {
			fmt.Print("\033[u\033[K")
			fmt.Println("current ", j)
		}
		time.Sleep(time.Second)
	}

	return nil
}
