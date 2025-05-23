package hotSearch

import (
	"encoding/json"
	"io"
	"net/http"
	"server/model/other"
	"time"
)

type Zhihu struct {
}

func (*Zhihu) GetHotSearchData(maxNum int) (other.HotSearchData, error) {
	resp, err := http.Get("https://orz.ai/api/v1/dailynews/?platform=weibo")
	if err != nil {
		return other.HotSearchData{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return other.HotSearchData{}, err
	}
	var zhihuHotNew other.ZhihuHotNew
	err = json.Unmarshal(body, &zhihuHotNew)
	if err != nil {
		return other.HotSearchData{}, err
	}
	updataTime := time.Now().Format("2006-01-02 15:04:05")
	var hotList []other.HotItem
	for i := 0; i < maxNum; i++ {
		hotList = append(hotList, other.HotItem{
			Index:      i + 1,
			Title:      zhihuHotNew.Data[i].Title,
			Image:      "/uploads/image/96d6f2e7e1f705ab5e59c84a6dc009b2-20250518155549.png",
			Popularity: zhihuHotNew.Data[i].Content[8:],
			URL:        zhihuHotNew.Data[i].URL,
		})
	}
	return other.HotSearchData{Source: "知乎热榜", UpdateTime: updataTime, HotList: hotList}, nil
}

// resp, err := http.Get("https://www.zhihu.com/topsearch")
// if err != nil {
// return other.HotSearchData{}, err
// }
// defer resp.Body.Close()
// body, err := io.ReadAll(resp.Body)
// if err != nil {
// return other.HotSearchData{}, err
// }
//
// var jsonStr string
// reg := regexp.MustCompile(`(?s)<script id="js-initialData" type="text/json">(.*?)</script>`) // 原来是正则表达式库啊regexp
// result := reg.FindAllStringSubmatch(string(body), -1)
// if len(result) > 0 && len(result[0]) > 1 {
// jsonStr = result[0][1]
// } else {
// return other.HotSearchData{}, errors.New("failed to get data")
// }
//
// updateTime := time.Now().Format("2006-01-02 15:04:05")
//
// var hotList []other.HotItem
// for i := 0; i < maxNum; i++ {
// if index := gjson.Get(jsonStr, "initialState.topsearch.data."+strconv.Itoa(i)+".id"); !index.Exists() {
// break
// }
// hotList = append(hotList, other.HotItem{
// Index:       i + 1,
// Title:       gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.titleArea.text").Str,
// Description: gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.excerptArea.text").Str,
// Image:       fmt.Sprintf(gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.imageArea.url").Str),
// Popularity:  gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.metricsArea.text").Str,
// URL:         fmt.Sprintf(gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.link.url").Str),
// })
// }
