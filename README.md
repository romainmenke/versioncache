# VersionCache

Dumb, no frills in memory cache in Go.

- no expiration of cached objects
- "versioning" of cached objects

A call to `Version` will clear the entire cache.

This is intended to cache things that have a clear action that invalidates the entire cache.
