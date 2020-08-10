package main

import (
        "encoding/json"
        "flag"
        "fmt"
        "os"

        "code.gitea.io/sdk/gitea"
)

var clientVersion = "1.0.2"
var configFile *os.File

type Giteaconf struct {
        Baseurl string `json:"baseurl"`
        Token   string `json:"token"`
}

func createUser(client *gitea.Client, name string, email string) {
        var setUserOptions gitea.CreateUserOption

        change := new(bool)
        *change = false

        setUserOptions.LoginName = name
        setUserOptions.Username = name
        setUserOptions.Email = email
        setUserOptions.Password = name + "!Q"
        setUserOptions.MustChangePassword = change

        usr, err := client.AdminCreateUser(setUserOptions)
        if err != nil {
                fmt.Println(err)
                return
        }
        fmt.Println(usr.Created)
}

func createOrg(client *gitea.Client, name string, description string) {
        // Get an organisation by name
        org, err := client.GetOrg(name)

        if err != nil {
                if err.Error() == "404 Not Found" {
                        org, err = client.CreateOrg(gitea.CreateOrgOption{UserName: name, Visibility: "private", Description: description})
                        if err != nil {
                                fmt.Println(err)
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
        var setRepoOptions gitea.CreateRepoOption
        setRepoOptions.AutoInit = true
        setRepoOptions.DefaultBranch = "master"
        setRepoOptions.Name = name
        setRepoOptions.Description = description
        setRepoOptions.Private = true

        if len(repos) == 0 {
                repo, err := client.CreateOrgRepo(organisation, setRepoOptions)
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

func branchProtection(client *gitea.Client, owner string, repo string) {
        var setBranchProcOpt gitea.CreateBranchProtectionOption
        setBranchProcOpt.BranchName = "master"
        //setBranchProcOpt.EnablePush = true

        _, err := client.CreateBranchProtection(owner, repo, setBranchProcOpt)
        if err != nil {
                fmt.Println(err)
                return
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
                fmt.Printf("%s \n", repo.Name)
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

func removeTeamRepo(client *gitea.Client, org string, team string, repo string) {
        id := getTeamID(client, org, team)
        if id == int64(404) {
                fmt.Println("Team ID not found")
                return
        }
        err := client.RemoveTeamRepository(id, org, repo)
        if err != nil {
                fmt.Println(err)
                return
        }
}

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

func removeTeamMember(client *gitea.Client, org string, team string, user string) {
        id := getTeamID(client, org, team)
        if id == int64(404) {
                fmt.Println("Team ID not found")
                return
        }
        err := client.RemoveTeamMember(id, user)
        if err != nil {
                fmt.Println(err)
                return
        }
}

func createUserPub(client *gitea.Client, user string, title string, pubkey string) {
        var setUserOptions gitea.CreateKeyOption

        setUserOptions.ReadOnly = false
        setUserOptions.Title = title
        setUserOptions.Key = pubkey

        _, err := client.AdminCreateUserPublicKey(user, setUserOptions)
        if err != nil {
                fmt.Println(err)
                return
        }
}

func getVersion(client *gitea.Client) (serverVersion string) {
        serverVersion, err := client.ServerVersion()
        if err != nil {
                fmt.Println(err)
                return
        }
        return serverVersion
}

func printUsage() {
        fmt.Println("Expected a subcommand:")
        fmt.Println("version\t\t\t- list version information of client and server")
        fmt.Println("listrepos\t\t- list repositories of an organisation")
        fmt.Println("createuser\t\t- to create an user")
        fmt.Println("createorg\t\t- to create an organisation")
        fmt.Println("createorgrepo\t\t- to create a repository in an organisation")
        fmt.Println("createteam\t\t- to create a team in an organisation")
        fmt.Println("createuserpub\t\t- to add a public key to an user")
        fmt.Println("branchprotection\t- to add branch protection for a repo")
        fmt.Println("addteamrepo\t\t- to add a repository to a team")
        fmt.Println("addteammember\t\t- to add a member to a team")
        fmt.Println("removeteamrepo\t\t- to remove a repository from a team")
        fmt.Println("removeteammember\t- to remove a member from a team")
}

func main() {
        if _, err := os.Stat("/etc/gitea/config.json"); !os.IsNotExist(err) {
                configFile, err = os.Open("/etc/gitea/config.json")
                if err != nil {
                        fmt.Println(err)
                }
        } else if _, err := os.Stat("config.json"); !os.IsNotExist(err) {
                configFile, err = os.Open("config.json")
                if err != nil {
                        fmt.Println(err)
                }
        } else {
                fmt.Println("Config file not found!")
                os.Exit(1)
        }
        // Defer the closing of the config file so that we can parse it later on
        defer configFile.Close()

        var giteaconf Giteaconf

        decoder := json.NewDecoder(configFile)
        err := decoder.Decode(&giteaconf)
        if err != nil {
                fmt.Println(err)
        }

        // Setup new API connection
        client := gitea.NewClient(giteaconf.Baseurl, giteaconf.Token)

        version := flag.NewFlagSet("version", flag.ExitOnError)

        listrepos := flag.NewFlagSet("listrepos", flag.ExitOnError)
        listReposFlag := listrepos.String("o", "", "Which organisation to list the repos of.")

        // Flag set create an user
        createuser := flag.NewFlagSet("createuser", flag.ExitOnError)
        createUserFlag := createuser.String("l", "", "Login name of the user")
        createUserMailFlag := createuser.String("m", "", "Email name of the user")

        // Flag set create an organisation
        createorg := flag.NewFlagSet("createorg", flag.ExitOnError)
        createOrgFlag := createorg.String("o", "", "Create a workspace (eq organisation).")
        createOrgDescFlag := createorg.String("d", "", "Workspace (eq organisation) description.")

        // Flag set create organisation repository
        createorgrepo := flag.NewFlagSet("createorgrepo", flag.ExitOnError)
        nameFlag := createorgrepo.String("n", "", "Repositoryname.")
        descFlag := createorgrepo.String("d", "", "Repository description.")
        orFlag := createorgrepo.String("o", "", "In which organisation to create the repo.")

        // Flag set create a team
        createteam := flag.NewFlagSet("createteam", flag.ExitOnError)
        orgFlag := createteam.String("o", "", "In which organisation to create the team.")
        teamNameFlag := createteam.String("n", "", "The name of the team.")

        // Flag set to add repo to team
        addteamrepo := flag.NewFlagSet("addteamrepo", flag.ExitOnError)
        orgTeamFlag := addteamrepo.String("o", "", "Which organisation contains the team.")
        nameTeamFlag := addteamrepo.String("n", "", "Name of the team")
        repoTeamFlag := addteamrepo.String("r", "", "Name of the repository to add")

        // Flag set to add repo to team
        removeteamrepo := flag.NewFlagSet("removeteamrepo", flag.ExitOnError)
        rmOrgTeamFlag := removeteamrepo.String("o", "", "Which organisation contains the team.")
        rmNameTeamFlag := removeteamrepo.String("n", "", "Name of the team")
        rmRepoTeamFlag := removeteamrepo.String("r", "", "Name of the repository to remove")

        addteammember := flag.NewFlagSet("addteammember", flag.ExitOnError)
        orgMemFlag := addteammember.String("o", "", "Which organisation contains the team.")
        teamFlag := addteammember.String("t", "", "Name of the team")
        userMemFlag := addteammember.String("u", "", "Name of the user to add")

        removeteammember := flag.NewFlagSet("removeteammember", flag.ExitOnError)
        rmOrgMemFlag := removeteammember.String("o", "", "Which organisation contains the team.")
        rmTeamFlag := removeteammember.String("t", "", "Name of the team")
        rmUserMemFlag := removeteammember.String("u", "", "Name of the user to remove")

        createuserpub := flag.NewFlagSet("createuserpub", flag.ExitOnError)
        userFlag := createuserpub.String("u", "", "Name of the user")
        titleFlag := createuserpub.String("i", "", "Title of the key to add")
        pubkeyFlag := createuserpub.String("p", "", "The public key to add")

        branchprotection := flag.NewFlagSet("branchprotection", flag.ExitOnError)
        ownerFlag := branchprotection.String("m", "", "Name of the owner (usually the team")
        repoFlag := branchprotection.String("r", "", "Name of the repository")

        // A subcommand is needed
        if len(os.Args) < 2 {
                printUsage()
                os.Exit(1)
        }

        switch os.Args[1] {
        case "version":
                version.Parse(os.Args[2:])
                fmt.Println("Client version: " + clientVersion)
                fmt.Println("Server version: " + getVersion(client))
        case "listrepos":
                listrepos.Parse(os.Args[2:])
                if *listReposFlag != "" {
                        listAllReposOrg(client, *listReposFlag)
                } else {
                        listrepos.Usage()
                }
        case "createuser":
                createuser.Parse(os.Args[2:])
                if *createUserFlag != "" && *createUserMailFlag != "" {
                        createUser(client, *createUserFlag, *createUserMailFlag)

                } else {
                        createuser.Usage()
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
                if *nameFlag != "" && *descFlag != "" && *orFlag != "" {
                        createOrgRepo(client, *nameFlag, *descFlag, *orFlag)
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
                if *orgTeamFlag != "" && *repoTeamFlag != "" && *nameTeamFlag != "" {
                        addTeamRepo(client, *orgTeamFlag, *nameTeamFlag, *repoTeamFlag)
                } else {
                        addteamrepo.Usage()
                }
        case "removeteamrepo":
                removeteamrepo.Parse(os.Args[2:])
                if *rmOrgTeamFlag != "" && *rmRepoTeamFlag != "" && *rmNameTeamFlag != "" {
                        removeTeamRepo(client, *rmOrgTeamFlag, *rmNameTeamFlag, *rmRepoTeamFlag)
                } else {
                        removeteamrepo.Usage()
                }
        case "addteammember":
                addteammember.Parse(os.Args[2:])
                if *orgMemFlag != "" && *userMemFlag != "" && *teamFlag != "" {
                        addTeamMember(client, *orgMemFlag, *teamFlag, *userMemFlag)
                } else {
                        addteammember.Usage()
                }
        case "removeteammember":
                removeteammember.Parse(os.Args[2:])
                if *rmOrgMemFlag != "" && *rmUserMemFlag != "" && *rmTeamFlag != "" {
                        removeTeamMember(client, *rmOrgMemFlag, *rmTeamFlag, *rmUserMemFlag)
                } else {
                        removeteammember.Usage()
                }
        case "createuserpub":
                createuserpub.Parse(os.Args[2:])
                if *userFlag != "" && *titleFlag != "" && *pubkeyFlag != "" {
                        createUserPub(client, *userFlag, *titleFlag, *pubkeyFlag)
                } else {
                        createuserpub.Usage()
                }
	case "branchprotection":
		branchprotection.Parse(os.Args[2:])
                if *ownerFlag != "" && *repoFlag != "" {
                        branchProtection(client, *ownerFlag, *repoFlag)
                } else {
                        branchprotection.Usage()
                }
        default:
                printUsage()
                os.Exit(1)
        }
}
