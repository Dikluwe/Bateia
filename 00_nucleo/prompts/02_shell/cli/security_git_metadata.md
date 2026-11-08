# Confined advisory Git remote metadata

Local GitHub identity hints must read .git/config (or relative gitdir/commondir resolution) through one retained os.Root descriptor. Never return an absolute pathname and reopen it unconstrained. Relative pointers may normalize within the declared worktree, but absolute/external/escaping/control-containing pointers fail closed. Reject symlinks at metadata path components and final descriptors; os.Root enforces confinement even under path races. Open O_NONBLOCK|O_NOFOLLOW; descriptors must be bounded regular single-link files (nlink=1) to refuse hardlinks/FIFOs/devices. Pointer files <=4096 bytes, config <=16MiB with read growth guard.

Only missing commondir allows reading gitdir/config. Present but unsafe/unreadable commondir must not silently fall back to another config. Unavailable/unsafe metadata yields nil identity hints and no remote remapping. Valid normal .git/config and relative internal gitdir/commondir remain supported. Nonlocal file authorities/query/fragment are not local worktrees.

Tradeoff: external administration used by linked Git worktrees is deliberately unsupported; remote identity remains unknown instead of reading another root. Hardlinks and symlinks inside the tree also lose hints under the conservative policy. Metadata is a local declaration, never authentication of a remote origin or a signed commit. This does not change aggregate tree-hash budgets or historical certificates.

After bounded reading, f.Stat must still show the same regular single-link descriptor, unchanged initial size/modtime and size equal to bytes read. Reject mutation or newly acquired hardlinks; this is consistency checking, not an atomic transaction or cryptographic proof of Git identity.

0097: readLocalGitConfigLimit accepts a positive caller-owned limit no larger than16MiB. Capture admission uses64KiB/config only with target; freezes normalized IDs and never rereads after collection. Default diagnostic reader retains16MiB. Hints do not authenticate origin or enter file digest; their separate fingerprint identifies the run.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
