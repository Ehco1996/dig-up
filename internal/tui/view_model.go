package tui

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	PageSize = 10

	seenNotCheck = "未知"
)

func (m *model) FetchVideoPage(page, pageSize int) error {

	helpMsg := fmt.Sprintf("n:下一页 | enter:将未观看的视频加入收藏夹(%d) | 按住 p:批量操作本页", m.favID)

	idC := table.Column{Title: "视频 AVID"}
	titleC := table.Column{Title: "标题", Width: len(helpMsg)}
	timeC := table.Column{Title: "播放时间", Width: len("2006-01-02 15:04:05")}
	seenC := table.Column{Title: "是否观看", Width: len("是否观看")}

	rows := []table.Row{}

	videoInfo, err := m.c.GetUPVideoList(context.TODO(), m.upUID, page, pageSize)
	if err != nil {
		return err
	}
	// video rows
	for _, v := range videoInfo.Data.List.Vlist {
		aid := fmt.Sprint(v.Aid)

		if len(aid) > idC.Width {
			idC.Width = len(aid)
		}

		title := v.Title
		if len(title) > titleC.Width {
			titleC.Width = len(title)
		}
		ctime := time.Unix(int64(v.Created), 0).Format("2006-01-02 15:04:05")
		seen := seenNotCheck
		rows = append(rows, []string{aid, title, ctime, seen})
	}

	m.totalPage = videoInfo.Data.Page.Count/pageSize + 1
	m.currentPage = page
	rows = append(rows, []string{})
	rows = append(rows, []string{"按键说明:", helpMsg, fmt.Sprintf("当前页:%d", m.currentPage), fmt.Sprintf("共%d页", m.totalPage)})
	m.tableRows = &rows

	columns := []table.Column{idC, titleC, timeC, seenC}
	width := 0
	for _, c := range columns {
		width += c.Width
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
		table.WithWidth(width+4), // padding
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m.table = t
	return nil
}

func (m *model) CheckAndSave() error {
	ctx := context.TODO()
	rows := *m.tableRows
	row := rows[m.table.Cursor()]
	// reduce http call
	if row[3] != seenNotCheck {
		return nil
	}

	seen, err := m.c.AlreadySeen(ctx, m.upUID, row[1])
	if err != nil {
		return err
	}
	if seen {
		row[3] = "已观看"
	} else {
		row[3] = "未观看"
	}
	m.table.SetRows(rows)
	if !seen {
		aID, _ := strconv.Atoi(row[0]) // must success
		return m.c.AddToFavorite(ctx, aID, m.favID)
	}
	return nil
}
