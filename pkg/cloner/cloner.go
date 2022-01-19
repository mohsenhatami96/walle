package cloner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mohsenhatami96/walle/pkg/dhttp"
)

type Cloner struct {
	url           string
	token         string
	apiURL        string
	sshAuth       bool
	sshPrivateKey string
	username      string
}

type group struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	WebURL string `json:"web_url"`
}

type project struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	WebURL  string `json:"web_url"`
	SSHURL  string `json:"ssh_url_to_repo"`
	HttpURL string `json:"http_url_to_repo"`
}

func New(url, token, username string, sshAuth bool, sshPrivateKeyPath string) *Cloner {
	return &Cloner{
		url:           url,
		token:         token,
		apiURL:        url + "/api/v4",
		sshAuth:       sshAuth,
		sshPrivateKey: sshPrivateKeyPath,
		username:      username,
	}
}

func (cloner *Cloner) CloneAll() {
	allGroups := cloner.getAllGroups()
	printGroups(allGroups)
	userWantedGroups := getUserWantedGroups(allGroups)
	if err := os.MkdirAll("projects", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(filepath.Join(currDir, "projects"))
	if err != nil {
		log.Fatal(err)
	}
	for _, group := range userWantedGroups {
		err = os.MkdirAll(getGroupPath(group), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir(filepath.Join(currDir, "projects", getGroupPath(group)))
		if err != nil {
			log.Fatal(err)
		}
		projects := cloner.getProjectsOfGroup(group)
		// fmt.Println(projects)
		cloner.cloneProjects(projects)
		err = os.Chdir(filepath.Join(currDir, "projects"))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.Chdir(currDir)
	if err != nil {
		log.Fatal(err)
	}
}

func (cloner *Cloner) getAllGroups() []group {
	resp, err := dhttp.Getter(cloner.apiURL+"/groups?per_page=100", &cloner.token)
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

func getGroupPath(grp group) string {
	groupWebURL := grp.WebURL
	splittedGroupWebURL := strings.Split(groupWebURL, "/groups/")
	if len(splittedGroupWebURL) < 2 {
		return ""
	}
	return splittedGroupWebURL[1]
}

func printGroups(groups []group) {
	fmt.Println("Index\tPoject Name\tProject URL")
	for index, group := range groups {
		fmt.Println(fmt.Sprintf("%d.", index+1), group.Name, group.WebURL)
	}
}

func (cloner *Cloner) getProjectsOfGroup(wantedGroup group) []project {
	projects := make([]project, 0)
	path := fmt.Sprintf("/groups/%d/projects?per_page=100", wantedGroup.ID)
	// fmt.Println(path)
	resp, err := dhttp.Getter(cloner.apiURL+path, &cloner.token)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	// fmt.Println("resp:", string(resp))
	projList := make([]project, 0)
	err = json.Unmarshal(resp, &projList)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	projects = append(projects, projList...)
	return projects
}

func (cloner *Cloner) cloneProjects(projects []project) {
	privateKeyFile := "/home/blood/.ssh/id_ed25519"
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		fmt.Errorf("read file %s failed %s\n", privateKeyFile, err.Error())
		return
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if err != nil {
		fmt.Errorf("generate publickeys failed: %s\n", err.Error())
		return
	}
	for _, project := range projects {
		fmt.Println(project.SSHURL)
		_, err := git.PlainClone(project.Name, false, &git.CloneOptions{
			Auth:     publicKeys,
			URL:      project.SSHURL,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
}
