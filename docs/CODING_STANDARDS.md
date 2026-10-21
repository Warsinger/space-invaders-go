# Coding standards

- Use the Go version in go.mod and gofmt; group standard library, local and
  third-party imports. Prefer small functions with explicit errors.
- Keep `main` focused on startup/flags, `game` on lifecycle/rules, `components`
  on ECS behavior, `assets` on embedded resources and `scores` on persistence.
- Define Donburi data as `XData`, its component as `X`, and constructors as `NewX`.
  Pass worlds/entries explicitly; avoid new mutable globals for gameplay state.
- Preserve balance formulas and input semantics unless a reviewed change calls
  for them. Do not process input or simulation while paused/game over except the
  existing reset/quit controls. One projectile must not hit several targets.
- Keep platform filesystem and browser APIs behind build-constrained files.
  Missing/blocked browser storage must not terminate the runtime.
- Prefer tests of deterministic rules using small ECS fixtures. Use fake
  renderers for bounds rather than starting a graphics loop in unit tests.
- Return startup/constructor errors; log recoverable persistence errors. Never
  mask an earlier entity-update error with subsequent work.
- Keep gameplay assets unchanged during infrastructure work. This repository
  does not currently establish asset provenance or a license: do not invent
  credits/ownership claims. Ask the owner before adding such claims.
- Update README/spec/testing/browser documentation with behavior changes. Report
  which graphical/audio checks were actually performed.
