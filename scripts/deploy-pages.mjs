import { updatePreviewComment } from './preview-comment.mjs';
import { pagesProject } from './pages-project.mjs';
import { deploymentBranch } from './pages-target.mjs';
import fs from 'node:fs';

const project = 'alien-invaders';
const account = process.env.CLOUDFLARE_ACCOUNT_ID;
const token = process.env.CLOUDFLARE_API_TOKEN;
const repo = 'Warsinger/space-invaders-go';
const sha = process.env.GITHUB_SHA;
if (!account || !token) throw new Error('Set CLOUDFLARE_ACCOUNT_ID and CLOUDFLARE_API_TOKEN repository secrets');
if (process.env.GITHUB_REPOSITORY !== repo || !/^[a-f0-9]{40}$/.test(sha || '')) throw new Error('Deployment requires this repository workflow context');
if (!['push', 'workflow_dispatch'].includes(process.env.GITHUB_EVENT_NAME)) throw new Error('PR events cannot publish');

async function github(endpoint, options = {}) {
  const response = await fetch(`https://api.github.com/repos/${repo}/${endpoint}`, {
    method: options.method || 'GET',
    headers: { Authorization: `Bearer ${process.env.GITHUB_TOKEN}`, Accept: 'application/vnd.github+json', 'Content-Type': 'application/json' },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });
  if (!response.ok) throw new Error(`GitHub validation failed: HTTP ${response.status}`);
  return response.json();
}

const branch = await deploymentBranch(process.env, github);

const config = await pagesProject(fetch, `https://api.cloudflare.com/client/v4/accounts/${encodeURIComponent(account)}/pages/projects`,
  { Authorization: `Bearer ${token}` });
const hostname = config.subdomain;
if (hostname !== 'alien-invaders.pages.dev') throw new Error('Unexpected Pages hostname');
const version = JSON.parse(fs.readFileSync('build/pages/version.json', 'utf8'));
if (version.revision !== sha) throw new Error('Artifact revision does not match the checked workflow revision');
const mode = process.argv[2];
if (mode === 'validate') {
  if (!process.env.GITHUB_OUTPUT) throw new Error('Missing Actions output file');
  fs.appendFileSync(process.env.GITHUB_OUTPUT, `branch=${branch}\n`);
} else if (mode === 'report') {
  const { deploymentLinks } = await import('./deployment-links.mjs');
  const { url, deploymentURL } = deploymentLinks(branch, process.env.DEPLOYMENT_URL, process.env.DEPLOYMENT_ALIAS_URL);
  if (process.env.GITHUB_STEP_SUMMARY) fs.appendFileSync(process.env.GITHUB_STEP_SUMMARY,
    `## Space Invaders deployment\n\n[Open the game](${url}) · [This deployment](${deploymentURL})\n\nRevision: \`${sha}\`\n`);
  if (branch !== 'main') {
    await updatePreviewComment(github, process.env.PREVIEW_PR, url, sha,
      `https://github.com/${repo}/actions/runs/${process.env.GITHUB_RUN_ID}`);
  }
} else {
  throw new Error('Choose validate or report');
}
