# Testing and coverage

Run from the repository root:

```sh
go test ./...
go vet ./...
go build ./...
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func coverage.out
go tool cover -html coverage.out -o coverage.html
node --test scripts/*.test.mjs
node scripts/test-browser-storage.mjs
go run ./cmd/web
node scripts/package-pages.mjs
```

Linux graphics/audio packages and a display are needed for packages that import
Ebitengine; CI installs the dependencies and uses Xvfb. Unit tests do not start a
game loop. CI tests/vets/builds Linux amd64, Windows amd64, macOS arm64,
reports Linux coverage and uploads binaries plus a browser-ready artifact.

The focused tests cover high-score parsing/native persistence, score accumulation
and reset, bullet expiry, player bounds, barrier damage/count, alien movement,
and collision consumption/player death. Browser tests exercise raw/gzip runtime
loading, automatic ready/error feedback, real HTML references and all five
localStorage scenarios (normal, throwing getter/getItem/setItem, missing).

Manual checklist: open the website and confirm automatic gameplay/loading;
click to focus/unlock audio; move left/right; fire; pause/resume; restart; observe
alien/barrier hits and game over; reload and check high score; test a narrow
viewport and confirm stable board geometry; check keyboard focus, error feedback
and footer link. Q must not terminate the browser. Check native quit and relative
score files separately. A successful Wasm build is not a graphical support claim.

Remaining limits: no exhaustive browser/OS graphics matrix, touch controls,
automated sound verification or CI browser gameplay automation. Coverage reports
are informational, not a guarantee that rendering and randomness are tested.

Install browser assets with `npm ci --ignore-scripts` before browser tests/builds. Shared loader/status/frame behavior tests live in web-shared; this repo tests its own packaged references and options/storage where applicable.
