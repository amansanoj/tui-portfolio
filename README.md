![OG Image](https://repo-og-generator.vercel.app/repo-og-generator?description=An%20SSH-accessible%20terminal%20portfolio%20app%20powered%20by%20Bubble%20Tea%20and%20data%20from%20Notion&url=ssh%20whoami.chaya.qzz.io&scale=2)

## Index
- [Live Demo](#live-demo)
- [Features](#features)
- [Project Structure](#project-structure)
- [Updating Content](#updating-content)
- [Local Development](#local-development)
- [Deployment](#deployment)
- [Credits & License](#credits--license)

## Live Demo
You can experience the interactive portfolio directly from your terminal. No installation or configuration is required:

```sh
ssh whoami.chaya.qzz.io
```

## Features
- Provides an interactive, SSH-accessible terminal interface built natively in Go
- Fetches project and certification data dynamically at runtime via the Notion API
- Eliminates the need for code redeployments when portfolio content is updated in Notion
- Utilizes Lip Gloss to deliver beautiful, customizable terminal layouts and color themes
- Integrates seamlessly with the Wish framework to serve Bubble Tea apps over SSH

## Project Structure
```text
.
├── .env.example
├── Dockerfile
├── README.md
├── content_store.go
├── fly.toml
├── go.mod
├── go.sum
├── main.go
├── model.go
├── notion.go
├── styles.go
└── views.go
```
`main.go` acts as the primary entry point and sets up the SSH server via Wish. Most UI logic and customization lives across `views.go` and `styles.go`.

## Updating Content
For updating your live portfolio, you do not need to redeploy the codebase:

- The TUI fetches records natively from your Notion databases based on `NOTION_API_KEY`
- Manage text, projects, and certifications directly within Notion
- The app automatically syncs fresh content on an interval defined by `NOTION_REFRESH_SECONDS`

If you want to customize the design or structure of the application itself:
- `styles.go`: modify Lip Gloss color definitions and dimensional stylings
- `views.go`: adjust the rendering logic for sidebars, content tabs, and footers
- `model.go`: modify the core Bubble Tea state management

## Local Development
To customize the repository or run your own instance locally, ensure you have Go (1.23+) installed and generate a local host key for the SSH server to bind to:

```sh
# Copy environment template
cp .env.example .env

# Generate a local SSH host key
ssh-keygen -t ed25519 -f ./host_key -N ""

# Export environment variables and run
export $(grep -v '^#' .env | xargs)
export HOST_KEY_PATH=./host_key
go run .
```
You can then access your local TUI by running `ssh localhost -p 2222`.

## Deployment
Deployment is streamlined via Fly.io. Because Fly.io application names must be globally unique, you will need to choose your own unique app name and update the `app` field in your `fly.toml` file before deploying. 

After setting up the Fly CLI, authenticating, and updating your `.toml` file, you can provision the app, mount a persistent volume for the SSH host key, and push to production:

```sh
# Replace 'your-unique-app-name' with the name you set in fly.toml
fly apps create your-unique-app-name
fly volumes create portfolio_data --region ams --size 1
fly secrets set NOTION_API_KEY=your_notion_api_key
fly deploy --remote-only
```
The application will instantly become available via SSH using your assigned fly domain or configured custom domain (like `whoami.chaya.qzz.io`). 

## Credits & License
Constructed using Go, Bubble Tea, Wish, and Lip Gloss. This codebase is open-source. Please check the included LICENSE file for redistribution rights and terms.
