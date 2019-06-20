<p>
<img src="./public/images/bishack.svg" width="200"/>
<br/><br/>
<span>JOIN US ON SLACK</span>
</p>

&nbsp;

[![CircleCI](https://circleci.com/gh/bis-hack/bishack.dev.svg?style=svg)](https://circleci.com/gh/bis-hack/bishack.dev)
[![codecov](https://codecov.io/gh/bis-hack/bishack.dev/branch/master/graph/badge.svg)](https://codecov.io/gh/bis-hack/bishack.dev)
[![Go Report Card](https://goreportcard.com/badge/github.com/bis-hack/bishack.dev)](https://goreportcard.com/report/github.com/bis-hack/bishack.dev)

&nbsp;

### Folder Structure

	├── assets/             // Static files (go-bindatable/private)
	│   ├── css/
	│   ├── scripts/
	│   └── templates/
	│
	├── handler/            // Handler package.
	|
	├── middleware/         // Middleware central :D
	│
	├── public/             // For publicly available static content.
	│   └── images/
	│
	├── services/           // Anything that does i/o goes here.
	│   └── user/
	│
	├── testing/            // Includes a nifty little tricks for testing.
	│
	└── utils/              // None service related stuff like: session, crypto, etc...
	    └── session/
&nbsp;


### Foreword

This application heavily relies on managed external services like: Cognito, DynamoDB, Lambda et cetera. So setting up a local development environment for this project needs a bit of an extra work.

If you want to contribute to this project, do **[let me know](https://github.com/penzur)** so I can assist you on setting up credentials and what not.

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

### Dev Setup

Install a live-reload command line utility called [**Gin**](https://github.com/codegangsta/gin) with the following command:

	$ go get -u github.com/codegangsta/gin


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
	  DYNAMO_TABLE_POSTS=posts \
	  DYNAMO_ENDPOINT=http://localhost:8000 \
	  make dev



Head to `http://localhost:3000/` on your browser.

&nbsp;

### DynamoDB Setup (local)

> Before you start you need to have **[docker](https://docker.io)** installed on your system first.

1. **Run the the following command:**

	`$ docker run -d -p 8000:8000 amazon/dynamodb-local -jar DynamoDBLocal.jar -inMemory -sharedDb`

	> The above command will run a localized version of dynamodb with shared database.


2. **Go to `http://localhost:8000/shell` and copy, paste and run every files inside the `./assets/dynamo` folder.**

	> See example below:

	![ss](https://cl.ly/2143f27589f7/Screen%252520Recording%2525202019-06-16%252520at%25252001.30%252520AM.gif)


&nbsp;


Happy Hacking! 🖖 🤓

&nbsp;


