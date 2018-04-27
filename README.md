<p align="center">
  <h3 align="center">G0-PhotoBot</h3>
  <p align="center">Move your Dropbox Camera Uploads to a shared folder.</p>
  <p align="center">
    <a href="https://github.com/garibaldi0/g0-photobot/releases/latest">
      <img alt="Release" src="https://img.shields.io/github/release/garibaldi0/g0-photobot.svg?style=flat-square">
    </a>
    <a href="/LICENSE.md">
      <img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square">
    </a>
    <a href="https://goreportcard.com/report/github.com/garibaldi0/g0-photobot">
      <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/garibaldi0/g0-photobot?style=flat-square">
    </a>
    <a href="http://godoc.org/github.com/garibaldi0/g0-photobot">
      <img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square">
    </a>
    <a href="https://github.com/goreleaser">
      <img alt="Powered By: GoReleaser" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square">
    </a>
  </p>
</p>

* * *

G0-PhotoBot uses API calls to move photos from your Dropbox Camera Uploads folder
to another (shared) folder in your Dropbox.

This project came about due to the pain involved in sharing photos quickly with my
family members.  I didn't like manually having to move photos to different folders
and forgetting to do so for days.

I have this application running on a linux box in cron every 5 minutes.  It then
cycles through the Camera Uploads folder for each of my family members and moves
the photos to our shared folder for everyone to see.

## Installation

1.  Download the current release for your desired platform from [here](https://github.com/garibaldi0/g0-photobot/releases/latest)
2.  Sign up for a Dropbox Developer account [here](https://www.dropbox.com/developers)
3.  Create a Dropbox Application under "My Apps"

    -   Use the Dropbox API (I have not tested with the Dropbox Business API.)
    -   For the Type of Access select "Full Dropbox"

4.  Once your application is created, copy the "App key" and "App secret" values
    to your g0-photobot.toml file.  A default has been provided.
5.  In the g0-photobot.toml file, update the following:
    -   Your LogFile location
    -   Your destination directory
6.  Now you can either place the executable and config file in the same directory
    or the config file can go in one of the following locations

            $HOME/go/etc
            /etc
            ../etc

7.  To get your initial token run with the -n option.

    -   This will open the Dropbox website for you to authorize your app on your account.
    -   Past in the authorization code.
    -   The access token will be printed to the screen.
    -   Update the g0-photobot.toml file with a name and this token in the Tokens section.

8.  You can add multiple tokens to the config file.  This allows for multiple users
    to have their photos moved to the same shared folder.
9.  Next you'll probably want to setup cron or some other scheduling system to run
    the photobot every few minutes.

## Stargazers over time

[![garibaldi0/g0-photobot stargazers over time](https://starcharts.herokuapp.com/garibaldi0/g0-photobot.svg)](https://starcharts.herokuapp.com/garibaldi0/g0-photobot)

* * *

Would you like to fix something in the documentation? Feel free to open an [issue](https://github.com/garibaldi0/g0-photobot/issues).
