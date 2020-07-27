package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

var (
	GIT_URL   = ""
	GIT_TOKEN = ""
)

var gl *gitlab.Client

func main() {
	client, err := gitlab.NewClient(GIT_TOKEN, gitlab.WithBaseURL(GIT_URL))
	if err != nil {
		log.Fatalf("Can't connect to gitlab: %v", err)
		return
	}
	gl = client
	repo := os.Args[1]
	fmt.Println(LastTag(FindProject(repo)))
}

func FindProject(url string) *gitlab.Project {
	projects, _, err := gl.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		log.Fatalf("Get git lab projects error: %v", err)
		return nil
	}

	groups, _, _ := gl.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	for _, group := range groups {
		groupProjects, _, _ := gl.Groups.ListGroupProjects(group.ID, &gitlab.ListGroupProjectsOptions{})
		projects = append(projects, groupProjects...)
	}

	for _, v := range projects {
		if v.HTTPURLToRepo == url {
			return v
		}
	}

	if err != nil {
		log.Fatalf("Get git lab projects error: %v", err)
		return nil
	}
	return nil
}

func LastTag(project *gitlab.Project) string {
	var lastTag string
	if project == nil {
		log.Fatalf("Can't find the project")
	}
	tags, _, err := gl.Tags.ListTags(project.ID, &gitlab.ListTagsOptions{})
	if err != nil {
		log.Fatalf("Tag error %v", err)
	}

	commitTime := tags[0].Commit.CommittedDate.Unix() - 1
	for _, tag := range tags {
		if tag.Commit.CommittedDate.Unix() > commitTime {
			commitTime = tag.Commit.CommittedDate.Unix()
			lastTag = tag.Name
		}
	}
	return lastTag
}
