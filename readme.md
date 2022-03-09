# SSH-Check

Because sometimes you just want to know if your logins work on a number of servers.

## Usage

```bash
Usage of ssh-check.exe:
  -d    activate debug log
  -p string
        password to use for authentication
  -s string
        path to server list to try out    
  -u string
        username to use for authentication
```

Normal usage:
```bash
ssh-check -u <username> -p <password> -s <path_to_list_of_servers>
```

With debug log:
```bash
ssh-check -u <username> -p <password> -s <path_to_list_of_servers> -d
```