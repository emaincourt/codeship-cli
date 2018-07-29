[ ![Codeship Status for emaincourt/codeship-cli](https://app.codeship.com/projects/926858d0-74c5-0136-87ef-06aa68e587da/status?branch=master)](https://app.codeship.com/projects/299734)

# <div style="display:flex; flex-direction:row"><img src="https://i.vimeocdn.com/portrait/5336365_300x300" width="50px" height="50px">&nbsp;codeship-cli</div>

## Description

This CLI tends to be a small tool for monitoring builds on Codeship without need to use the web interface, for efficiency purpose.

<p align="center">
  <img width="60%" src="https://s1.gifyu.com/images/2018-07-28-13.34.42.gif">
</p>

## Installation

```bash
go get github.com/emaincourt/codeship-cli
```

## Usage

There are currently two ways to provide your [Codeship](https://codeship.com/) credentials :

From env vars :
* `CODESHIP_USERNAME`
* `CODESHIP_PASSWORD`

From flags :
* `--username`
* `--password`

> Flags will always be prevalent on env vars

Keep in mind that it is currently not possible to use 2FA with Codeship's API.

Then you can run :

```bash
codeship-cli --org=<YOUR_ORGANIZATION_NAME>
```

