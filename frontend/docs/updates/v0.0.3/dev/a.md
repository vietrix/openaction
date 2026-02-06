# v0.0.3 dev a

## Highlights
- Faster cold start for the CI controller.
- Improved dependency cache reuse across pipelines.
- Clearer step status messages in the build view.

## Changes
- Added fallback retry for artifact upload timeouts.
- Trimmed log noise for successful steps.
- Updated default runner labels for dev builds.

## Notes
If you are running custom runners, verify that the agent reports version metadata after upgrade.
