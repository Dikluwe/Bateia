# Confined Git metadata corpus

Synthetic files only, no Git commands/network/secrets. Preserve normal .git/config, relative internal gitdir and commondir with an internal parent segment. Refuse absolute/relative-escaping gitdir and commondir, .git/config/parent symlinks, external hardlinks (including pointer files), FIFO config, oversized pointer, and remote file authorities. Unsafe common pointer must suppress remapping rather than use the otherwise valid direct config.

Rename the filesystem pathname after opening os.Root and replace it with different config; a read through the retained descriptor must still return original root bytes. Treat hints as declarations, not source authentication.
