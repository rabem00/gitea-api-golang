package main

import (
	"flag"
	"fmt"

	"code.gitea.io/sdk/gitea"
)

var listReposFlag = flag.String("lr", "", "List all repositories for certain workspace (eq organisation)")
var createOrgFlag = flag.String("co", "", "Create a workspace (eq organisation). Needs flag -cod")
var createOrgDescFlag = flag.String("cod", "", "Workspace (eq organisation) description. Needs flag -co")
var createOrgRepoFlag = flag.String("cr", "", "Create a repository in workspace (eq organisation). Needs flags -crd/-org")
var createOrgRepoDescFlag = flag.String("crd", "", "Repository description. Needs flag -cr/-org")
var orgFlag = flag.String("org", "", "Organisation to use. Needs flag -cr/-crd")

func createOrg(client *gitea.Client, name string, description string) {
	// Get an organisation by name
	org, err := client.GetOrg(name)

	if err != nil {
		if err.Error() == "404 Not Found" {
			org, err = client.CreateOrg(gitea.CreateOrgOption{UserName: name, Visibility: "private"})
			if err != nil {
				// TODO: return values
				return
			}
			fmt.Printf("Organisation %s created.\n", org.UserName)
		}
		return
	}
	fmt.Printf("Organisation %s already exist with ID %d.\n", name, org.ID)
}

func createOrgRepo(client *gitea.Client, name string, description string, organisation string) {
	repos, err := client.SearchRepos(gitea.SearchRepoOptions{Keyword: name, Private: true})
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(repos) == 0 {
		// TODO: hardcoded
		repo, err := client.CreateOrgRepo(organisation, gitea.CreateRepoOption{Name: name, Description: description, Private: true})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Repo created on %s\n", repo.Created)
	} else {
		fmt.Printf("Repo %s already exist\n", name)
	}
}

func listAllReposOrg(client *gitea.Client, name string) {
	// List organisations repositories (in this case default pagenation options)
	repos, err := client.ListOrgRepos(name, gitea.ListOrgReposOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(repos) == 0 {
		fmt.Println("No repositories found!")
		return
	}
	// Print name of each repo we got in repos
	for _, repo := range repos {
		fmt.Printf("Repo name %s \n", repo.Name)
	}
}

// TODO: should be LDAP
func createUser(client *gitea.Client, username string, emailname string) *gitea.User {
	bFalse := false
	user, _ := client.GetUserInfo(username)
	if user.ID != 0 {
		return user
	}
	user, err := client.AdminCreateUser(gitea.CreateUserOption{Username: username, Password: username + "!Q", Email: emailname + "@belastingdienst.nl", MustChangePassword: &bFalse, SendNotify: bFalse})
	if err != nil {
		fmt.Println(err)
	}
	return user
}

func main() {
	// Token example: 6787577dd7665afeb801d653935a101f962d9da1

	// Setup new API connection
	client := gitea.NewClient("http://192.168.2.41:3000", "6787577dd7665afeb801d653935a101f962d9da1")

	flag.Parse()

	if *listReposFlag != "" {
		listAllReposOrg(client, *listReposFlag)
	} else {
		flag.Usage()
	}

	if *createOrgDescFlag != "" && *createOrgFlag != "" {
		createOrg(client, *createOrgFlag, *createOrgDescFlag)
	} else {
		flag.Usage()
	}

	if *createOrgRepoFlag != "" && *createOrgRepoDescFlag != "" && *orgFlag != "" {
		createOrgRepo(client, *createOrgRepoFlag, *createOrgRepoDescFlag, *orgFlag)
	} else {
		flag.Usage()
	}

	//user := createUser(client, "test01", "m.rabelink")
	//fmt.Println("%s\n", user.Created)
}
