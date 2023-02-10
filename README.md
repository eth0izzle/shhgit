<p align="center">
<img src="./images/logo.png" height="30%" width="30%" />

# 🚨 shhgit is no longer maintained 🚨

## **shhgit helps secure forward-thinking development, operations, and security teams by finding secrets across their code before it leads to a security breach.**

![Go](https://github.com/eth0izzle/shhgit/workflows/Go/badge.svg) ![](https://img.shields.io/docker/cloud/build/eth0izzle/shhgit.svg)

<a href="https://www.shhgit.com?utm_source=github&utm_medium=opensource&utm_campaign=readme"><img src="https://user-images.githubusercontent.com/97316/90076719-5ee09c80-dcf8-11ea-87c7-c5f3b454f246.gif" /></a>
</p>

Accidentally leaking secrets — usernames and passwords, API tokens, or private keys — in a public code repository is a developers and security teams worst nightmare. Fraudsters constantly scan public code repositories for these secrets to gain a foothold in to systems. Code is more connected than ever so often these secrets provide access to private and sensitive data — cloud infrastructures, database servers, payment gateways, and file storage systems to name a few.

shhgit can constantly scan your code repositories to find and alert you of these secrets. 

## Installation

You have two options. I'd recommend the first as it will give you access to the shhgit live web interface. Use the second option if you just want the command line interface.

### via Docker

1. Clone this repository: `git clone https://github.com/eth0izzle/shhgit.git`
2. Build via Docker compose: `docker-compose build`
3. Edit your `config.yaml` file (i.e. adding your GitHub tokens)
4. Bring up the stack: `docker-compose up`
5. Open up http://localhost:8080/

### via Go get

_Note_: this method does not include the shhgit web interface

1. Install [Go](https://golang.org/doc/install) for your platform.
2. `go get github.com/eth0izzle/shhgit` will download and build shhgit automatically. Or you can clone this repository and run `go build -v -i`.
3. Edit your `config.yaml` file and see usage below.

## Usage

shhgit can work in two ways: consuming the public APIs of GitHub, Gist, GitLab and BitBucket  or by processing files in a local directory.

By default, shhgit will run in the former 'public mode'. For GitHub and Gist, you will need to obtain and provide an access token (see [this guide](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line); it doesn't require any scopes or permissions. And then place it under `github_access_tokens` in `config.yaml`). GitLab and BitBucket do not require any API tokens.

You can also forgo the signatures and use shhgit with your own custom search query, e.g. to find all AWS keys you could use `shhgit --search-query AWS_ACCESS_KEY_ID=AKIA`. And to run in local mode (and perhaps integrate in to your CI pipelines) you can pass the `--local` flag (see usage below).

### Options

```
--clone-repository-timeout
        Maximum time it should take to clone a repository in seconds (default 10)
--config-path
        Searches for config.yaml from given directory. If not set, tries to find if from shhgit binary's and current directory
--csv-path
        Specify a path if you want to write found secrets to a CSV. Leave blank to disable
--debug
        Print debugging information
--entropy-threshold
        Finds high entropy strings in files. Higher threshold = more secret secrets, lower threshold = more false positives. Set to 0 to disable entropy checks (default 5.0)
--local
        Specify local directory (absolute path) which to scan. Scans only given directory recursively. No need to have Github tokens with local run.
--maximum-file-size
        Maximum file size to process in KB (default 512)
--maximum-repository-size
        Maximum repository size to download and process in KB) (default 5120)
--minimum-stars
        Only clone repositories with this many stars or higher. Set to 0 to ignore star count (default 0)
--path-checks
        Set to false to disable file name/path signature checking, i.e. just match regex patterns (default true)
--process-gists
        Watch and process Gists in real time. Set to false to disable (default true)
--search-query
        Specify a search string to ignore signatures and filter on files containing this string (regex compatible)
--silent
        Suppress all output except for errors
--temp-directory
        Directory to store repositories/matches (default "%temp%\shhgit")
--threads
        Number of concurrent threads to use (default number of logical CPUs)
```

### Config

The `config.yaml` file has 7 elements. A [default is provided](https://github.com/eth0izzle/shhgit/blob/master/config.yaml).

```
github_access_tokens: # provide at least one token
  - 'token one'
  - 'token two'
webhook: '' # URL to a POST webhook.
webhook_payload: '' # Payload to POST to the webhook URL
blacklisted_strings: [] # list of strings to ignore
blacklisted_extensions: [] # list of extensions to ignore
blacklisted_paths: [] # list of paths to ignore
blacklisted_entropy_extensions: [] # additional extensions to ignore for entropy checks
signatures: # list of signatures to check
  - part: '' # either filename, extension, path or contents
    match: '' # simple text comparison (if no regex element)
    regex: '' # regex pattern (if no match element)
    name: '' # name of the signature
```

#### Signatures

shhgit comes with 150 signatures. You can remove or add more by editing the `config.yaml` file.

```
1Password password manager database file, Amazon MWS Auth Token, Apache htpasswd file, Apple Keychain database file, Artifactory, AWS Access Key ID, AWS Access Key ID Value, AWS Account ID, AWS CLI credentials file, AWS cred file info, AWS Secret Access Key, AWS Session Token, Azure service configuration schema file, Carrierwave configuration file, Chef Knife configuration file, Chef private key, CodeClimate, Configuration file for auto-login process, Contains a private key, Contains a private key, cPanel backup ProFTPd credentials file, Day One journal file, DBeaver SQL database manager configuration file, DigitalOcean doctl command-line client configuration file, Django configuration file, Docker configuration file, Docker registry authentication file, Environment configuration file, esmtp configuration, Facebook access token, Facebook Client ID, Facebook Secret Key, FileZilla FTP configuration file, FileZilla FTP recent servers file, Firefox saved passwords DB, git-credential-store helper credentials file, Git configuration file, GitHub Hub command-line client configuration file, Github Key, GNOME Keyring database file, GnuCash database file, Google (GCM) Service account, Google Cloud API Key, Google OAuth Access Token, Google OAuth Key, Heroku API key, Heroku config file, Hexchat/XChat IRC client server list configuration file, High entropy string, HockeyApp, Irssi IRC client configuration file, Java keystore file, Jenkins publish over SSH plugin file, Jetbrains IDE Config, KDE Wallet Manager database file, KeePass password manager database file, Linkedin Client ID, LinkedIn Secret Key, Little Snitch firewall configuration file, Log file, MailChimp API Key, MailGun API Key, Microsoft BitLocker recovery key file, Microsoft BitLocker Trusted Platform Module password file, Microsoft SQL database file, Microsoft SQL server compact database file, Mongoid config file, Mutt e-mail client configuration file, MySQL client command history file, MySQL dump w/ bcrypt hashes, netrc with SMTP credentials, Network traffic capture file, NPM configuration file, NuGet API Key, OmniAuth configuration file, OpenVPN client configuration file, Outlook team, Password Safe database file, PayPal/Braintree Access Token, PHP configuration file, Picatic API key, Pidgin chat client account configuration file, Pidgin OTR private key, PostgreSQL client command history file, PostgreSQL password file, Potential cryptographic private key, Potential Jenkins credentials file, Potential jrnl journal file, Potential Linux passwd file, Potential Linux shadow file, Potential MediaWiki configuration file, Potential private key (.asc), Potential private key (.p21), Potential private key (.pem), Potential private key (.pfx), Potential private key (.pkcs12), Potential PuTTYgen private key, Potential Ruby On Rails database configuration file, Private SSH key (.dsa), Private SSH key (.ecdsa), Private SSH key (.ed25519), Private SSH key (.rsa), Public ssh key, Python bytecode file, Recon-ng web reconnaissance framework API key database, remote-sync for Atom, Remote Desktop connection file, Robomongo MongoDB manager configuration file, Rubygems credentials file, Ruby IRB console history file, Ruby on Rails master key, Ruby on Rails secrets, Ruby On Rails secret token configuration file, S3cmd configuration file, Salesforce credentials, Sauce Token, Sequel Pro MySQL database manager bookmark file, sftp-deployment for Atom, sftp-deployment for Atom, SFTP connection configuration file, Shell command alias configuration file, Shell command history file, Shell configuration file (.bashrc, .zshrc, .cshrc), Shell configuration file (.exports), Shell configuration file (.extra), Shell configuration file (.functions), Shell profile configuration file, Slack Token, Slack Webhook, SonarQube Docs API Key, SQL Data dump file, SQL dump file, SQLite3 database file, SQLite database file, Square Access Token, Square OAuth Secret, SSH configuration file, SSH Password, Stripe API key, T command-line Twitter client configuration file, Terraform variable config file, Tugboat DigitalOcean management tool configuration, Tunnelblick VPN configuration file, Twilo API Key, Twitter Client ID, Twitter Secret Key, Username and password in URI, Ventrilo server configuration file, vscode-sftp for VSCode, Windows BitLocker full volume encrypted data file, WP-Config
```

## Contributing

1. Fork it, baby!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request.

## Disclaimer

I take no responsibility for how you use this tool. Don't be a dick.

## License

MIT. See [LICENSE](https://github.com/eth0izzle/shhgit/blob/master/LICENSE)
