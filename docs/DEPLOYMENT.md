# Cloudflare preparation — publishing deferred

The static artifact is ready for later deployment, but CI intentionally has no
deployment job and this work does not create a Cloudflare project, upload code,
or publish a site. Follow issue #4 after the owner provisions the application.

Owner setup checklist:

1. Create a free Pages Direct Upload project with production branch `main`.
2. Confirm the actual project name and assigned pages.dev hostname. Do not assume
   a name is available or reuse another app's project.
3. Add `CLOUDFLARE_ACCOUNT_ID` and `CLOUDFLARE_API_TOKEN` as repository Actions
   secrets. The token needs Pages edit access for the intended account; other
   repositories' secret values cannot be copied/read by the builder.
4. Authorize deployment. Add a pinned Wrangler deployment job with target/name/
   hostname and revision validation, serialized publication, main publishing
   after successful tests/builds, and explicitly requested same-repository PR
   previews. Post preview links in a bot-owned PR comment and build summary,
   following Tower Defense's approach. Do not enable triggers before this step.
5. Test the preview and production origin: gameplay, sound after interaction,
   score persistence, version.json, missing assets and cache behavior.

Current preparation: CI uploads `space-invaders-web` containing `build/pages`.
`version.json` records the revision and content-addressed asset directory.
Old/new runtime, JavaScript and Wasm must stay in the same release directory;
do not mix files between builds. No Cloudflare CLI/login/secret is needed to
build, package, test or download this artifact.

After two known-good production deployments, document/test Cloudflare's
production rollback and roll-forward path with publishing idle. Preview
deployments are not production rollback targets; the next main deployment can
supersede a rollback. High scores belong to each site origin and are not copied
between previews and production.
