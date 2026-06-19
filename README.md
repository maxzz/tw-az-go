# twaz

Check and automatically fix Tailwind CSS utility class order in JSX/TSX files.

`twaz` enforces a consistent, readable order for `className`, `class`, and `cn()` arguments. It scans source files, reports violations, and can reorder classes in place with `--fix`.

This package ships a native Go binary for fast CLI usage. The Go library in `twaz/` can also be imported directly from Go code.

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
# Check class order (default: current directory)
twaz

# Check a specific directory or file
twaz src
twaz src/components/Button.tsx

# Automatically fix class order in place
twaz --fix src
```

### Go library

```go
import "twaz/twaz"

violations := twaz.CheckClassString("bg-muted text-sm absolute")
sorted := twaz.SortClassString("bg-muted text-sm absolute top-0")
result := twaz.RunScan([]string{"src"}, twaz.ScanOptions{}, os.Stdout)
```

### Development

```bash
# Run CLI from source
npm run dev
npm run fix

# Go tests
npm test

# Build Windows binary (default)
npm run build:go

# Build all platform binaries (for npm publish)
npm run build:all

# Bump patch version, sync version.go, and rebuild Windows binary
npm run build
```

## Publishing to npm

1. Log in: `npm login`
2. Build all platform binaries: `npm run build:all`
3. Publish: `npm publish` (or `npm run to-npm`)

`prepublishOnly` runs `build:all` automatically before publish.

## What gets scanned

By default, `twaz` scans `.tsx` and `.jsx` files. It looks for class strings in:

- `className="..."` / `className='...'`
- `className={`...`}` / `className={"..."}`
- `class="..."`
- `cn("...")` / `classNames("...")`

Directories `node_modules`, `dist`, and `.git` are skipped.

## Tailwind class order rules

See the [full class order reference](https://github.com/maxzz/tw-az) for details. In short, utilities are sorted into 19 groups from position anchor through children layout to cursor/z-index.

Unrecognized tokens are ignored for ordering (they do not trigger violations and stay in place during `--fix`).

## In conclusion

twaz is an acronym standing for "Tailwind from A to Z" (referring not to alphabetical order, but to the order of values in the CSS layout model).

Happy coding!
