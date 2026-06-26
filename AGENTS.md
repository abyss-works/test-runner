# test-runner

test-runner — scaffolded by idp-starter

## Tech Stack

- **Backend**: Go 1.23, net/http
- **Frontend**: Vue 3, Vite, Tailwind CSS
- **Database**: PostgreSQL
- **Deployment**: Docker, Kubernetes

## Getting Started

```bash
# Frontend dev server
cd frontend && npm install && npm run dev

# Go server (development)
go run .

# Docker
docker build -t test-runner .
docker run -p 8080:8080 test-runner
```

## Configuration

Copy `config.example.yaml` to `config.yaml` and adjust values for your environment.

## CI/CD

GitHub Actions workflows are in `.github/workflows/`. Replace placeholder values in `k8s/` manifests before deploying.

---

*This project was scaffolded by the idp-starter platform.*
