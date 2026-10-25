# Sakura deployment fork

Based on upstream v0.7.1, commit 619082a825fbb28a07950b8fe75c1b69282578d0.

This fork requires current read/create permissions and channel membership before linking a Todo to an existing post. All later Bot thread replies recheck current membership and permissions, including after a user has left a channel. Unlinked personal tasks remain available. The HTTP handler validates links before list mutation.

The server authorization regression tests cover missing/deleted posts, archived channels, membership, read/create permissions, early rejection, and revocation before a later Bot reply. No Mattermost license behavior is changed.

The Linux amd64 deployment bundle uses a locally built server binary and the source-built webapp with mobile touch targets, profile-language menu labels, and a main-menu Tasks entry. The plugin ID stays com.mattermost.plugin-todo and the distribution version is 0.7.1-sakura.1. Run the standard manifest generator before Go tests/builds. The frontend replaces node-sass with pinned Dart Sass for modern Node builds; Webpack 4 still needs NODE_OPTIONS=--openssl-legacy-provider during the build.

A previously observed TaskItems CSS-selector exception came from Mattermost core LCP telemetry, not this Todo UI. The companion Mattermost mobile.3.2 fork handles malformed browser metric selectors without interrupting application callbacks.

Upstream LICENSE.txt remains in force. This fork is tested with Mattermost Team Edition 11.7.14 and the Sakura mobile UI. Physical iPhone acceptance is recorded separately from browser emulation.

The main menu and post menu follow the current profile language. Mobile inputs use 16px text and primary buttons have 44px targets. Japanese IME confirmation is kept in add/edit fields before the following normal Enter submits. The test bootstrap is local; the unused archived full-webapp Git development dependency is removed.
