# Optional scan CLI regressions

Inject a fake BinaryVulnerabilityScanner and use temporary artifact bytes. Verify stdout JSON, atomic optional file output, exact artifact path forwarding and preservation. Reject missing/relative/nonregular inputs, source flag, positional trailing arguments, zero/excessive timeout, output identical to binary including symlink and hardlink aliases before invoking runner. Nil scanner yields NOT_SCANNED/incomplete/advisory-only, never clean coverage. No network and no artifact execution in these tests.
