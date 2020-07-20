package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"code.gitea.io/sdk/gitea"
)

type Giteaconf struct {
	Baseurl string `json:"baseurl"`
	Token   string `json:"token"`
}

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

func printUsage() {
	fmt.Println("Expected a subcommand:")
	fmt.Println("listrepos\t- list repositories of an organisation")
	fmt.Println("createorg\t- to create an organisation")
	fmt.Println("createorgrepo\t- to create a repository in an organisation")
}

func main() {
	// Open the config file
	configFile, err := os.Open("api-config.json")
	if err != nil {
		fmt.Println(err)
	}
	// Defer the closing of the config file so that we can parse it later on
	defer configFile.Close()

	var giteaconf Giteaconf

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&giteaconf)
	if err != nil {
		fmt.Println(err)
	}

	// Setup new API connection
	client := gitea.NewClient(giteaconf.Baseurl, giteaconf.Token)

	listrepos := flag.NewFlagSet("listrepos", flag.ExitOnError)
	listReposFlag := listrepos.String("o", "", "Which organisation to list the repos of.")

	// Flag set create an organisation
	createorg := flag.NewFlagSet("createorg", flag.ExitOnError)
	createOrgFlag := createorg.String("co", "", "Create a workspace (eq organisation).")
	createOrgDescFlag := createorg.String("cod", "", "Workspace (eq organisation) description.")

	// Flag set create organisation repository
	createorgrepo := flag.NewFlagSet("createorgrepo", flag.ExitOnError)
	createOrgRepoFlag := createorgrepo.String("r", "", "Create a repository in a workspace (eq organisation).")
	createOrgRepoDescFlag := createorgrepo.String("d", "", "Repository description.")
	orgFlag := createorgrepo.String("o", "", "In which organisation to create the repo.")

	// A subcommand is needed
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "listrepos":
		listrepos.Parse(os.Args[2:])
		if *listReposFlag != "" {
			listAllReposOrg(client, *listReposFlag)
		} else {
			listrepos.Usage()
		}
	// Create an organisation
	case "createorg":
		createorg.Parse(os.Args[2:])
		if *createOrgDescFlag != "" && *createOrgFlag != "" {
			createOrg(client, *createOrgFlag, *createOrgDescFlag)
		} else {
			createorg.Usage()
		}
	// Create a repository in an organisation
	case "createorgrepo":
		createorgrepo.Parse(os.Args[2:])
		if *createOrgRepoFlag != "" && *createOrgRepoDescFlag != "" && *orgFlag != "" {
			createOrgRepo(client, *createOrgRepoFlag, *createOrgRepoDescFlag, *orgFlag)
		} else {
			createorgrepo.Usage()
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

/*
	//user := createUser(client, "test01", "m.rabelink")
	//fmt.Println("%s\n", user.Created)

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
*/
