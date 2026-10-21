# Browser build

Use Go 1.27.2 (go.mod) and run `go run ./cmd/web -serve`. Open
http://localhost:8080; use `-port 8085` if another app is running. The server binds
to localhost and disables caching. Do not open HTML via `file://`.

`go run ./cmd/web` builds `build/web` without serving. It compiles the game for
js/wasm and copies wasm_exec.js from the exact same Go installation. No Go
cross-compilation environment variables need to be set. The old root HTML and
wasmpackage.sh prototype are replaced by this portable command.

The game loads automatically, with instructions and the shared project footer
below it. Keyboard focus moves to the game canvas once it is ready. Click inside
if focus is lost; sound still requires browser interaction. Left/Right move,
Space fires, P pauses and R restarts. Q is native-only. Space/arrows do not scroll
while gameplay has focus. The 600×800 logical board scales with its frame.

High scores save under localStorage key `space-invaders.highscore.v1`. Storage is
origin-scoped (including port), separate from native files and other apps. A
blocked getter, getItem/setItem failure or missing storage switches to an
in-memory high score for this session. Clearing site data removes stored scores.
No score data is uploaded. High-score improvements save without waiting for quit.

For future static hosting, run `node scripts/package-pages.mjs` (Node 24) to
produce `build/pages`. Use a fresh destination, e.g. `node scripts/package-pages.mjs
build/pages-next`, when repeating packaging. The packager gzip-compresses Wasm
to meet the Pages 25 MiB file limit, uses content-addressed asset paths, marks
assets immutable and root/version metadata revalidatable, and provides a real
404 response page. The loader handles raw local Wasm and explicitly packaged
gzip payloads. No service worker or backend is needed. See DEPLOYMENT.md.
