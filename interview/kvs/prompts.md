M1. implement a simple Key-Value Storage. This storage will behave like an in-memory database where you can store and retrieve key-value pairs.

- Getting a non‑existent key must return an explicit error or distinct optional type. Do not silently return an empty or default value.
- ⚠️ **Use a map or dictionary for committed data.** Candidates who do not follow this approach almost certainly will not finish the required milestones.

M2. Add support for transactions. Extend your Key-Value Storage to support BEGIN, COMMIT, and ROLLBACK operations.

    - Maintain a separate transactional state—either a copy or a delta map.
    - Copy is simpler—allow the candidate to proceed with this approach, but discuss the trade-offs at some point.
    - Delta map is more efficient but complicates delete handling. Keep this in mind.
    - See **Gotcha 4** for more details: [**Full copy vs overlay approach** — Copying the entire committed map on BEGIN is simple but has O(n) space complexity where n is the number of keys. An overlay/delta map that only stores changed keys is more efficient with O(k) space where k is the number of changes. For large datasets with small transactions, this difference matters significantly. On the other hand, it will make handling the deletes and rollbacks more difficult. Watch if candidates consider this trade-off.](https://www.notion.so/Full-copy-vs-overlay-approach-Copying-the-entire-committed-map-on-BEGIN-is-simple-but-has-O-n-spa-28d5b28f309180ef905dcb43cd61e628?pvs=21)
    - Snapshots - make a copy of a map, and then apply the changes when commit, or discard the snapshot (similar to the copy)
- `PUT` operates on the transactional map when a transaction is active.
- `GET` reads uncommitted values first (transactional map), then falls back to committed map.
- `COMMIT` merges the transactional map into the committed map and clears transactional state.
- `ROLLBACK` discards the transactional map.

**Error cases:**

- `ROLLBACK` without active transaction → error
- `COMMIT` without active transaction → error

M3. Implement deletion
**Behavior:**
- Use a tombstone or delete marker so that a key deleted in a transaction is not visible via fallback to the committed map.
- `DELETE` on non-existent key → treat as no-op (preferred) or error. Be consistent and document your choice.
- `ROLLBACK` after `DELETE` restores committed value visibility.