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
		// TODO: for now only private repos
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

func createTeam(client *gitea.Client, org string, name string) {
	var setTeamOptions gitea.CreateTeamOption

	setTeamOptions.Name = name
	setTeamOptions.Description = "Team for workspace to work with multiple persons on one repository"
	setTeamOptions.Permission = "write"
	setTeamOptions.Units = []string{"repo.code", "repo.issues", "repo.pulls", "repo.releases"}

	team, err := client.CreateTeam(org, setTeamOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(team.Name + " created.")
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

func getTeamID(client *gitea.Client, org string, teamName string) (id int64) {
	// API http team search does not work and gitea sdk doesn't have
	// this option. So some assumptions: orgname == teamname
	teams, err := client.ListOrgTeams(org, gitea.ListTeamsOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, team := range teams {
		if team.Name == teamName {
			return team.ID
		}
	}
	return int64(404)
}

func addTeamRepo(client *gitea.Client, org string, team string, repo string) {
	// Find team ID
	id := getTeamID(client, org, team)
	if id == int64(404) {
		fmt.Println("Team ID not found")
		return
	}
	err := client.AddTeamRepository(id, org, repo)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//
func addTeamMember(client *gitea.Client, org string, team string, user string) {
	id := getTeamID(client, org, team)
	if id == int64(404) {
		fmt.Println("Team ID not found")
		return
	}
	err := client.AddTeamMember(id, user)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func printUsage() {
	fmt.Println("Expected a subcommand:")
	fmt.Println("listrepos\t- list repositories of an organisation")
	fmt.Println("createorg\t- to create an organisation")
	fmt.Println("createorgrepo\t- to create a repository in an organisation")
	fmt.Println("createteam\t- to create a team in an organisation")
	fmt.Println("addteamrepo\t- to add a repository to a team")
	fmt.Println("addteammember\t- to add a member to a team")
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

	// Flag set create a team
	createteam := flag.NewFlagSet("createteam", flag.ExitOnError)
	orgFlag = createteam.String("o", "", "In which organisation to create the team.")
	teamNameFlag := createteam.String("n", "", "The name of the team.")

	// Flag set to add repo to team
	addteamrepo := flag.NewFlagSet("addteamrepo", flag.ExitOnError)
	orgFlag = addteamrepo.String("o", "", "Which organisation contains the team.")
	teamFlag := addteamrepo.String("t", "", "Name of the team")
	repoFlag := addteamrepo.String("r", "", "Name of the repository")

	addteammember := flag.NewFlagSet("addteammember", flag.ExitOnError)
	orgFlag = addteammember.String("o", "", "Which organisation contains the team.")
	teamFlag = addteammember.String("t", "", "Name of the team")
	userFlag := addteammember.String("u", "", "Name of the user to add")

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
	case "createteam":
		createteam.Parse(os.Args[2:])
		if *orgFlag != "" && *teamNameFlag != "" {
			createTeam(client, *orgFlag, *teamNameFlag)
		} else {
			createteam.Usage()
		}
	case "addteamrepo":
		addteamrepo.Parse(os.Args[2:])
		if *orgFlag != "" && *repoFlag != "" && *teamFlag != "" {
			addTeamRepo(client, *orgFlag, *teamFlag, *repoFlag)
		} else {
			addteamrepo.Usage()
		}
	case "addteammember":
		addteammember.Parse(os.Args[2:])
		if *orgFlag != "" && *userFlag != "" && *teamFlag != "" {
			addTeamMember(client, *orgFlag, *teamFlag, *userFlag)
		} else {
			addteammember.Usage()
		}
	default:
		printUsage()
		os.Exit(1)
	}
}
