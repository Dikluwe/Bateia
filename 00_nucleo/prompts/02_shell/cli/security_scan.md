# Optional security scan CLI

RunSecurityScanCommand(args, stdout, domain.BinaryVulnerabilityScanner) accepts --binary ABSOLUTE_REGULAR_FILE, optional --output JSON and positive --timeout <=5m. No source-mode, arguments or arbitrary database/tool configuration is exposed. Deadline context reaches the injected scanner; nil/unavailable scanner produces an honest NOT_SCANNED report. Invalid requests are errors, operational scanner failures remain structured data.

JSON on stdout or bounded serialized output (16 MiB), atomic file replacement using the existing protected atomic helper. Output must never overwrite artifact: compare absolute spelling and file identity for aliases, symlinks and hardlinks, check again before write. Keep initial artifact identity for the atomic writer. No integration with automatic decision gates. Wiring is a separate responsibility.
