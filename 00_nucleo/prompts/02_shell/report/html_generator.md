# Specification: HTML Report Generator (L2 Shell)

## Purpose
Render domain SearchReport as a standalone, modern HTML page with Dark Mode styles, responsive tables, cards, and expandable details.

## Rules
- L2 package (`02_shell/report`).
- Pure formatting of `domain.SearchReport`.
- Accepts `*domain.SearchReport` and returns `string` (HTML).
- Escape every untrusted text in its HTML context, including graph source and timezone names. Validate URL schemes before escaping href attributes; protect target=_blank with noopener noreferrer. Unsafe links remain plain labels.
- Fully offline output: no JavaScript or remote CDN. Show Mermaid source as escaped text instead of executing a renderer. This intentionally trades visual rendering for offline determinism and removes third-party script trust.
- Defensive meta CSP: default/script/network/object/frame sources denied, no base or forms; inline CSS is allowed only to preserve the existing static presentation. CSP supplements escaping; it does not replace it or control manual navigation to hyperlinks.
- Always display the scope/limits of security scans. Empty advisories do not mean a secure repository.
- This hardens the HTML artifact, not arbitrary external Markdown renderers. The full Markdown report is evidence from untrusted repositories: renderers must disable raw HTML and run Mermaid in strict security mode. Its other dynamic fields are not a general-purpose Markdown sanitizer.
