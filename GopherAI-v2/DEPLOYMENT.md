# Local CI/CD Deployment

This deployment runs the Vue frontend, Go API, MySQL, Redis Stack, RabbitMQ, and the MCP weather service on one Docker Compose network. The frontend is available at `http://localhost:8080`; the API health endpoint is `http://localhost:9090/healthz`.

## First deployment

1. Install Docker Desktop and start it.
2. From `GopherAI-v2`, run:

```powershell
.\scripts\deploy-local.ps1
```

The script creates an ignored `.env` with cryptographically random local credentials on its first run, validates configuration, builds images, starts dependencies in health-check order, and waits for the API health endpoint. It preserves named-volume data across normal redeployments. Add `OPENAI_API_KEY`, `OPENAI_MODEL_NAME`, and `OPENAI_BASE_URL` to `.env` before using chat or RAG.

To enable image recognition, set the three ONNX paths in `.env` and deploy with `./scripts/deploy-local.ps1 -WithImageRecognition`. The base deployment intentionally does not mount model files, so chat and user flows can start without them.

## Local CI and rollback

Run the local quality gate before deployment:

```powershell
.\scripts\ci-local.ps1
```

It runs Go tests, containerized frontend lint/build, Compose validation, and image builds. Stop the stack while preserving data with `./scripts/stop-local.ps1`; add `-RemoveData` only when intentionally deleting local MySQL, Redis, RabbitMQ, and uploaded-file volumes. On a self-hosted runner, the script automatically uses `GOPHERAI_LOCAL_ENV_FILE`; otherwise pass `-EnvFile <path>`.

## GitHub Actions CI/CD

`.github/workflows/ci.yml` validates every push and pull request. After a successful push to `main`, `.github/workflows/cd-local.yml` automatically deploys the exact tested commit on a self-hosted GitHub Actions runner labeled `gopherai-local` on the target machine. It can also be started manually with `workflow_dispatch`.

Set the runner-level `GOPHERAI_LOCAL_ENV_FILE` variable to an absolute `.env` path outside the checkout workspace; the workflow copies it only for deployment and removes the working copy afterward. It is never copied to GitHub. The runner must have Docker Desktop, Docker Compose, and permission to bind ports 8080 and 9090.
