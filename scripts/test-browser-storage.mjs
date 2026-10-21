// Exercise the actual browser storage implementation under Go/Wasm without
// importing the graphics engine, which requires a browser DOM.
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { spawnSync } from 'node:child_process';
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
function run(command, args, options = {}) {
  const result = spawnSync(command, args, { encoding: 'utf8', windowsHide: true, ...options });
  if (result.error || result.status !== 0) {
    throw new Error(result.error?.message || result.stderr || result.stdout || `${command} failed`);
  }
  return result.stdout.trim();
}

if (process.argv[2] === '--case') {
  const [, , , mode, wasm, runtime] = process.argv;
  let stored = null;
  const storage = {
    getItem() {
      if (mode === 'getItem-throws') throw new Error('getItem blocked');
      return stored;
    },
    setItem(key, value) {
      if (mode === 'setItem-throws') throw new Error('quota exceeded');
      stored = value;
    },
  };
  Object.defineProperty(globalThis, 'localStorage', {
    get() {
      if (mode === 'getter-throws') throw new DOMException('Storage blocked', 'SecurityError');
      if (mode === 'missing-storage') return undefined;
      return storage;
    },
  });
  createRequire(import.meta.url)(runtime);
  const go = new Go();
  let exitCode = 0;
  go.exit = code => { exitCode = code; };
  const { instance } = await WebAssembly.instantiate(fs.readFileSync(wasm), go.importObject);
  await go.run(instance);
  if (exitCode !== 0) throw new Error(`${mode}: Go exited with ${exitCode}`);
  if (mode === 'normal' && stored !== '42') throw new Error('Normal persistence failed');
  console.log(`PASS ${mode}`);
} else {
  const temp = fs.mkdtempSync(path.join(os.tmpdir(), 'invaders-storage-test-'));
  try {
    const source = fs.readFileSync(path.join(root, 'scores/storage_js.go'), 'utf8');
    // Only the package name changes; the storage functions are tested verbatim.
    fs.writeFileSync(path.join(temp, 'storage_js.go'), source.replace('package score', 'package main'));
    fs.writeFileSync(path.join(temp, 'score.go'), fs.readFileSync(path.join(root, 'scores/score.go'), 'utf8').replace('package score', 'package main'));
    fs.writeFileSync(path.join(temp, 'main.go'), `package main
func main() {
  Load()
  if err := Save(42); err != nil { panic(err) }
  if Load() != 42 { panic("score round trip failed") }
}
`);
    const wasm = path.join(temp, 'storage.wasm');
    const env = { ...process.env, GOOS: 'js', GOARCH: 'wasm', GOWORK: 'off' };
    run('go', ['build', '-o', wasm, path.join(temp, 'main.go'), path.join(temp, 'storage_js.go'), path.join(temp, 'score.go')], { env });
    const goRoot = run('go', ['env', 'GOROOT']);
    const runtime = path.join(goRoot, 'lib', 'wasm', 'wasm_exec.js');
    for (const mode of ['normal', 'getter-throws', 'getItem-throws', 'setItem-throws', 'missing-storage']) {
      console.log(run(process.execPath, [fileURLToPath(import.meta.url), '--case', mode, wasm, runtime]));
    }
  } finally {
    if (path.dirname(path.resolve(temp)) !== path.resolve(os.tmpdir()) || !path.basename(temp).startsWith('invaders-storage-test-')) {
      throw new Error('Refusing to remove unexpected temporary directory');
    }
    fs.rmSync(temp, { recursive: true, force: true });
  }
}
