# Heka Deadlock Example

These plugins when run together with the example.toml cause heka to freeze.

This is caused by the router freezing on trying to send a pack to filter that
is requesting a large number of pipeline packs (more than available in the pool).

Because the router is stuck, no other plugins are available to free up packs that the greedy filter is injecting.