export function deploymentLinks(branch, deploymentURL, aliasURL) {
  const hostname = 'alien-invaders.pages.dev';
  const expected = branch === 'main' ? hostname : branch + '.' + hostname;
  const deployment = new URL(deploymentURL);
  if (deployment.protocol !== 'https:' || deployment.username || deployment.password ||
      !/^[a-z0-9-]+\.alien-invaders\.pages\.dev$/.test(deployment.hostname) ||
      deployment.port || deployment.pathname !== '/' || deployment.search || deployment.hash) {
    throw new Error('Unexpected deployment URL');
  }
  if (aliasURL) {
    const alias = new URL(aliasURL);
    if (alias.origin !== 'https://' + expected || alias.pathname !== '/' || alias.search || alias.hash || alias.username || alias.password) {
      throw new Error('Unexpected deployment alias');
    }
  }
  return { url: 'https://' + expected + '/', deploymentURL: deployment.href };
}
