import test from 'node:test';
import assert from 'node:assert/strict';
import { deploymentLinks } from './deployment-links.mjs';
test('Pages outputs report production and the expected preview alias', () => {
  const deployment = 'https://abc123.alien-invaders.pages.dev';
  assert.equal(deploymentLinks('main', deployment, '').url, 'https://alien-invaders.pages.dev/');
  assert.equal(deploymentLinks('pr-10', deployment, 'https://pr-10.alien-invaders.pages.dev').url, 'https://pr-10.alien-invaders.pages.dev/');
});
test('missing, foreign or mismatched action outputs fail closed', () => {
  for (const url of ['', 'http://abc.alien-invaders.pages.dev', 'https://evil.example', 'https://abc.alien-invaders.pages.dev.evil.example', 'https://abc.alien-invaders.pages.dev/path']) {
    assert.throws(() => deploymentLinks('main', url, ''));
  }
  assert.throws(() => deploymentLinks('pr-10', 'https://abc.alien-invaders.pages.dev', 'https://pr-11.alien-invaders.pages.dev'));
});
