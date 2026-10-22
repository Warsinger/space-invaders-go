# Cloudflare Pages deployment

The owner has provisioned the free Direct Upload project `alien-invaders`, with
production branch `main` and hostname `alien-invaders.pages.dev`, and installed
repository Actions secrets `CLOUDFLARE_ACCOUNT_ID` / `CLOUDFLARE_API_TOKEN`.
The workflow validates that exact project, production branch, Direct Upload mode
and hostname. It never creates projects or changes account settings.

## Publish and preview

Pushes to main publish automatically after all native tests/vet/coverage/builds
and browser checks succeed. Deployment jobs are serialized and a superseded
main revision is rejected. The downloaded artifact's version.json must match
the workflow SHA. Publishing, project validation, preview comments and their tests
live in [cloudflare-pages-deploy](https://github.com/Warsinger/cloudflare-pages-deploy).
This repository retains its build/package checks and a small workflow caller.
Both the reusable workflow and its script checkout are pinned to the same full
shared commit; update both references together after reviewing shared changes.
The shared workflow wraps Cloudflare's official Wrangler action, pinned with
Wrangler 4.151.1, and receives only the explicit repository secrets below.

For a preview, manually dispatch **Go build** on the open PR branch with `deploy`
enabled and `preview_pr` set to its number. Only the current head of an open,
same-repository PR targeting main can publish. Ordinary PR events cannot deploy
and do not receive deployment secrets. The URL is
`https://pr-N.alien-invaders.pages.dev/`; production is
`https://alien-invaders.pages.dev/`. Manual dispatch without deploy only builds.

The action writes deployment links and the commit to the job summary. A preview
also creates or updates one github-actions[bot] PR comment with the preview and
build links. Production publication does not add a PR preview comment.

## Setup and verification

For future credential rotation/project recovery:

1. Create a free Pages Direct Upload project with production branch `main`.
2. Confirm the actual project name and assigned pages.dev hostname. Do not assume
   a name is available or reuse another app's project.
3. Add `CLOUDFLARE_ACCOUNT_ID` and `CLOUDFLARE_API_TOKEN` as repository Actions
   secrets. The token needs Pages edit access for the intended account; other
   repositories' secret values cannot be copied/read by the builder.
4. Use the existing publishing workflow. A changed account/project/hostname
   requires a reviewed configuration change; a missing project fails closed.
5. Test the preview and production origin: gameplay, sound after interaction,
   score persistence, version.json, missing assets and cache behavior.

CI uploads `space-invaders-web` containing `build/pages`.
`version.json` records the revision and content-addressed asset directory.
Old/new runtime, JavaScript and Wasm must stay in the same release directory;
do not mix files between builds. No Cloudflare CLI/login/secret is needed to
build, package, test or download this artifact.

After two known-good production deployments, document/test Cloudflare's
production rollback and roll-forward path with publishing idle. Preview
deployments are not production rollback targets; the next main deployment can
supersede a rollback. High scores belong to each site origin and are not copied
between previews and production.

The browser job installs the integrity-locked shared browser package with npm ci --ignore-scripts. Shared modules/styles are included in the app asset hash and Pages upload. No hosting target, credentials or production publication policy changes are required.
