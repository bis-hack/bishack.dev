![bishack](./public/images/bishack.png)

&nbsp;

[![CircleCI](https://circleci.com/gh/bis-hack/bishack.dev.svg?style=svg)](https://circleci.com/gh/bis-hack/bishack.dev)
[![codecov](https://codecov.io/gh/bis-hack/bishack.dev/branch/master/graph/badge.svg)](https://codecov.io/gh/bis-hack/bishack.dev)


### Foreword

This application heavily relies on managed external services like: Cognito, DynamoDB, Lambda et cetera. So setting up a local development environment for this project needs a bit of an extra work.

If you want to contribute to this project, do let me know. So I can assist you on setting up credentials and what not.

&nbsp;

### Prerequisites


- [**Go**](https://golang.org) - 1.12 or higher

- **AWS buffet** - you can always sign up for a free tier
	
	- **IAM Key Pair** - to be loaded into your `~/.aws/credentials` config file.
	
		> On mac you can actually just `$ brew install aws-cli` and then run `$ aws config` from your terminal. The prompt will ask you to input the credentials and will load them to the file I mentioned above.
		
	- **Cognito User Pool (with client key and secret)** - you can create one from AWS console.
		
		> Ping me on slack if you need help on this one.
		
	- **DynamoDB Tables** - TBA
	
	
- **Github oauth credentials** this one is easy
- **Slack token** (optional)

&nbsp;

### Setup

Install a live-reload command line utility called [**Gin**](https://github.com/codegangsta/gin) with the following command:

	$ https://github.com/codegangsta/gin
	
	
And then install [**up**](https://up.docs.apex.sh/).

	$ curl -sf https://up.apex.sh/install | sh


Start the server with this command:

	$ SLACK_TOKEN=<slack api token (optional)> \
	  SESSION_KEY=<32-bytes-key> \
	  CSRF_KEY=<32-bytes-key> \
	  COGNITO_CLIENT_ID=<key> \
	  COGNITO_CLIENT_SECRET=<secret> \
	  GITHUB_CLIENT_ID=<id> \
	  GITHUB_CLIENT_SECRET=<secret> \
	  GITHUB_CALLBACK=http://localhost:3000/signup \
	  make dev



Head to `http://localhost:3000/` on your browser.

&nbsp;

> Write your test and submit a pull-request. 🖖 🤓

&nbsp;

Happy Hacking!
