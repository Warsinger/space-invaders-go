import {startGoGame} from './shared/go-runtime.mjs';
startGoGame({
  wasmURL: new URL('game.wasm', import.meta.url),
  readyCallback: 'spaceInvadersReady',
  stoppedMessage: 'Game stopped. Reload the page to restart.',
  preventKeys: ['Space', 'Enter', 'Escape', 'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight']
});
