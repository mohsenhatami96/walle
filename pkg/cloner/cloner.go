package cloner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mohsenhatami96/dobby/pkg/dhttp"
)

type Cloner struct {
	url    string
	token  string
	apiURL string
}

type group struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	WebURL string `json:"web_url"`
}

func New(url string, token string) *Cloner {
	return &Cloner{url: url, token: token, apiURL: url + "/api/v4"}
}

func (cloner *Cloner) CloneAll() {
	allGroups := cloner.getAllGroups()
	printGroups(allGroups)
	userWantedGroups := getUserWantedGroups(allGroups)
	fmt.Println("\n\nYou wanted these groups:")
	for _, g := range userWantedGroups {
		fmt.Println(g.Name)
	}
}

func (cloner *Cloner) getAllGroups() []group {
	resp, err := dhttp.Getter(cloner.apiURL+"/groups", &cloner.token)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	groups := make([]group, 0)
	err = json.Unmarshal(resp, &groups)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	return groups
}

func getUserAnswerList() []int {
	fmt.Println("\n\nWhich group do you want to clone all the projects of? (comma seperated) [default: all]:")
	reader := bufio.NewReader(os.Stdin)
	groupList, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	trimmedGroupList := strings.TrimSpace(groupList)
	if trimmedGroupList == "" {
		return make([]int, 0)
	}
	// fmt.Println(groupList)
	groups := strings.Split(trimmedGroupList, ",")
	// fmt.Print(groups)
	groupsNumList := make([]int, 0)
	for _, g := range groups {
		gint, err := strconv.ParseInt(strings.TrimSpace(g), 10, 32)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		groupsNumList = append(groupsNumList, int(gint))
	}
	// fmt.Println(groupsNumList)
	return groupsNumList
}

func getUserWantedGroups(allGroups []group) []group {
	userWantedGroups := make([]group, 0)
	userAnswerList := getUserAnswerList()
	// fmt.Println(userAnswerList)
	if len(userAnswerList) == 0 {
		return allGroups
	}
	for _, index := range userAnswerList {
		// fmt.Println(index)
		userWantedGroups = append(userWantedGroups, allGroups[index-1])
	}
	return userWantedGroups
}

func printGroups(groups []group) {
	fmt.Println("Index\tPoject Name\tProject URL")
	for index, group := range groups {
		fmt.Println(fmt.Sprintf("%d.", index+1), group.Name, group.WebURL)
	}
}
