# <div style="display:flex; flex-direction:row"><img src="https://i.vimeocdn.com/portrait/5336365_300x300" width="50px" height="50px"><p>&nbsp;codeship-cli<p/></div>

## Description

This CLI tends to be a small tool for monitoring builds on Codeship without need to use the web interface, for efficiency purpose.

## Installation

```bash
go get github.com/emaincourt/codeship-cli
```

## Usage

You currently need to set the following env vars to allow access to [Codeship](https://codeship.com/) :
* `CODESHIP_USERNAME`
* `CODESHIP_PASSWORD`

Keep in mind that it is currently not possible to use 2FA with Codeship's API.

Then you can run :

```bash
codeship-cli --org=<YOUR_ORGANIZATION_NAME>
```
