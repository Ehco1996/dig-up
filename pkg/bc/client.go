package bc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BaseRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

type GetUPVideoListRes struct {
	BaseRes

	Data struct {
		List struct {
			Tlist struct {
				Num160 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"160"`
				Num234 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"234"`
			} `json:"tlist"`
			Vlist []struct {
				Comment        int    `json:"comment"`
				Typeid         int    `json:"typeid"`
				Play           int    `json:"play"`
				Pic            string `json:"pic"`
				Subtitle       string `json:"subtitle"`
				Description    string `json:"description"`
				Copyright      string `json:"copyright"`
				Title          string `json:"title"`
				Review         int    `json:"review"`
				Author         string `json:"author"`
				Mid            int    `json:"mid"`
				Created        int    `json:"created"`
				Length         string `json:"length"`
				VideoReview    int    `json:"video_review"`
				Aid            int    `json:"aid"`
				Bvid           string `json:"bvid"`
				HideClick      bool   `json:"hide_click"`
				IsPay          int    `json:"is_pay"`
				IsUnionVideo   int    `json:"is_union_video"`
				IsSteinsGate   int    `json:"is_steins_gate"`
				IsLivePlayback int    `json:"is_live_playback"`
			} `json:"vlist"`
		} `json:"list"`
		Page struct {
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
			Count int `json:"count"`
		} `json:"page"`
		IsRisk      bool        `json:"is_risk"`
		GaiaResType int         `json:"gaia_res_type"`
		GaiaData    interface{} `json:"gaia_data"`
	} `json:"data"`
}

type AlreadySeenRes struct {
	BaseRes

	Data struct {
		HasMore bool `json:"has_more"`
		Page    struct {
			Pn    int `json:"pn"`
			Total int `json:"total"`
		} `json:"page"`
		List []struct {
			Title     string      `json:"title"`
			LongTitle string      `json:"long_title"`
			Cover     string      `json:"cover"`
			Covers    interface{} `json:"covers"`
			URI       string      `json:"uri"`
			History   struct {
				Oid      int    `json:"oid"`
				Epid     int    `json:"epid"`
				Bvid     string `json:"bvid"`
				Page     int    `json:"page"`
				Cid      int    `json:"cid"`
				Part     string `json:"part"`
				Business string `json:"business"`
				Dt       int    `json:"dt"`
			} `json:"history"`
			Videos     int    `json:"videos"`
			AuthorName string `json:"author_name"`
			AuthorFace string `json:"author_face"`
			AuthorMid  int    `json:"author_mid"`
			ViewAt     int    `json:"view_at"`
			Progress   int    `json:"progress"`
			Badge      string `json:"badge"`
			ShowTitle  string `json:"show_title"`
			Duration   int    `json:"duration"`
			Total      int    `json:"total"`
			NewDesc    string `json:"new_desc"`
			IsFinish   int    `json:"is_finish"`
			IsFav      int    `json:"is_fav"`
			Kid        int    `json:"kid"`
			TagName    string `json:"tag_name"`
			LiveStatus int    `json:"live_status"`
		} `json:"list"`
	} `json:"data"`
}

type Client struct {
	inner *http.Client

	credential []*http.Cookie

	csrfToken string
}

func NewClient(curl string) (*Client, error) {
	cookies := parseCurl2Credential(curl)
	if len(cookies) != 4 {
		return nil, fmt.Errorf("can't parse cookies from curl %v", cookies)
	}

	token := ""
	for _, c := range cookies {
		if c.Name == KEY_bili_jct {
			token = c.Value
		}
	}
	return &Client{inner: http.DefaultClient, credential: cookies, csrfToken: token}, nil
}

func (c *Client) sendRequest(req *http.Request) ([]byte, error) {
	// add default header
	req.Header.Add("Referer", "https://www.bilibili.com")
	req.Header.Add("User-Agent", "Mozilla/5.0")

	// add auth cookie
	for _, c := range c.credential {
		req.AddCookie(c)
	}
	resp, err := c.inner.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := &BaseRes{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	if res.Code != 0 {
		return nil, fmt.Errorf(res.Message)
	}
	return body, nil
}

func (c *Client) GetUPVideoList(ctx context.Context, upUID, pageNumber, pageSize int) (*GetUPVideoListRes, error) {

	url := "https://api.bilibili.com/x/space/arc/search"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	params := map[string]any{
		"mid":     upUID,
		"ps":      pageSize,
		"tid":     0,
		"pn":      pageNumber,
		"keyword": "",
		"order":   "pubdate",
	}

	query := req.URL.Query()

	for k, v := range params {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = query.Encode()

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}
	res := &GetUPVideoListRes{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) AlreadySeen(ctx context.Context, upUID int, title string) (bool, error) {
	url := "https://api.bilibili.com/x/web-goblin/history/search"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	params := map[string]any{
		"pn":       1,
		"keyword":  title,
		"business": "all",
	}
	query := req.URL.Query()
	for k, v := range params {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = query.Encode()

	body, err := c.sendRequest(req)
	if err != nil {
		return false, err
	}

	res := &AlreadySeenRes{}
	if err := json.Unmarshal(body, res); err != nil {
		return false, err
	}
	for _, v := range res.Data.List {
		if v.AuthorMid == upUID {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) AddToFavorite(ctx context.Context, videoAID, favID int) error {
	url := "https://api.bilibili.com/x/v3/fav/resource/deal"
	reqData := fmt.Sprintf("rid=%d&type=2&add_media_ids=%d&del_media_ids=&jsonp=jsonp&csrf=%s&platform=web", videoAID, favID, c.csrfToken)
	var data = strings.NewReader(reqData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, data)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	_, err = c.sendRequest(req)
	return err
}
