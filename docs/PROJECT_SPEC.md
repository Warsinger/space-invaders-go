# Space Invaders project specification

Space Invaders is a single-player arcade game written in Go using Ebitengine for
graphics/input/audio and Donburi for ECS state. `main.go` parses `-level` (default
1, minimum 1), creates and initializes the game, then runs the Ebitengine loop.

The board is 600×800. A player near the bottom moves horizontally with arrow keys
and fires one shot for each Space press. A hit kills the player; R resets the run
to the configured starting level. P toggles pause. Q saves and exits on native
builds and is ignored in browsers.

Aliens begin in 4 rows × 8 columns, alternating two embedded sprites. Levels
above 10 add one row. Horizontal speed is `level/3 + 1`; aliens reverse at their
75-pixel horizontal travel limit and descend 10 pixels. Alien shots are random.
Each alien is worth `10 + level*5`. Destroying all aliens advances the level,
clears remaining bullets/barriers and creates a fresh formation and barriers.
An alien reaching the player's vertical position also ends the game.

Barriers consist of four destructible lines. Their count is
`max(2, (level-1)%6+1)`. Each bullet consumes at most one target/hit; either side's
bullets damage barriers. Player bullets kill aliens; alien bullets kill the
player. Bullets expire above/below the board. Existing bullet horizontal motion
is not applied, even though alien bullets carry an x velocity; changing that
behavior is future gameplay work.

The current score resets with each run; the high score survives restarts. Native
storage is relative to the working directory at `score/highscore.txt`; missing
directories are created on save. Browser storage uses
`space-invaders.highscore.v1` on the current origin, with a session fallback when
localStorage is blocked or fails. High-score improvements save immediately;
restart, game-over and native quit also attempt saving. Save failures warn rather
than stopping gameplay. Browser/native stores are independent.

The browser uses the same gameplay/assets in a fixed logical board scaled by the
frame. The website loads automatically with loading/error feedback, instructions
and the shared project footer. Sound requires browser interaction. Keyboard and
desktop browsers are the initial target; touch controls, multiplayer and PWA
installation are not implemented. Cloudflare publishing uses the owner-provisioned
alien-invaders.pages.dev project; see DEPLOYMENT.md.

Architecture: `game` owns lifecycle, collisions, scoring and overlays;
`components` owns movement, bullets, barriers and sprite renderers; `assets`
embeds images/sounds/fonts; `scores` isolates storage; `platform` isolates browser
hooks; `cmd/web` builds and serves the browser bundle. There is no application
backend or remote score collection.
