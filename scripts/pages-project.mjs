export async function pagesProject(fetcher, base, headers) {
  const project = 'alien-invaders';
  const response = await fetcher(`${base}/${project}`, { headers });
  if (!response.ok) {
    const detail = await response.json().catch(() => ({}));
    const codes = (detail.errors || []).map(error => error.code).filter(code => Number.isInteger(code)).join(",");
    throw new Error(`Cloudflare project lookup failed: HTTP ${response.status}; error codes: ${codes || "unavailable"}`);
  }
  const result = await response.json();
  if (!result.success || result.result?.name !== project || result.result.production_branch !== 'main' || result.result.source) {
    throw new Error('Expected the alien-invaders Direct Upload project with production branch main');
  }
  return result.result;
}
