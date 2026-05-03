# stmt2redis

[![CI](https://github.com/its-the-vibe/stmt2redis/actions/workflows/ci.yaml/badge.svg)](https://github.com/its-the-vibe/stmt2redis/actions/workflows/ci.yaml)

A Go command-line utility that parses bank statement CSV files and streams each transaction as a JSON object into a Redis list for downstream processing.

## Supported CSV Types

| `--type`    | Bank / Product          |
|-------------|-------------------------|
| `starling`  | Starling Bank           |
| `amex`      | American Express        |
| `monzo`     | Monzo                   |
| `monzo-flex`| Monzo Flex              |

## Prerequisites

- Go 1.24+
- A running Redis instance (unless using `--stdout`)

## Setup

### 1. Clone and build

```bash
git clone https://github.com/its-the-vibe/stmt2redis.git
cd stmt2redis
make build
```

### 2. Create your configuration file

Copy the example config and edit it with your Redis settings and desired list names:

```bash
cp config.example.yaml config.yaml
```

`config.yaml` (gitignored — never committed):

```yaml
redis:
  host: localhost
  port: 6379
  db: 0

lists:
  starling: transactions:starling
  amex: transactions:amex
  monzo: transactions:monzo
  monzo_flex: transactions:monzo-flex
```

### 3. Set up your environment file

Copy the example `.env` and add your Redis password (leave blank if none):

```bash
cp .env.example .env
```

`.env` (gitignored — never committed):

```
REDIS_PASSWORD=your_redis_password_here
```

## Usage

### Push transactions to Redis

```bash
./stmt2redis push --type starling --file statement.csv
./stmt2redis push --type amex     --file amex.csv
./stmt2redis push --type monzo    --file monzo.csv
./stmt2redis push --type monzo-flex --file monzo_flex.csv
```

### Print JSON to stdout (no Redis required)

```bash
./stmt2redis push --type starling --file statement.csv --stdout
```

### Use a custom config or env file

```bash
./stmt2redis push --type monzo --file monzo.csv --config /path/to/config.yaml --env-file /path/to/.env
```

### All flags

| Flag         | Short | Default       | Description                                         |
|--------------|-------|---------------|-----------------------------------------------------|
| `--type`     | `-t`  | *(required)*  | CSV type: `starling`, `amex`, `monzo`, `monzo-flex` |
| `--file`     | `-f`  | *(required)*  | Path to the CSV file                                |
| `--stdout`   |       | `false`       | Print JSON to stdout instead of pushing to Redis    |
| `--config`   |       | `config.yaml` | Path to the YAML config file                        |
| `--env-file` |       | `.env`        | Path to the `.env` file                             |

## Makefile targets

| Target  | Description               |
|---------|---------------------------|
| `build` | Compile the binary        |
| `test`  | Run all unit tests        |
| `lint`  | Run `go vet`              |
| `clean` | Remove the compiled binary|

```bash
make build
make test
make lint
make clean
```

## JSON Output

Each transaction is output as a JSON object with fields derived from the CSV columns (see [CSV Headers](#csv-headers) below). Two additional fields are always included:

- `"filename"` – the base name of the source CSV file, useful for data provenance when processing multiple files.
- `"index"` – the 0-based position of the record within the source CSV file (first data row is `0`), allowing downstream consumers to reconstruct the original order or detect missing/duplicate records.

For example:

```json
{
  "date": "2026-03-21",
  "counter_party": "Example Ltd",
  "amount_gbp": "100.00",
  "filename": "statement.csv",
  "index": 0
}
```

Both fields are present in Redis and stdout output modes.

## CSV Headers

**Starling:**
```
Date,Counter Party,Reference,Type,Amount (GBP),Balance (GBP),Spending Category,Notes
```

> **Date handling:** Starling exports dates in `DD/MM/YYYY` format (e.g. `21/03/2026`). The parser automatically converts these to ISO `YYYY-MM-DD` format (e.g. `2026-03-21`) in the output JSON.

**Amex:**
```
Date,Description,Amount,Extended Details,Appears On Your Statement As,Address,Town/City,Postcode,Country,Reference,Category
```

**Monzo / Monzo Flex:**
```
Transaction ID,Date,Time,Type,Name,Emoji,Category,Amount,Currency,Local amount,Local currency,Notes and #tags,Address,Receipt,Description,Category split,Money Out,Money In
```

## Security

- **Never commit** `config.yaml` or `.env` — both are gitignored.
- Use `config.example.yaml` and `.env.example` as templates.
- The Redis password is loaded exclusively from the `REDIS_PASSWORD` environment variable (set via `.env`).
