# Project workflow

- Use a dedicated `codex/` issue branch and a PR targeting `main`; do not push implementation directly to main.
- Link completed issues in GitHub Development using `Closes #N` and verify `closingIssuesReferences` before handoff.
- Assign every issue/PR to `Warsinger`; request `Warsinger` review on new PRs. Use appropriate documentation, testing, dependencies, platform and build labels. Apply priority only when the owner sets it.
- Discuss gameplay changes before implementing them. Record agreed scope/design in the issue. Preserve existing controls and balance unless a fix is required and documented.
- Run meaningful checks and inspect CI results. Never merge without user authorization.
- Cloudflare provisioning and publishing are deferred under issue #4. Do not create a project or upload/deploy until the owner confirms the project name/hostname and authorizes publishing.

# GitHub identity

Use Warsinger Code Builder for GitHub interactions, not connected GitHub tools or default gh authentication. Read `C:\Users\eelwe\.codex\github-identities\README.md` first. On this machine the shared helper's allowlist does not yet include this repository; use the compatibility copy outside the repository:

```powershell
node D:\eelwe\dev\codex\github-role-space-invaders.mjs run builder --repo Warsinger/space-invaders-go -- gh <arguments>
```

The copy uses the existing identity home and only extends the repository allowlist. Pass the same repository to gh options. Keep credentials outside the repository; never print/copy tokens or keys, persist temporary authentication, or put tokens in remote URLs. API identity and Git push/commit identity are separate. Publish commits with the intended builder attribution.

# Development

- Read `docs/PROJECT_SPEC.md`, `docs/CODING_STANDARDS.md` and `docs/TESTING.md` before relevant changes.
- Target the Go version in go.mod. Format changed Go files, then run applicable tests, vet and builds.
- Keep assets embedded and rendering/input separate from deterministic rules where practical.
- Browser high scores use their own namespaced localStorage with a session fallback; native scores remain in `score/highscore.txt`.
- Keep README and browser/deployment docs aligned with implemented behavior. Report manual graphics/audio validation limits honestly.
