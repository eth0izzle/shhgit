# **shhgit**: find GitHub secrets in real time

**Shhgit finds secrets and sensitive files across GitHub code and Gists committed in *near* real time by listening to the [GitHub Events API](https://developer.github.com/v3/activity/events/).**

<p align="center">
<img src="https://www.darkport.co.uk/assets/img/shhgit.png" alt="shhgit" width="200" />
</p>

## **[NEW: LIVE VERSION. Find GitHub secrets straight from your browser!](https://shhgit.darkport.co.uk)**

Finding secrets in GitHub is nothing new. There are many great tools available to help with this depending on which side of the fence you sit. On the adversarial side, popular tools such as <a href="https://github.com/michenriksen/gitrob">gitrob</a> and <a href="https://github.com/dxa4481/truffleHog">truggleHog</a> focus on digging in to commit history to find secret tokens from specific repositories, users or organisations. On the defensive side, GitHub themselves are actively scanning for secrets through their [token scanning](https://help.github.com/en/articles/about-token-scanning) project. Their objective is to identify secret tokens within committed code in real-time and notify the service provider to action. So in theory if any AWS secret keys are committed to GitHub, Amazon will be notified and automatically revoke them.

I developed shhgit to raise awareness and bring to life the prevalence of this issue. I hope GitHub will do more to prevent bad actors using the treasure trove of information across the platform. I don't know the inner-workings of their [token scanning](https://help.github.com/en/articles/about-token-scanning) project but delaying the real-time feed API until the pipeline has completed and posing SLAs on the providers seems like a step in the right direction.

**With some tweaking of the signatures shhgit would make an awesome addition to your bug bounty toolkit.**

<img src="https://www.darkport.co.uk/assets/img/shhgit-example.png" alt="shhgit" />
<img src="https://www.darkport.co.uk/assets/img/shhgit-live-example.png" alt="shhgit live!" />

## Installation

You can use the [precompiled binaries](https://www.github.com/eth0izzle/shhgit/releases) **or** build from source:

1. Install [Go](https://golang.org/doc/install) for your platform.
2. `$ go get github.com/eth0izzle/shhgit` will download and build shhgit.
3. See usage.

## Usage

shhgit needs to access the public GitHub API so you will need to obtain and provide an access token. The API has a hard rate limit of 5,000 requests per hour per account, regardless what token is used. The more account-unique tokens you provide, the faster you can process the events. Follow [this guide](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line) to generate a token; it doesn't require any scopes or permissions. And then place it under `github_access_tokens` in `config.yaml`. **Note that it is against the GitHub terms to bypass their rate limits. Use multiple tokens at your own risk**.

Unlike other tools, you don't need to pass any targets with shhgit. Simply run `$ shhgit` to start watching GitHub commits and find secrets or sensitive files matching the included 120 signatures.

Alternatively, you can forgo the signatures and use shhgit with a search query, e.g. to find all AWS keys you could use `shhgit --search-query AWS_ACCESS_KEY_ID=AKIA`

### Options

```
--check-owner
        Will check owner details before processing repo. Set to true to enable.
--clone-repository-timeout
        Maximum time it should take to clone a repository in seconds (default 10)
--csv-path
        Specify a path if you want to write found secrets to a CSV. Leave blank to disable
--debug
        Print debugging information
--entropy-threshold
        Finds high entropy strings in files. Higher threshold = more secret secrets, lower threshold = more false positives. Set to 0 to disable entropy checks (default 5.0)
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

The `config.yaml` file has 6 elements. A [default is provided](https://github.com/eth0izzle/shhgit/blob/master/config.yaml).

```
github_access_tokens: # provide at least one token
  - 'token one'
  - 'token two'
slack_webhook: '' # url to your slack webhook. Found secrets will be sent here
blacklisted_extensions: [] # list of extensions to ignore
blacklisted_paths: [] # list of paths to ignore
blacklisted_entropy_extensions: [] # additional extensions to ignore for entropy checks
organizations: []
signatures: # list of signatures to check
  - part: '' # either filename, extension, path or contents
    match: '' # simple text comparison (if no regex element)
    regex: '' # regex pattern (if no match element)
    name: '' # name of the signature
```

#### Signatures

shhgit comes with 120 signatures. You can remove or add more by editing `config.yaml`.

```
Chef private key, Potential Linux shadow file, Potential Linux passwd file, Docker configuration file, NPM configuration file, Environment configuration file, Contains a private key, AWS Access Key ID Value, AWS Access Key ID, AWS Account ID, AWS Secret Access Key, AWS Session Token, Artifactory, CodeClimate, Facebook access token, Google (GCM) Service account, Stripe API key, Google OAuth Key, Google Cloud API Key
Google OAuth Access Token, Picatic API key, Square Access Token, Square OAuth Secret, PayPal/Braintree Access Token, Amazon MWS Auth Token, Twilo API Key, MailGun API Key, MailChimp API Key, SSH Password, Outlook team, Sauce Token, Slack Token, Slack Webhook, SonarQube Docs API Key, HockeyApp, Username and password in URI, NuGet API Key, Potential cryptographic private key, Log file, Potential cryptographic key bundle, Potential cryptographic key bundle
Potential cryptographic key bundle, Potential cryptographic key bundle, Pidgin OTR private key, OpenVPN client configuration file, Azure service configuration schema file, Remote Desktop connection file, Microsoft SQL database file, Microsoft SQL server compact database file, SQLite database file, SQLite3 database file, Microsoft BitLocker recovery key file
Microsoft BitLocker Trusted Platform Module password file, Windows BitLocker full volume encrypted data file, Java keystore file, Password Safe database file, Ruby On Rails secret token configuration file, Carrierwave configuration file, Potential Ruby On Rails database configuration file, OmniAuth configuration file, Django configuration file
1Password password manager database file, Apple Keychain database file, Network traffic capture file, GnuCash database file, Jenkins publish over SSH plugin file, Potential Jenkins credentials file, KDE Wallet Manager database file, Potential MediaWiki configuration file, Tunnelblick VPN configuration file, Sequel Pro MySQL database manager bookmark file, Little Snitch firewall configuration file, Day One journal file, Potential jrnl journal file, Chef Knife configuration file, cPanel backup ProFTPd credentials file
Robomongo MongoDB manager configuration file, FileZilla FTP configuration file, FileZilla FTP recent servers file, Ventrilo server configuration file, Terraform variable config file, Shell configuration file, Shell configuration file, Shell configuration file, Private SSH key, Private SSH key, Private SSH key, Private SSH key, SSH configuration file, Potential cryptographic private key, Shell command history file
MySQL client command history file, PostgreSQL client command history file, PostgreSQL password file, Ruby IRB console history file, Pidgin chat client account configuration file, Hexchat/XChat IRC client server list configuration file, Irssi IRC client configuration file, Recon-ng web reconnaissance framework API key database, DBeaver SQL database manager configuration file, Mutt e-mail client configuration file, S3cmd configuration file, AWS CLI credentials file, SFTP connection configuration file, T command-line Twitter client configuration file, Shell configuration file
Shell profile configuration file, Shell command alias configuration file, PHP configuration file, GNOME Keyring database file, KeePass password manager database file, SQL dump file, Apache htpasswd file, Configuration file for auto-login process, Rubygems credentials file, Tugboat DigitalOcean management tool configuration, DigitalOcean doctl command-line client configuration file, git-credential-store helper credentials file, GitHub Hub command-line client configuration file, Git configuration file
```

## Contributing

1. Fork it, baby!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request.

## Credits

Some code borrowed from [Gitrob](https://github.com/michenriksen/gitrob) by [Michael Henriksen](https://michenriksen.com/).

## Disclaimer

I take no responsibility for how you use this tool. Don't be a dick.

## License

MIT. See [LICENSE](https://github.com/eth0izzle/shhgit/blob/master/LICENSE)
