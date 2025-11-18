## Open WebUI + Ollama (Docker) with optional Web Search

A **local AI** guide to setup use Ollama with Open WebUI & Web Search

---

## Table of contents

1. [Prerequisites](#1-prerequisites)
2. [Quick start (run Open WebUI)](#2-quick-start-run-open-webui)
3. [Connect Open WebUI to Ollama](#3-connect-open-webui-to-ollama)
4. [Recommended prompts (System + User)](#4-recommended-prompts)
5. [Advanced parameters (deterministic vs creative)](#5-advanced-parameters-good-defaults)
6. [Enable Web Search (Google PSE)](#6-enable-web-search-google-pse)
7. [Health check & troubleshooting](#7-health-check--troubleshooting)
8. [Maintenance (update/stop/remove)](#8-maintenance-updatestopremove)

---

## 1) Prerequisites

- Install Ollama: https://ollama.com/download
- Install Docker Desktop: https://docs.docker.com/get-started/get-docker/
- Pull at least one model in Ollama (choose what you need):

```bash
ollama pull mistral:7b          # lightweight general-purpose
ollama pull gpt-oss:20b         # larger reasoning-focused
```

---

## 2) Quick start (run Open WebUI)

```bash
# Pull the latest image
docker pull ghcr.io/open-webui/open-webui:main

# Create a volume to persist Open WebUI data (users, settings)
docker volume create open-webui-data

# Run the container on port 6969
docker run -d -p 6969:8080 \
	-v open-webui-data:/app/backend/data \
	--name open-webui \
	ghcr.io/open-webui/open-webui:main
```

Open http://localhost:6969 and create the admin account when prompted. You can add more users later.

---

## 3) Connect Open WebUI to Ollama

Because Open WebUI runs in Docker and Ollama runs on your macOS host, use the special hostname `host.docker.internal` inside the container.

Steps:
- In Open WebUI, go to: Settings > Admin settings > Connections
- Add a new provider: “Ollama”
- URL: `http://host.docker.internal:11434`
- Click “Test” to verify, then Save

You should see your local models in the model picker. If not, see Troubleshooting below.

Full docs: https://docs.openwebui.com/getting-started/quick-start/starting-with-ollama/

---

## 4) Recommended prompts

Use Markdown for both the System Prompt and your prompts. Keep it short, concrete, and tool-aware.

System Prompt (paste into Settings > General > System Prompt):

```md
You are a pragmatic, accurate AI assistant.

Principles:
- Be concise and actionable; use Markdown formatting.
- Prefer best practices and performance-minded solutions.
- Show trade-offs when they matter.
- Verify factual claims; when using web search, cite 1–3 high-quality sources with titles and links.
- Ask at most one clarifying question only if it unblocks a significantly better answer.
- For code: use fenced blocks with language, minimal dependencies, and brief comments.
```

Optional task-oriented User Prompt template:

```md
Task: <what you want>
Constraints: <limits, versions, APIs, environment>
Deliverable: <format, file(s), steps to run>
Notes: <context, preferences>

If anything is ambiguous, state your single most important assumption and proceed.
```

---

## 5) Advanced parameters (good defaults)

Start conservative and adjust as you go. The exact names can vary by model/backend, but these profiles work well with Ollama-backed models.

- Deterministic (repeatable, precise)
	- temperature: 0.2–0.3
	- top_k: 40–100
	- top_p: 0.9
	- repeat_penalty: 1.1–1.2
	- num_predict (max tokens): 1024–2048 (raise if you expect long outputs)

- Balanced/Creative (more exploration)
	- temperature: 0.6–0.8
	- top_k: 100–200
	- top_p: 0.95
	- repeat_penalty: 1.1
	- num_predict: 1024–4096 (depends on model/context limits)

Notes:
- “Mirostat” is typically off (0) unless you know you want it; temperature/top_p is simpler to tune.
- Bigger “reasoning” models may need higher context or longer outputs; be mindful of speed.

---

## 6) Enable Web Search (Google PSE)

Open WebUI supports multiple search engines. Below uses Google Programmable Search Engine (PSE).

1) Create a Google PSE: https://programmablesearchengine.google.com/controlpanel/create
2) Copy your Search Engine ID (cx)
3) Get an API Key from the Custom Search JSON API: https://developers.google.com/custom-search/v1/introduction

In Open WebUI: Settings > Admin settings > Web Search
- Web Search Engine: `google_pse`
- API Key: your key
- Search Engine ID: your cx
- Suggested:
	- Search Result Count: 5
	- Concurrent Requests: 6

Tips:
- Respect API quotas; start with Result Count = 3 if you hit limits.
- Encourage the model to cite sources when search is used (see System Prompt above).

Full docs: https://docs.openwebui.com/tutorials/web-search/google-pse/

---

## 7) Health check & troubleshooting

Basic checks

```bash
# Is the container up?
docker ps --filter name=open-webui

# View logs
docker logs --tail=200 open-webui
```

Ollama connectivity from container

```bash
# Should return JSON with your local models
docker exec open-webui curl -s http://host.docker.internal:11434/api/tags | jq '.'
```

If models don’t appear in Open WebUI:
- Verify Ollama is running: `ollama serve` (usually started automatically by the app)
- Re-check the URL: `http://host.docker.internal:11434` (macOS-specific hostname for Docker)
- In Open WebUI, re-Test and Save the Ollama connection

Port notes
- UI: http://localhost:6969
- Ollama (default): http://localhost:11434 (from host) / http://host.docker.internal:11434 (from container)

---

## 8) Maintenance (update/stop/remove)

Update image

```bash
docker pull ghcr.io/open-webui/open-webui:main
docker rm -f open-webui && \
docker run -d -p 6969:8080 \
	-v open-webui-data:/app/backend/data \
	--name open-webui \
	ghcr.io/open-webui/open-webui:main
```

Stop / remove

```bash
docker stop open-webui
docker rm open-webui
```

Remove persisted data (irreversible)

```bash
docker volume rm open-webui-data
```

---

Et voilà, happy local chatting!
