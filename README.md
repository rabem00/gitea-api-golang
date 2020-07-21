## Gitea
Using gitea-sdk go packages to create content in gitea.
Only the some API endpoints are are used for a specific project.

```bash
Usage:
% ./gitea-api-golang                                           
Expected a subcommand:
listrepos       - list repositories of an organisation
createorg       - to create an organisation
createorgrepo   - to create a repository in an organisation
createteam      - to create a team in an organisation
createuserpub   - to add a public key to an user
createbranchprotection  - to add branch protection for a repo
addteamrepo     - to add a repository to a team
addteammember   - to add a member to a team

% ./gitea-api-golang createorg
Usage of createorg:
  -d string
        Workspace (eq organisation) description.
  -o string
        Create a workspace (eq organisation).


% ./gitea-api-golang createorg -o myorg -d "My organisation"
```
