# twaz

Check and automatically fix Tailwind CSS utility class order in JSX/TSX files.

`twaz` enforces a consistent, readable order for `className`, `class`, and `cn()` arguments. It scans source files, reports violations, and reorders classes in place. **Fix mode is enabled by default.**

Built with Go. Distributed on npm with a small Node.js launcher.

- **Repository:** [github.com/maxzz/tw-az](https://github.com/maxzz/tw-az)
- **npm package:** [npmjs.com/package/twaz](https://www.npmjs.com/package/twaz)

## Contents

- [Installation](#installation)
- [Usage](#usage)
  - [CLI](#cli)
  - [Go library](#go-library)
- [Options](#options)
- [What gets scanned](#what-gets-scanned)
- [Tailwind class order rules](#tailwind-class-order-rules)
- [Development](#development)
- [Publishing to npm](#publishing-to-npm)
- [License](#license)

## Installation

```bash
npm install -D twaz
```

Or run without installing:

```bash
npx twaz src
```

## Usage

### CLI

```bash
twaz [options] [paths...]
```

```bash
# Fix class order in the current directory (default)
twaz

# Fix a specific directory or file
twaz src
twaz src/components/Button.tsx

# Report violations only, without changing files
twaz --check src
```

On Windows, the CLI prints a colored report and waits for a key press before closing the window.

### Go library

```go
import "twaz/twaz"

violations := twaz.CheckClassString("bg-muted text-sm absolute")
sorted := twaz.SortClassString("bg-muted text-sm absolute top-0")
result := twaz.RunScan([]string{"src"}, twaz.ScanOptions{Fix: true})
```

## Options

| Option | Short | Default | Description |
|--------|-------|---------|-------------|
| `--fix` | `-f` | **on** | Reorder utility classes automatically in place. This is the default behavior; the flag is kept for explicit use in scripts. |
| `--check` | `-c` | off | Report class order violations only. Does not modify files. Same as `--no-fix`. |
| `--no-fix` | — | off | Alias for `--check`. |
| `--help` | `-h` | off | Show usage, options, and examples. |

**Arguments**

| Argument | Default | Description |
|----------|---------|-------------|
| `paths...` | current directory | One or more files or directories to scan. When a file path is given, its parent folder name is shown in the fix report. |

**Examples**

```bash
twaz src
twaz --check src
twaz src/App.tsx
```

## What gets scanned

By default, `twaz` scans `.tsx` and `.jsx` files. It looks for class strings in:

- `className="..."` / `className='...'`
- `className={`...`}` / `className={"..."}`
- `class="..."`
- `cn("...")` / `classNames("...")`

Directories `node_modules`, `dist`, and `.git` are skipped.

## Tailwind class order rules

See the [full class order reference](https://github.com/maxzz/tw-az) for details. In short, utilities are sorted into 19 groups from position anchor through children layout to cursor/z-index.

Unrecognized tokens are ignored for ordering (they do not trigger violations and stay in place during fix).

## Development

```bash
# Run CLI from source
npm run dev

# Go tests
npm test

# Build Windows binary (default)
npm run build

# Build all platform binaries (for npm publish)
npm run build:all
```

## Publishing to npm

1. Log in: `npm login`
2. Build all platform binaries: `npm run build:all`
3. Publish: `npm publish` (or `npm run to-npm`)

`prepublishOnly` runs `build:all` automatically before publish.

## License

MIT — see the `license` field in [package.json](./package.json).
