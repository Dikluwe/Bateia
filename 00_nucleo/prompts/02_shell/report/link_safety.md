# Report links: context and scheme safety

L2 pure formatting only. Untrusted repository and advisory links are parsed with net/url and validated before rendering. Web links accept only http and https with a host, no credentials or opaque URLs. Local links are an explicit separate option: file with an absolute path, empty/localhost authority, no query or fragment. Reject raw or percent-encoded controls and backslashes. This validation does not authorize reading, execution or guarantee that a web target is trustworthy.

Escape text and href attributes independently with html.EscapeString; rejected URLs become inert labels. All target=_blank links use rel=noopener noreferrer. sourceFileURL must reject paths escaping the repository and percent-encode local file names; Markdown's raw HTML details use the same link helper.
