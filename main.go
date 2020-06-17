package main

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

func createOrg(client *gitea.Client, name string, description string) {
	// Get an organisation by name
	org, err := client.GetOrg(name)

	if err != nil {
		if err.Error() == "404 Not Found" {
			org, err = client.CreateOrg(gitea.CreateOrgOption{UserName: name, Visibility: "public"})
			if err != nil {
				//TODO: return values
				return
			}
			fmt.Printf("Organisation %s created.\n", org.UserName)
		}
		return
	}
	fmt.Printf("Organisation %s already exist with ID %d.\n", name, org.ID)
}

func createOrgRepo(client *gitea.Client, name string, description string) {
	repos, err := client.SearchRepos(gitea.SearchRepoOptions{Keyword: name, Private: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("%s\n", repos)
	if len(repos) == 0 {
		repo, err := client.CreateOrgRepo("werkgebieden", gitea.CreateRepoOption{Name: name, Description: description, Private: true})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Repo created on %s\n", repo.Created)
	}
}

func main() {
	// Token example: 3f0bf456ab473c30cdcc67b460989c30f015536c

	// Setup new API connection
	client := gitea.NewClient("http://192.168.2.40:3000", "3f0bf456ab473c30cdcc67b460989c30f015536c")

	createOrg(client, "test", "Alle werkgebieden in business")

	// List organisations repositories (in this case default pagenation options)
	repos, err := client.ListOrgRepos("werkgebieden", gitea.ListOrgReposOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print name of each repo we got in repos
	for _, repo := range repos {
		fmt.Printf("Repo name %s \n", repo.Name)
	}

	createOrgRepo(client, "w00006", "Werkgebieden-00006")

	/*
		// Under authorized repo
		newRepo, err := client.CreateRepo(gitea.CreateRepoOption{Name: "self", Description: "For my-self", Private: true})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s\n", newRepo.Created)
	*/

	// New repo in orginisation (only Owner team)
	/*
		newOrgRepo, err := client.CreateOrgRepo("werkgebieden", gitea.CreateRepoOption{Name: "w00003", Description: "Werkgebieden-00003", Private: true})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s\n", newOrgRepo.Created)
	*/
}
