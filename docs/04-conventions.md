# Conventions

Every package follows the rules below. They are non-negotiable. A contribution that breaks one without justification is rejected.

## Naming

### Packages

- Lowercase, single word, no underscores or camelCase.
- Singular when the package is about a single concept (`shell`).
- Plural when the package yields collections or grouped resources (`paths`, `files`).
- Never prefix with `onyx` — the import path already provides that context.

### Items

- **Be verbose. Be explicit. Never abbreviate.**
  - `ResolveExecutable`, never `RslvExec`.
  - `RevealInFileManager`, never `ShowFile`.
- Functions are `PascalCase` (exported) or `camelCase` (unexported), start with a verb in the imperative mood.
  - `OpenPath`, `Set`, `Resolve`.
- Types are `PascalCase`. `Resolver`, `Platform`, `AppPaths`.
- Constructors return a configured value, named after the value or the domain.
  - `NewResolver()`, `paths.For("hyperion")`.
- Predicate methods start with `Is`, `Has`, or `Can`.
  - `IsDarwin`, `HasExtension`, `CanWrite`.

### Errors

- One or more sentinel errors per package, declared at the top of the file.
- Sentinels follow the `Err<Reason>` convention: `ErrNotFound`, `ErrEmptyService`, `ErrUnavailable`.
- `errors.Is`-able. Wrap backend failures via `fmt.Errorf("...: %w", err)` — never replace the cause.
- Distinguish input validation, not-found, unavailable backend, and wrapped backend errors as separate sentinels.
- Never `panic` for recoverable failures. Never `log.Fatal` in library code.

## Function design

- **Single responsibility.** A function does one thing the name describes.
- **`(T, error)` over `bool`-second-return for fallible operations.** Reserve `(T, bool)` for genuinely optional reads (map lookups, etc.).
- **No boolean flag parameters.** Branching on flags means two functions wearing one name — split them.
- **No `interface{}` / `any` in public signatures unless dispatch is genuinely dynamic.** Prefer typed parameters.
- **Builders for 3+ parameters.** A function that needs many options exposes a builder so call sites stay readable.
- **`string` for paths is acceptable in Go** (the stdlib uses `string` everywhere) but document expectations — absolute vs relative, trailing slash, etc.

## Comments

- **Code is the documentation.** If the name does not explain the function, the name is wrong.
- **No inline `//` narration.** Comments explaining *what* the next line does are forbidden.
- **`//` doc comments on every exported item** — one or two sentences describing intent, not implementation. Standard Godoc rules.
- **`// TODO(@handle):` is allowed**, one line, with the author's GitHub handle.

## Documentation

- Every package has a markdown file in `docs/2N-<package>.md`.
- The file covers: purpose, public API table, examples, error catalog, dependencies, related packages, navigation footer.
- `README.md` is the adoption hook — long-form text belongs in `docs/`.
- Each package also has a `doc.go` with a `// Package <name> ...` block describing the package surface.

## SOLID, DRY, KISS

- **Single responsibility** — one package = one concern. No "utils" grab bag.
- **Open/closed** — extend by adding new types or implementations, not by editing existing ones.
- **Liskov** — interfaces are tiny, behavior is consistent across implementations.
- **Interface segregation** — prefer many small interfaces. A consumer should depend only on what it uses.
- **Dependency inversion** — packages depend on contracts in `onyx`, not on transitive third-party libraries.
- **DRY** — if the same OS switch appears twice, it lives in `osinfo` or a shared helper.
- **KISS** — the simplest correct API wins. Reflection, init-time side effects, and channels are last resorts.

## Testing

- Every public function has a test.
- Tests are `TestXxx`, scenario-first: `TestResolveFailsWhenNothingMatches`.
- Use `t.TempDir()` and the real filesystem; no mocking stdlib.
- OS-facing tests skip cleanly when no backend is reachable (return early on `Unavailable`). The suite stays green on minimal CI.

## Lint gate

Before opening a PR, all three pass:

```sh
gofmt -l . | grep -q '^' && exit 1 || true
go vet ./...
go test ./...
```

CI enforces the same gate on every push and pull request.

## Module layout

```
onyx/
├── go.mod
├── README.md
├── <package>/
│   ├── <package>.go         public API
│   ├── <package>_test.go    tests
│   └── doc.go               package-level Godoc block
└── docs/
    ├── 00-index.md          TOC
    ├── 01-installation.md   …
    └── 2N-<package>.md      per-package reference
```

No `util.go`, `helpers.go`, or `common.go`.

---

Navigation: [← Architecture](03-architecture.md) · **Conventions** · [Release flow →](05-release-flow.md)
