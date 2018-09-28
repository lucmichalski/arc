# Environment Variables

Plugins might require certain environment variables to be in order initialize the components it 
needs for its functioning. Those variables can be declared in any file. The path to that file
must be provided via the `--env` flag.

**Note:** `ES_CLUSTER_URL` is by all the plugins that are interacting with elasticsearch. `USER_ID` and
 `PASSWORD` are temporary entry point master credentials in order to test the plugins. 

List of specific env vars required by respective plugins are listed below:

##### 1. Users
- `USER_ES_INDEX`
- `USER_ES_TYPE`

##### 2. Permissions
- `PERMISSIONS_ES_INDEX`
- `PERMISSIONS_ES_TYPE`

##### 3. Auth
- `USERS_ES_INDEX`
- `USERS_ES_TYPE`
- `PERMISSIONS_ES_INDEX`  
- `PERMISSIONS_ES_TYPE`