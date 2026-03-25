# tui-portfolio

An SSH-accessible terminal portfolio app powered by Bubble Tea and data from Notion.

Live demo: `ssh ssh.chaya.qzz.io`

## Preview

Add a screenshot or GIF here.

`docs/demo.gif`

## Tech Stack

- Go
- Bubble Tea
- Wish (SSH server for Bubble Tea apps)
- Lip Gloss
- Notion API
- fly.io

## Project Structure

- `main.go`: Application entry point. Configures and launches the SSH server using the Wish framework.
- `model.go`: Core state management. Defines the Bubble Tea model, application initialization, and update loops.
- `views.go`: User interface rendering. Handles the layout for all UI elements including the sidebar, content pages, and status bar.
- `styles.go`: Application styling. Contains the LipGloss design definitions and color themes.
- `notion.go`: API integration. Manages fetching and parsing data from the Notion API.
- `go.mod`: Go module definition and direct dependencies.
- `go.sum`: Go dependency checksums and lock file.
- `Dockerfile`: Container configuration for multi-stage application builds.
- `fly.toml`: Deployment configuration for the fly.io platform.
- `.gitignore`: Version control exclusions for secrets and local artifacts.
- `.env.example`: Template outlining required environment variables.

## Local Development

### Prerequisites

- Go 1.23+

### 1. Configure environment

```bash
cp .env.example .env
```

Set `NOTION_API_KEY` in `.env`.

### 2. Create a local SSH host key

```bash
ssh-keygen -t ed25519 -f ./host_key -N ""
```

### 3. Run locally

```bash
export $(grep -v '^#' .env | xargs)
export HOST_KEY_PATH=./host_key
go run .
```

### 4. Connect to the local SSH app

```bash
ssh localhost -p 22
```

If port 22 is unavailable on your machine, run the app in a container or map a different external port to container port 22.

## Environment Variables

| Name | Required | Default | Description |
| --- | --- | --- | --- |
| `NOTION_API_KEY` | Yes | none | Notion integration token used to fetch projects and certifications. |
| `HOST_KEY_PATH` | No | `/data/host_key` | Path to the SSH host private key used by Wish. |

## Deploy to fly.io

1. Install and authenticate Fly CLI:

```bash
fly auth login
```

2. Create the app (one-time):

```bash
fly apps create tui-portfolio
```

3. Create a persistent volume for the SSH host key:

```bash
fly volumes create portfolio_data --region sin --size 1
```

4. Set secrets:

```bash
fly secrets set NOTION_API_KEY=your_notion_api_key
```

5. Deploy:

```bash
fly deploy
```

6. Test:

```bash
ssh ssh.chaya.qzz.io
```

## Updating Content

Projects and certifications are fetched from Notion at runtime. Update content directly in Notion; no code or redeploy is needed for content changes.

## License

MIT
