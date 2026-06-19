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
  - [Precedence overview](#precedence-overview)
  - [Visual order](#visual-order)
  - [Group details](#group-details)
  - [Examples](#examples)
  - [How violations are detected](#how-violations-are-detected)
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

When writing or editing `className`, `class`, or `cn()` arguments, order utility classes in this sequence. Separate groups with a single space. Keep variant prefixes attached to each utility (e.g. `hover:bg-primary`, not `hover: bg-primary`).

Lower group numbers sort **earlier (left)**. Higher numbers sort **later (right)**. Within the same group, the original relative order is preserved (stable sort).

Unrecognized tokens are ignored for ordering — they do not trigger violations and stay in place during fix.

### Precedence overview

| # | Group | Examples |
|---|-------|----------|
| 1 | Position anchor | `relative`, `absolute`, `fixed`, `sticky`, `static` |
| 2 | Position offsets | `inset-*`, `top-*`, `right-*`, `bottom-*`, `left-*` |
| 3 | Self & group | `self-*`, `group`, `group/name` |
| 4 | Element | `shrink-*`, `grow-*`, `select-*`, `whitespace-*`, `compress-zero` |
| 5 | Margin & padding | `m-*`, `p-*`, and axis variants |
| 6 | Width & height | `w-*`, `h-*`, `min-*`, `max-*`, `size-*`, `aspect-*` |
| 7 | Display | `block`, `inline`, `hidden`, `visible`, `isolate` |
| 8 | Text size | `text-xs`, `text-sm`, `text-base`, `text-lg`, … `text-9xl` |
| 9 | Font | `font-medium`, `font-mono`, `font-*` |
| 10 | Text color | `text-red-500`, `text-muted-foreground`, `text-center` |
| 11 | Background & fill color | `bg-*`, `fill-*`, `stroke-*`, gradients, `opacity-*` |
| 12 | Variant modifiers | `hover:`, `focus:`, `disabled:`, `aria-*:`, `data-*:`, `dark:`, `md:` |
| 13 | Transition | `transition-*`, `duration-*`, `animate-*` |
| 14 | Border | `border`, `border-*`, `outline-*`, `ring-*`, `divide-*` |
| 15 | Rounding | `rounded-*` |
| 16 | Shadow | `shadow-*` |
| 17 | Truncate & overflow | `truncate`, `overflow-*`, `text-ellipsis` |
| 18 | Children (grid & flex) | `grid-*`, `flex`, `gap-*`, `items-*`, `justify-*`, etc. |
| 19 | End | `cursor-*`, `pointer-events-*`, `z-*` (always last) |

### Visual order

Groups are listed top → bottom (1 first, 19 last). In a `className` string, earlier groups sort to the left.

```mermaid
flowchart TB
    subgraph R01[ ]
        direction LR
        H1["1 Position anchor"] --> M1["relative<br/>absolute<br/>fixed<br/>sticky<br/>static"]
    end
    subgraph R02[ ]
        direction LR
        H2["2 Offsets"] --> M2["inset-*<br/>top-*<br/>right-*<br/>bottom-*<br/>left-*"]
    end
    subgraph R03[ ]
        direction LR
        H3["3 Self & group"] --> M3["self-*<br/>group<br/>group/name"]
    end
    subgraph R04[ ]
        direction LR
        H4["4 Element"] --> M4["shrink-*<br/>grow-*<br/>basis-*<br/>select-*<br/>whitespace-*<br/>compress-zero"]
    end
    subgraph R05[ ]
        direction LR
        H5["5 Spacing"] --> M5["m-* mx-* my-*<br/>mt-* mr-* mb-* ml-*<br/>p-* px-* py-*<br/>pt-* pr-* pb-* pl-*"]
    end
    subgraph R06[ ]
        direction LR
        H6["6 Size"] --> M6["w-* h-*<br/>min-w-* max-w-*<br/>min-h-* max-h-*<br/>size-* aspect-*"]
    end
    subgraph R07[ ]
        direction LR
        H7["7 Display"] --> M7["block<br/>inline<br/>hidden<br/>visible<br/>isolate"]
    end
    subgraph R08[ ]
        direction LR
        H8["8 Text size"] --> M8["text-xs<br/>text-sm<br/>text-base<br/>text-lg<br/>… text-9xl"]
    end
    subgraph R09[ ]
        direction LR
        H9["9 Font"] --> M9["font-medium<br/>font-mono<br/>font-*"]
    end
    subgraph R10[ ]
        direction LR
        H10["10 Text color"] --> M10["text-red-500<br/>text-muted-*<br/>text-center<br/>text-left"]
    end
    subgraph R11[ ]
        direction LR
        H11["11 Background"] --> M11["bg-* fill-*<br/>stroke-* from-*<br/>to-* via-*<br/>opacity-* accent-*"]
    end
    subgraph R12[ ]
        direction LR
        H12["12 Variants"] --> M12["hover:* focus:*<br/>disabled:* md:*<br/>dark:* aria-*:*<br/>data-*:* group-hover:*"]
    end
    subgraph R13[ ]
        direction LR
        H13["13 Transition"] --> M13["transition-*<br/>duration-*<br/>animate-*"]
    end
    subgraph R14[ ]
        direction LR
        H14["14 Border"] --> M14["border<br/>border-* outline-*<br/>ring-* divide-*"]
    end
    subgraph R15[ ]
        direction LR
        H15["15 Rounding"] --> M15["rounded-sm<br/>rounded-md<br/>rounded-full<br/>rounded-*"]
    end
    subgraph R16[ ]
        direction LR
        H16["16 Shadow"] --> M16["shadow-sm<br/>shadow-md<br/>shadow-lg<br/>shadow-*"]
    end
    subgraph R17[ ]
        direction LR
        H17["17 Overflow"] --> M17["truncate<br/>overflow-*<br/>text-ellipsis"]
    end
    subgraph R18[ ]
        direction LR
        H18["18 Flex/grid"] --> M18["grid flex<br/>gap-* items-*<br/>justify-* content-*<br/>order-* col-* row-*"]
    end
    subgraph R19[ ]
        direction LR
        H19["19 End"] --> M19["cursor-*<br/>pointer-events-*<br/>z-*"]
    end

    R01 --> R02 --> R03 --> R04 --> R05 --> R06 --> R07 --> R08 --> R09 --> R10 --> R11 --> R12 --> R13 --> R14 --> R15 --> R16 --> R17 --> R18 --> R19
```

ASCII summary of the same pipeline:

```
┌─────────────┬──────────────┬─────────────┬──────────┬───────────────┐
│ 1–7 Layout  │ 8–11 Type &  │ 12 Variants │ 13–17    │ 18 Children  │
│ position →  │ color        │ hover: md:  │ surface  │ grid flex    │
│ display     │ text bg      │ dark: …     │ border … │ gap items    │
└─────────────┴──────────────┴─────────────┴──────────┴──────┬───────┘
                                                             │
                                                    19 End: cursor z-*
```

### Group details

#### 1. Position anchor

Positioning mode comes first — before any offset values.

`relative` · `absolute` · `fixed` · `sticky` · `static`

Variant-prefixed anchors (e.g. `md:absolute`) are treated as **variant modifiers** (group 12).

#### 2. Position offsets

`inset-*` · `top-*` · `right-*` · `bottom-*` · `left-*`

#### 3. Self & group

`self-*` · `group` · `group/accordion-trigger`

Named groups like `group/accordion-trigger` are recognized as group utilities, not variant prefixes.

#### 4. Element

`shrink-*` · `grow-*` · `basis-*` · `select-*` · `whitespace-*` · `compress-zero`

#### 5. Margin & padding

**Unprefixed only** — responsive or state-prefixed spacing (e.g. `md:px-4`, `hover:p-2`) is treated as a **variant modifier** (group 12).

`m-*` · `mx-*` · `my-*` · `mt-*` · `mr-*` · `mb-*` · `ml-*` · `p-*` · `px-*` · `py-*` · …

#### 6. Width & height

`w-*` · `h-*` · `min-w-*` · `max-w-*` · `min-h-*` · `max-h-*` · `size-*` · `aspect-*`

#### 7. Display

`block` · `inline` · `hidden` · `visible` · `isolate`

Variant-prefixed display utilities (e.g. `md:hidden`) sort as **variant modifiers**.

#### 8. Text size

Font size tokens only — not color or alignment.

`text-xs` · `text-sm` · `text-base` · `text-lg` · `text-xl` · `text-2xl` … `text-9xl`

#### 9. Font

`font-medium` · `font-mono` · `font-*`

#### 10. Text color

Text color and text-related non-size utilities. **Unprefixed only.**

`text-red-500` · `text-muted-foreground` · `text-center` · `text-left`

Anything matching `text-*` that is **not** a text size token (group 8) belongs here.

#### 11. Background & fill color

Surface and decorative color. **Unprefixed only.**

`bg-*` · `fill-*` · `stroke-*` · `from-*` · `to-*` · `via-*` · `opacity-*` · `accent-*` · `caret-*` · `decoration-*`

#### 12. Variant modifiers

Any utility with a variant prefix not already placed in an earlier group:

`hover:*` · `focus:*` · `disabled:*` · `aria-*:*` · `data-*:*` · `dark:*` · `md:*` · `lg:*` · `group-hover:*`

- Prefixed margin/padding, display, position, border, rounding, shadow, and children utilities land here.
- Variant-prefixed transition utilities also land here (base `transition-*` without a prefix is group 13).
- The `:` prefix stays attached to the utility name.

#### 13. Transition

Motion and animation for the element itself (unprefixed).

`transition-*` · `duration-*` · `animate-*`

#### 14. Border

Borders, outlines, rings, and dividers. `rounded-*` is excluded (group 15).

`border` · `border-*` · `outline-*` · `ring-*` · `divide-*`

#### 15. Rounding

`rounded-*` · `rounded-sm` · `rounded-md` · `rounded-full`

#### 16. Shadow

`shadow-*` · `shadow-sm` · `shadow-md` · `shadow-lg`

#### 17. Truncate & overflow

Before children layout utilities.

`truncate` · `overflow-*` · `text-ellipsis`

#### 18. Children (grid & flex)

Layout that affects children. Within this group, `grid-*` sorts before `flex` when both are present.

`grid` · `grid-*` · `inline-grid` · `flex` · `inline-flex` · `gap-*` · `items-*` · `justify-*` · `content-*` · `place-*` · `order-*` · `col-*` · `row-*` · `space-x-*` · `space-y-*` · `list-*`

Variant-prefixed children utilities (e.g. `md:flex`) sort as **variant modifiers**.

#### 19. End

Interaction and stacking — **always last**.

`cursor-*` · `pointer-events-*` · `z-*`

### Examples

#### Position before offsets; text color before background

```tsx
// ❌ BAD
<div className="top-0 left-0 absolute bg-muted text-sm text-muted-foreground" />

// ✅ GOOD
<div className="absolute top-0 left-0 text-sm text-muted-foreground bg-muted" />
```

#### Font after text size; hover after base color

```tsx
// ❌ BAD
<button className="font-medium text-sm hover:bg-primary flex px-4 h-9 bg-background border rounded-md" />

// ✅ GOOD
<button className="px-4 h-9 text-sm font-medium text-foreground bg-background hover:bg-primary border rounded-md flex" />
```

#### Group near position; z-index last

```tsx
// ❌ BAD
<div className="px-3 py-2 group/accordion-trigger relative z-50 flex cursor-pointer" />

// ✅ GOOD
<div className="relative group/accordion-trigger px-3 py-2 text-sm font-medium text-muted-foreground bg-muted flex cursor-pointer z-50" />
```

#### Overflow before flex children

```tsx
// ❌ BAD
<div className="flex gap-2 overflow-auto truncate w-full" />

// ✅ GOOD
<div className="w-full overflow-auto truncate flex gap-2" />
```

### How violations are detected

`twaz` walks each class string left to right. Each token is classified into a group number. If a token's group number is **lower** than a token that appeared earlier, it is reported as a violation — it should have appeared earlier in the string.

Example: in `bg-muted text-sm absolute`, `text-sm` (group 8) appears after `bg-muted` (group 11), and `absolute` (group 1) appears after `bg-muted` — both are violations.

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
