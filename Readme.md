# gomult

Multi-language code compilation and execution service with nsjail sandboxing.

## Supported Languages (28)

| Key | Language | Type |
|-----|----------|------|
| `asm` | Assembly (NASM x86-64) | Compiled |
| `c` | C | Compiled |
| `cpp` | C++ | Compiled |
| `cs` | C# (Mono) | Compiled |
| `d` | D | Compiled |
| `erl` | Erlang | Compiled |
| `ex` | Elixir | Interpreted |
| `f90` | Fortran | Compiled |
| `go` | Go | Compiled |
| `groovy` | Groovy | Interpreted |
| `hs` | Haskell | Compiled |
| `java` | Java | Compiled |
| `js` | JavaScript | Interpreted |
| `kt` | Kotlin | Interpreted |
| `lua` | Lua | Interpreted |
| `ml` | OCaml | Compiled |
| `pas` | Pascal | Compiled |
| `php` | PHP | Interpreted |
| `pl` | Perl | Interpreted |
| `pro` | Prolog | Interpreted |
| `py` | Python 3 | Interpreted |
| `r` | R | Interpreted |
| `rb` | Ruby | Interpreted |
| `rkt` | Racket | Interpreted |
| `rs` | Rust | Compiled |
| `scala` | Scala | Compiled |
| `sh` | Bash | Interpreted |
| `ts` | TypeScript | Interpreted |

## Features

- Single binary, config-driven — add a language with one YAML block
- nsjail sandbox: network isolation, rlimits (memory, filesize, nproc), user separation
- Fallback to direct execution when nsjail is unavailable
- Health check and language listing endpoints

## Quick Start

```bash
go build -o gomult .
./gomult
```

Starts on `http://localhost:8080` in `direct` sandbox mode (no nsjail needed for local dev).

### Docker

```bash
docker build -t gomult .
docker run -p 8080:8080 gomult
```

The Docker image includes nsjail and all 28 language runtimes.

## API

### `POST /compile`

```bash
curl -X POST http://localhost:8080/compile \
  -H "Content-Type: application/json" \
  -d '{"code":"print(\"hello\")","input":"","language":"py"}'
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `code` | string | yes | Source code |
| `input` | string | no | Stdin input |
| `language` | string | yes | Language key (see table) |

Response is `text/plain`. Errors are prefixed: `[COMPILE ERROR]`, `[RUNTIME ERROR]`, `[TIME LIMIT EXCEEDED]`.

### `GET /health` — `{"status":"ok"}`

### `GET /languages` — JSON map of language keys to names

## Configuration

See `config.yaml`. Key settings:

| Setting | Default | Description |
|---|---|---|
| `server.port` | 8080 | Listen port |
| `server.max_code_size` | 1048576 | Max code length (1 MB) |
| `sandbox.mode` | `auto` | `nsjail`, `direct`, or `auto` |
| `sandbox.time_limit` | 5 | Timeout (seconds) |
| `sandbox.max_memory` | 268435456 | Memory limit (256 MB) |
| `sandbox.max_file_size` | 10485760 | Max file size (10 MB) |
| `sandbox.max_processes` | 32 | Max processes per run |
| `sandbox.runtime_user` | 1000 | UID for sandboxed processes |

### Adding a language

```yaml
languages:
  zig:
    name: Zig
    extension: .zig
    compile_cmd: [zig, build-exe, "{file}"]
    run_cmd: ["{output}"]
```

Placeholders: `{file}` (source path), `{output}` (binary path), `{dir}` (work directory). Omit `compile_cmd` for interpreted languages. Set `source_file` to override the written filename (defaults to `code.<ext>`).

## License

Apache 2.0 — see [LICENSE](./LICENSE)
