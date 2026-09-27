# Present security analysis without false assurance

Pure L2 formatter for domain.SecurityScanReport. Always show status, scanner, scope, database date and limitations, even when no advisory exists. An absent record means not evaluated. NOT_SCANNED is not assessed, SCANNED_NO_MATCH is no match within the declared scan/base (not absence of vulnerabilities), VULNERABILITY_FOUND is a scanner signal requiring validation of impact. Unknown statuses are not assessed.

HTML escapes all metadata as text. Markdown suppresses raw HTML and active Markdown link/fence syntax in scanner metadata. No renderer upgrades heuristic signals into proof of software safety.
