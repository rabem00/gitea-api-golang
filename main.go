package main

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

func main() {
	// Token: 3f0bf456ab473c30cdcc67b460989c30f015536c

	// Setup new API connection
	client := gitea.NewClient("http://192.168.2.40:3000", "3f0bf456ab473c30cdcc67b460989c30f015536c")

	// Get an organisation by name
	org, err := client.GetOrg("werkgebieden")
	if err != nil {
		fmt.Println(err)
	}
	// Print field of organisation; fullname, description, etc
	fmt.Printf("%s\n", org.FullName)

	// List organisations repositories (in this case default pagenation options)
	repos, err := client.ListOrgRepos("werkgebieden", gitea.ListOrgReposOptions{})
	if err != nil {
		fmt.Println(err)
	}
	// Print name of each repo we got in repos
	for _, repo := range repos {
		if repo.HasPullRequests {
			fmt.Printf("Repo %s has pull-request\n", repo.Name)
		}
	}

}
