# Adversarial report regressions

Use synthetic payloads only; never execute scripts. Test HTML breakout payloads in query, repository/dependency names, timezone, file names, advisories and scanner metadata. The generated HTML must have no script/active payload, no remote Mermaid import, a script-blocking CSP, inert graph source, and protected anchors. Test scheme/attribute policy with javascript, data, relative/opaque URLs, remote file shares, credentials, controls, backslashes and quote injection; preserve ordinary http/https and explicitly enabled absolute local links.

Test Mermaid labels cannot add directives, graph lines or Markdown fences, and IDs are generated rather than colliding on hostile names. Test local source URLs stay within the declared repository and preserve encoded file names. Test absent/no-match/vulnerability scan states remain distinct in HTML and Markdown even without advisories.
