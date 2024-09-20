package page

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/plant99/telegraphcl/pkg/user"
	"github.com/plant99/telegraphcl/pkg/util"
)

func ListPages() {
	// get access_token
	accessToken, err := util.FetchAccessToken()
	if err != nil {
		fmt.Println(err)
	}
	// make request parameters
	getPageListRequest := GetPageListRequest{
		AccessToken: accessToken,
		Offset:      0,
		Limit:       200,
	}
	parser := jsoniter.ConfigFastest
	// get page list
	data, err := util.MakeRequest("getPageList", getPageListRequest)
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	pageList := new(PageList)

	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	if err = parser.Unmarshal(data, pageList); err != nil {
		fmt.Println("Couldn't handle api.telegra.ph response !!!.")
		return
	}

	fmt.Println("Index | URL | Title")
	fmt.Println("--------------------")
	for i := 0; i < len(pageList.Pages); i++ {
		fmt.Println(i, ")", pageList.Pages[i].URL, "|", pageList.Pages[i].Title)
	}
}

func GetViews(path string) {
	// make request parameters
	requestPageViews := map[string]string{
		"path": path,
	}

	responsePageViews := new(PageViews)
	// create page
	data, err := util.MakeRequest("getViews", requestPageViews)

	parser := jsoniter.ConfigFastest

	if err = parser.Unmarshal(data, &responsePageViews); err != nil {
		fmt.Println("Couldn't handle api.telegra.ph response. Is the Telegra.ph path correct?")
	}
	fmt.Println(responsePageViews.Views)

}

func CreatePage(path string, title string) {
	// get access_token
	accessToken, err := util.FetchAccessToken()
	if err != nil {
		fmt.Println("Error fetching access token:", err)
		return
	}

	// get []Nodes from markdown file in path
	nodes, err := MarkdownFileToNodes(path)
	if err != nil {
		fmt.Println("Error converting markdown to nodes:", err)
		return
	}

	userInfo := user.GetCurrentUserNameAndURL()
	if len(userInfo) < 2 {
		fmt.Println("Error: User information is incomplete.")
		return
	}

	createPageRequestInstance := createPageRequest{
		Title:         title,
		AuthorName:    userInfo[0],
		AuthorUrl:     userInfo[1],
		Content:       nodes,
		ReturnContent: true,
		AccessToken:   accessToken,
	}

	createPageResponseInstance := new(Page)
	// get total views on page

	data, err := util.MakeRequest("createPage", createPageRequestInstance)
	parser := jsoniter.ConfigFastest
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	if data == nil {
		fmt.Println("Received nil data from MakeRequest")
		return
	}

	if err := parser.Unmarshal(data, &createPageResponseInstance); err != nil {
		fmt.Println("Couldn't handle api.telegra.ph response. Is the Telegra.ph path correct?")
		return
	}

	fmt.Println(createPageResponseInstance.URL)

}

func GetPage(path string) {

	requestGetPage := map[string]string{
		"path": path,
	}

	responseGetPage := new(Page)
	// get total views on page
	data, err := util.MakeRequest("getPage", requestGetPage)
	parser := jsoniter.ConfigFastest
	if err = parser.Unmarshal(data, &responseGetPage); err != nil {
		fmt.Println("Couldn't handle api.telegra.ph response. Is the Telegra.ph path correct?")
	}
	fmt.Println(responseGetPage.Title, responseGetPage.URL)
}

func EditPage(path string, markdownPath string, title string) {
	// get access_token
	accessToken, err := util.FetchAccessToken()
	if err != nil {
		fmt.Println(err)
	}

	// get []Nodes from markdown file in path
	nodes, err := MarkdownFileToNodes(markdownPath)

	if err != nil {
		panic(err)
	}

	userInfo := user.GetCurrentUserNameAndURL()

	editPageRequestInstance := editPageRequest{
		Title:         title,
		AuthorName:    userInfo[0],
		AuthorUrl:     userInfo[1],
		Content:       nodes,
		ReturnContent: true,
		AccessToken:   accessToken,
		Path:          path,
	}

	editPageResponseInstance := new(Page)
	// get total views on page
	data, err := util.MakeRequest("editPage", editPageRequestInstance)
	parser := jsoniter.ConfigFastest
	if err = parser.Unmarshal(data, &editPageResponseInstance); err != nil {
		fmt.Println("Couldn't handle api.telegra.ph response. Is the Telegra.ph path correct?")
	}
	fmt.Println(editPageResponseInstance.URL)

}
