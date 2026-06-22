# gomult

Multi-language code compilation and execution service. Single Go binary, per-language Docker images for minimal footprint.

## Languages

### Base (28)

| Key | Language | Key | Language | Key | Language | Key | Language |
|-----|----------|-----|----------|-----|----------|-----|----------|
| `py` | Python 3 | `js` | JavaScript | `rb` | Ruby | `php` | PHP |
| `pl` | Perl | `lua` | Lua | `sh` | Bash | `ts` | TypeScript |
| `ex` | Elixir | `pro` | Prolog | `r` | R | `rkt` | Racket |
| `groovy` | Groovy | `kt` | Kotlin | `c` | C | `cpp` | C++ |
| `go` | Go | `java` | Java 17 | `rs` | Rust | `cs` | C# |
| `d` | D | `hs` | Haskell | `ml` | OCaml | `pas` | Pascal |
| `f90` | Fortran | `erl` | Erlang | `scala` | Scala | `asm` | Assembly |

### Versioned

| Image | Language | Runtime |
|-------|----------|---------|
| `gomult-py2` | Python 2.7 | Debian bullseye |
| `gomult-java8` | Java 8 | Adoptium |
| `gomult-java11` | Java 11 | Adoptium |
| `gomult-java17` | Java 17 | Adoptium |
| `gomult-java21` | Java 21 | Adoptium |
| `gomult-node20` | Node.js 20 | Official binary |
| `gomult-node22` | Node.js 22 | Official binary |

**35 images total.** Each contains one Go binary + one language runtime.

## Quick Start

### Per-language Docker (recommended)

```bash
# Build one language
docker build -f docker/Dockerfile.py -t gomult-py .
docker run -p 8080:8080 gomult-py

# Build a versioned language  
docker build -f docker/Dockerfile.java21 -t gomult-java21 .
docker run -p 8080:8080 gomult-java21
```

### All-in-one (28 languages, large image)

```bash
docker build -t gomult .          # uses root Dockerfile + config.yaml
docker run -p 8080:8080 gomult
```

### Local dev (no Docker, requires local runtimes)

```bash
go build -o gomult .
./gomult
```

## API

`POST /compile` — JSON in, plain text out.

```bash
curl -X POST http://localhost:8080/compile \
  -H "Content-Type: application/json" \
  -d '{"code":"print(42)","language":"py"}'
```

| Field | Required | Description |
|-------|----------|-------------|
| `code` | yes | Source code |
| `input` | no | Stdin for the program |
| `language` | yes | Language key (see tables above) |
| `analyze` | no | Set `true` to get AI analysis via Deepseek (requires `--ai-key` flag or `DEEPSEEK_API_KEY` env var) |

Successful requests return program output. Errors are prefixed:
- `[COMPILE ERROR]` — compilation failed
- `[RUNTIME ERROR]` — non-zero exit code
- `[TIME LIMIT EXCEEDED]` — execution timed out

`GET /health` → `{"status":"ok"}`  
`GET /languages` → JSON map of supported keys

## AI Analysis

When enabled, the server can send your code and its execution output to Deepseek for analysis. The response includes bug identification, fix suggestions (on errors), or code quality feedback (on success).

```bash
# Start server with AI enabled
./gomult --ai-key sk-your-deepseek-key
# or
DEEPSEEK_API_KEY=sk-... ./gomult
```

```bash
# Request analysis
curl -X POST http://localhost:8080/compile \
  -H "Content-Type: application/json" \
  -d '{"code":"print(1/0)","language":"py","analyze":true}'
```

The AI analysis appears in the response under `[AI ANALYSIS]`. If the API call fails, the error is logged server-side and `[AI ANALYSIS FAILED]` is returned — the execution result is never blocked.

## Project Structure

```
gomult/
├── main.go                          # Entry point
├── go.mod / go.sum
├── config.yaml                      # All-in-one config (28 langs)
├── Dockerfile                       # All-in-one Dockerfile
├── internal/
│   ├── config/config.go             # YAML config loading
│   ├── sandbox/                     # Execution engine
│   │   ├── sandbox.go               #   Types + Executor interface
│   │   ├── engine.go                #   Compile → run orchestration
│   │   ├── executor_direct.go       #   Direct execution (dev/fallback)
│   │   └── executor_nsjail.go       #   nsjail sandbox (production)
│   ├── ai/ai.go                       # Deepseek analysis client
│   └── server/server.go               # HTTP handlers
├── configs/                         # Per-language configs
│   ├── py.yaml, c.yaml, java.yaml ...
│   └── py2.yaml, java8.yaml ...     # Version-specific
└── docker/                          # Per-language Dockerfiles
    ├── Dockerfile.py, Dockerfile.c ...
    └── Dockerfile.java21, Dockerfile.node22 ...
```

## Configuration

### CLI flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to config file |
| `--ai-key` | (env `DEEPSEEK_API_KEY`) | Deepseek API key for AI analysis |

Per-language configs follow the same format. Key settings:

```yaml
server:
  port: 8080
  max_code_size: 1048576          # 1 MB
sandbox:
  mode: direct                    # direct | nsjail | auto
  time_limit: 5                   # seconds
  max_memory: 268435456           # 256 MB
  max_file_size: 10485760         # 10 MB
  max_processes: 32
languages:
  py:
    name: Python 3
    extension: .py
    run_cmd: [python3, "{file}"]  # omit compile_cmd for interpreted
```

### Adding a language

1. Create `configs/zig.yaml` with a language block
2. Create `docker/Dockerfile.zig` installing the runtime + copying the config

```dockerfile
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y zig
COPY --from=builder /app/gomult /usr/local/bin/gomult
COPY configs/zig.yaml /etc/gomult/config.yaml
...
```

Placeholders in commands: `{file}` (source path), `{output}` (binary path), `{dir}` (work directory).

## Sandbox

Set `sandbox.mode: nsjail` in config for production. Requires nsjail installed in the image. Adds network isolation, cgroup memory limits, and seccomp filtering. Falls back to `direct` mode when nsjail is unavailable.

## License

Apache 2.0 — see [LICENSE](./LICENSE)
