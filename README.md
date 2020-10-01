<p align="center">
  <img alt="pwndrop logo" src="https://raw.githubusercontent.com/kgretzky/pwndrop/master/media/pwndrop-logo-512.png" height="120" />
  <p align="center">
    <img alt="pwndrop title" src="https://raw.githubusercontent.com/kgretzky/pwndrop/master/media/pwndrop-title-black-512.png" height="40" />
  </p>
</p>

**pwndrop** is a self-deployable file hosting service for sending out red teaming payloads or securely sharing your private files over HTTP and WebDAV.

If you've ever needed to quickly set up an nginx/apache web server to host your files and you were never happy with the limitations of `python -m SimpleHTTPServer`, **pwndrop** is definitely for you!

<p align="center">
  <img alt="demo" src="https://raw.githubusercontent.com/kgretzky/pwndrop/master/media/demo1.gif" height="500" />
</p>

With **pwndrop** you can:
- [x] Upload and immediately share multiple files using your own private VPS, using drag & drop.
- [x] Decide to make files available or unavailable for download with a single click.
- [x] Set up custom download URLs, for shared files, without playing with directory structure.
- [x] Set up facade files, which will be served instead of the original file whenever you feel like it.
- [x] Set up automatic redirects to spoof the file's extension in a shared link.
- [x] Change MIME type of the served file to change browser's behavior when a download link is clicked.
- [x] Serve files over HTTP, HTTPS and WebDAV.
- [x] Install and setup everything using a bash oneliner.
- [x] Set up **pwndrop** to work as a nameserver and respond with a valid DNS A record to any sub-domain you choose.
- [x] Protect your admin panel behind a custom secret URL path and log in securely with your own username and password.
- [x] Never worry about setting up HTTPS certificates as **pwndrop** does everything for you in the background (including auto-renewals).

Its main goal is to make file sharing as easy and intuitive as possible, while implementing extra features to aid in red team assessments.

Frontend of **pwndrop** is developed in pure Vue.js + Bootstrap with no npm or webpack dependencies. The backend serves REST API and manages a local database, powered by GO language.

## Write-up

If you want to learn how to use **pwndrop** or you want to learn what new features were implemented in recent releases, make sure to check out the posts on my blog:

https://breakdev.org/pwndrop

## Video guide

Take a look at the fantastic video made by Luke Turvey ([@TurvSec](https://twitter.com/TurvSec)), which fully explains how to get started using **pwndrop**.

[![File and Phishing Payload Hosting using PwnDrop (Red Team) - Luke Turvey](https://img.youtube.com/vi/e3veSyIFvOE/0.jpg)](https://www.youtube.com/watch?v=e3veSyIFvOE)

## Prerequisites

If you don't yet have the server to deploy to I highly recommend Digital Ocean. The cheapest $5/mo Debian 9 server with 25GB of storage space will work wonders for you. You can use my referral link to [get an extra $100 to spend on your servers in 60 days for free](https://m.do.co/c/50338abc7ffe).

Register a new domain and point its DNS A records to your VPS IP. You can also register a domain and point its `ns1` and `ns2` nameservers to **pwndrop** instance IP - it will automatically respond with valid DNS A replies.

1. Registered domain name pointing to **pwndrop** instance IP as a DNS A records or as a nameserver.
2. Server with at least 512 MB RAM.

If you want to set up **pwndrop** without a domain, check below how to set up a local instance, which will not auto-generate HTTPS certificates.

## Installation

Make sure there aren't any DNS or HTTP(S) servers running before you attempt to install **pwndrop**.

#### Oneliner

I do not recommend running oneliners, before downloading and checking the script code, but if you are really in a hurry, here it is:
```
curl https://raw.githubusercontent.com/kgretzky/pwndrop/master/install_linux.sh | sudo bash
```

This will download the latest amd64 release binary and fully install a daemon running in a background.

#### From binary

First you need to download the release package you want from: https://github.com/kgretzky/pwndrop/releases

Then do the following (this performs same actions to the oneliner):

```
tar zxvf pwndrop-linux-amd64.tar.gz
./pwndrop stop
./pwndrop install
./pwndrop start
./pwndrop status
```

#### From source code

First of all, make sure you have installed GO with version at least **1.13**: https://golang.org/doc/install

Then do the following:

```
sudo apt-get -y install git make
git clone https://github.com/kgretzky/pwndrop
cd pwndrop
make
make install
```

## Quickstart

Make sure the **pwndrop** is running.

1. Open the secret URL to authorize your browser: `https://yourdomain.com/pwndrop` (this is a default value; make sure to use the secret path, you've pre-configured)
2. Open the admin panel URL in your browser: `https://yourdomain.com/` (since you've authorized your browser, you will now see an admin panel login page)
3. Create your admin account or login.
4. Click the configuration cog in top-left corner and make sure you change the secret path to something other than `/pwndrop`.

You're good to go!

## Running from CLI

You don't have to install **pwndrop** as a daemon and you can run it straight from the console.

```
usage: pwndrop [start|stop|install|remove|status] [-config <config_path>] [-debug] [-no-autocert] [-no-dns] [-h]

daemon management:
    start           : start the daemon
    stop            : stop the daemon
    install         : install the daemon using the available system manager (systemd, systemv and upstart supported)
    remove          : uninstall the daemon
    status          : check status of the installed daemon

parameters:
    -config         : specify a custom path to a config file (def. 'pwndrop.ini' in same directory as the executable)
    -debug          : enable debug output 
    -no-autocert    : disable automatic TLS certificate retrieval from LetsEncrypt; useful when you want to connect over IP or/and in a local network
    -no-dns         : do not run a DNS server on port 53 UDP; use this if you don't want to use pwndrop as a nameserver
    -h              : usage help
```

## Configuration

On first launch, **pwndrop**, by default, will create a new configuration file `pwndrop.ini` in the same directory as an executable. You can later modify it or supply your own, for example to pre-configure **pwndrop** before the installation to automate the deployment of a tool even better.

Here is an example config file with all available config variables with commentary:
```
[pwndrop]
listen_ip = "190.33.86.22"                  # the external IP of your pwndrop instance (must be set if you want to use the nameserver feature)
http_port = 80                              # listening port for HTTP and WebDAV
https_port = 443                            # listening port for HTTPS
data_dir = "./data"                         # directory path where data storage will reside (relative paths are from executable directory path)
admin_dir = "./admin"                       # directory path where the admin panel files reside (relative paths are from executable directory path)

[setup]                                     # optional: put in if you want to pre-configure pwndrop (section will be deleted from the config file on first run)
username = "admin"                          # username of the admin account
password = "secretpassword"                 # password of the admin account
redirect_url = "https://www.somedomain.com" # URL to which visitors will be redirected to if they supply a path, which doesn't point to any shared file (put blank if you want to return 404)
secret_path = "/pwndrop"                    # secret URL path, which upon visiting will allow your browser to access the login page of the admin panel (make sure to change the default value)
```

If you want to pre-configure your **pwndrop** instance before deployment using any of the installation scripts, put your configuration file at `/usr/local/pwndrop/pwndrop.ini` and it will be parsed the moment **pwndrop** daemon is first executed.

## Credits

Huge thanks to [**@jaredhaight**](https://twitter.com/jaredhaight) for inspiring me to learn Vue, with his [Faction C2](https://www.factionc2.com/) framework!

Also much thanks to all the people who gave me pre-release feedback and supported me with their opinions on the tool!

## License

**pwndrop** is made by Kuba Gretzky ([@mrgretzky](https://twitter.com/mrgretzky)) and it's released under GPL3 license.
