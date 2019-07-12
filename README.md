<p align="center">
<img src="https://github.com/lc/secretz/raw/master/secretz.png" alt="secretz" width="250" />
</p>

# secretz
[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/lc/secretz.svg?branch=master)](https://travis-ci.org/lc/secretz)

`secretz` is a tool that minimizes the large attack surface of Travis CI. It automatically fetches repos, builds, and logs for any given organization. 

Built during and for our research on TravisCI: https://edoverflow.com/2019/ci-knew-there-would-be-bugs-here/


## Usage:
`secretz -t Organization [options]`


### Flags:
| Flag | Description | Example |
|------|-------------|---------|
| `-t` | Organization to get repos, builds, and logs for | `secretz -t ExampleCo` |
| `-c` | Limit the number of workers that are spawned | `secretz -t ExampleCo -c 3` |
| `-delay` | delay between requests + random delay/2 jitter | `secretz -t ExampleCo -delay 900`|
| `-members [list \| scan]` | Get all GitHub members belonging to Organization and list/scan them | `secretz -t ExampleCo -members scan` |
| `-timeout` | How long to wait for HTTP Responses from Travis CI | `secretz -t ExampleCo -timeout 20` |
| `-setkey` | Set API Key for api.travis-ci.org | `secretz -setkey yourapikey` |

## Installation:

### Via `go get`
```
go get -u github.com/lc/secretz
```

### Via `git clone`

```
go get -u github.com/json-iterator/go
git clone git@github.com:lc/secretz
cd secretz && go build -o secretz main.go
```


### Generate an API-Key: 
```
travis login
travis token --org
```

### Create config file
`secretz -setkey <API-KEY>`


### Note:
Please keep your delay high and your workers low out of respect for TravisCI and their APIs. This will also help you from being rate-limited by them. 
