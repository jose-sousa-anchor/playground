# Asset Registry Tagging - Diagrams Guide

## Diagrams Created for Presentation

### 1. **asset_tagging_flow.puml** ⭐ PRIMARY TAGGING DIAGRAM
**Purpose**: Shows the complete tagging flow from definition to usage

**4 Phases:**
1. **Tag Definition** - Engineer creates PR with new tag (REBASING, STABLE, etc.)
2. **Apply Tags to Assets** - Ops/Admin uses Retool to tag assets
3. **Client Cache Population** - Client loads tags every 10s via ListTagsWithAssetTypeIDs
4. **Service Usage** - Services query assets by tags using GetByTags()

**Use in presentation**: Primary diagram to explain tagging system flow

---

### 2. **asset_tagging_architecture.puml**
**Purpose**: Shows the architectural components and data flow

**Key components:**
- Tag Definition Layer (code-based tags)
- Tag Management Layer (Retool UI + Service)
- Database (supported_asset + asset_tag tables)
- gRPC APIs (ListTags, AddTags, RemoveTags, ListTagsWithAssetTypeIDs)
- Client Layer (AssetTypes Client with dual caches)
- Consumer Services

**Use in presentation**: Technical deep-dive on architecture

---

### 3. **before_after_tagging.puml** ⭐ VALUE PROPOSITION DIAGRAM
**Purpose**: Side-by-side comparison showing the problem and solution

**WITHOUT TAGGING (Current):**
```go
// Hardcoded asset IDs
if assetID == "ethereum:steth" ||
   assetID == "ethereum:reth" {
    // Rebasing logic
}
```
Problems: Code changes, easy to miss assets, duplicated logic

**WITH TAGGING (Future):**
```go
// Query by tag
rebasingAssets := client.GetByTags("REBASING")
// Apply logic to all
```
Benefits: No code changes, automatic discovery, single source of truth

**Use in presentation**: Show this to demonstrate value and motivate the project

---

### 4. **integrated_asset_and_tagging_flow.puml**
**Purpose**: Shows how tagging integrates with existing dynamic asset system

**Complete 5-step flow:**
1. Add New Asset (existing flow: GraphQL → HSM → Special Org)
2. Asset Registry Sync (existing: Cron → Database)
3. **NEW: Apply Tags** (Retool → asset_tag table)
4. **Enhanced: Client Refresh** (loads both assets AND tags)
5. **Enhanced: Service Usage** (query by ID or by tags)

**Use in presentation**: Show how tagging builds on existing infrastructure

---

## Recommended Presentation Flow

### Opening: Set the Stage
1. Show **dynamic_asset_simple_flow.puml** (created earlier)
   - "Here's how we dynamically add assets today"
   - Emphasize the 4-phase flow

### Problem Statement
2. Show **before_after_tagging.puml** (WITHOUT side)
   - "But we have a problem: no way to categorize or discover assets"
   - Show hardcoded asset IDs example
   - List the problems (code changes, duplication, fragility)

### Solution Overview
3. Show **before_after_tagging.puml** (WITH side)
   - "Tagging solves this with a metadata layer"
   - Show GetByTags() example
   - Highlight benefits

### Solution Detail
4. Show **asset_tagging_flow.puml**
   - Walk through the 4 phases
   - Explain tag governance (PR required)
   - Show how ops can apply tags via Retool
   - Explain automatic cache propagation

### Technical Architecture
5. Show **asset_tagging_architecture.puml** (optional, for technical audience)
   - Component breakdown
   - Database schema
   - gRPC APIs
   - Cache structure

### Integration Story
6. Show **integrated_asset_and_tagging_flow.puml**
   - "Tagging builds on our existing dynamic asset infrastructure"
   - Show how steps 1-2 are unchanged
   - Highlight new steps 3-5
   - Emphasize backward compatibility

### Progress Update
7. Show Linear tickets summary (asset_registry_tagging_tickets.md)
   - Current status: 8 completed, 5 in progress/review
   - Key accomplishments
   - Timeline to completion

---

## Key Talking Points for Each Diagram

### For asset_tagging_flow.puml

**Phase 1: Tag Definition**
- "Tags are defined in code via PR - this ensures governance"
- "Each tag must have clear description and behavioral implications"
- "Prevents tag sprawl and maintains quality"

**Phase 2: Apply Tags**
- "Once tags are defined, ops can apply them via Retool"
- "No code changes needed to tag an asset"
- "Stored in asset_tag table with audit trail"

**Phase 3: Cache Population**
- "Client automatically picks up tags every 10 seconds"
- "Same refresh mechanism as assets - just enhanced"
- "Tags cached in memory for fast lookup"

**Phase 4: Service Usage**
- "Services call GetByTags('REBASING') and get all rebasing assets"
- "Works for current AND future assets"
- "No code changes when new rebasing token is added"

---

### For before_after_tagging.puml

**WITHOUT (Pain Points):**
- "Today, if we want rebasing logic, we hardcode every asset ID"
- "When we add a new rebasing token, we must update code in multiple services"
- "Easy to miss an asset - causes bugs"
- "Logic duplicated across services"

**WITH (Benefits):**
- "With tagging, we query by characteristic: GetByTags('REBASING')"
- "Tag a new asset once in Retool, all services immediately see it"
- "Single source of truth for asset categorization"
- "No deployments needed"

---

### For integrated_asset_and_tagging_flow.puml

**Key Messages:**
- "Tagging doesn't replace dynamic assets - it enhances them"
- "Steps 1-2 unchanged: assets still flow through Special Org"
- "Step 3 NEW: Tag management via Retool"
- "Step 4 enhanced: Client loads assets + tags together"
- "Step 5 enhanced: Services can query by ID or by tags"

**Integration Points:**
- "Reuses existing client refresh mechanism (10s)"
- "Reuses existing database (adds asset_tag table)"
- "Reuses existing gRPC infrastructure"
- "Backward compatible: existing Get(id) still works"

---

## Technical Details for Q&A

### Database Schema

**asset_tag table:**
```sql
CREATE TABLE asset_tag (
    id UUID PRIMARY KEY,
    asset_id VARCHAR NOT NULL,  -- FK to supported_asset.id
    tag VARCHAR NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(asset_id, tag),
    FOREIGN KEY (asset_id) REFERENCES supported_asset(id)
);

CREATE INDEX idx_asset_tag_asset_id ON asset_tag(asset_id);
CREATE INDEX idx_asset_tag_tag ON asset_tag(tag);
```

### Cache Structure

**tagsCache (sync.Map):**
```go
// Key: tag (string)
// Value: []assettype.ID

Example:
tagsCache["REBASING"] = [
    "ethereum-mainnet:0x...",  // stETH
    "ethereum-mainnet:0x...",  // rETH
]

tagsCache["STABLE"] = [
    "ethereum-mainnet:0x...",  // USDC
    "ethereum-mainnet:0x...",  // USDT
]
```

### gRPC APIs

**Read APIs:**
- `ListTags()` - Get all available tag definitions
- `ListTagsWithAssetTypeIDs()` - Get map of tag → []assetTypeIDs

**Write APIs:**
- `AddTagsToAsset(assetID, tags[])` - Add tags to an asset
- `RemoveTagsFromAsset(assetID, tags[])` - Remove tags from an asset

**Authorization:**
- Write APIs restricted to EM google group (TBD exact group)
- Read APIs available to all services

---

## Use Case Examples for Demo

### Example 1: Rebasing Tokens
```go
// Service needs to handle rebasing tokens differently
rebasingAssets := client.GetByTags("REBASING")
for _, asset := range rebasingAssets {
    // Apply rebasing-specific balance update logic
    updateRebasingBalance(asset)
}
```

### Example 2: Finance Billing
```go
// Finance team applies different billing for DeFi assets
defiAssets := client.GetByTags("DEFI")
for _, asset := range defiAssets {
    applyDeFiBilling(asset)
}
```

### Example 3: UI Filtering
```go
// Hide deprecated assets from UI
allAssets := client.List()
activeAssets := filterOut(allAssets, "DEPRECATED")
```

### Example 4: Multi-tag Query
```go
// Get all stable + privacy assets
stablePrivacyAssets := client.GetByTags("STABLE", "PRIVACY")
// Apply special KYC requirements
```

---

## Answers to Anticipated Questions

**Q: Why not put tags in the Special Org alongside assets?**
A: Tags change more frequently than asset definitions. Separating them avoids HSM signing overhead and enables faster iteration.

**Q: Why PR for tag creation but UI for tag application?**
A: Tag *definitions* are governance-critical (affects code behavior). Tag *application* is operational (doesn't change code). Different access patterns.

**Q: What if a service is down when tags are updated?**
A: Client refreshes every 10s. When service comes back up, it gets latest tags. Eventual consistency is acceptable.

**Q: Can we query by multiple tags with AND logic?**
A: Yes, GetByTags() supports both OR and AND logic (match any vs match all).

**Q: What about tag versioning/history?**
A: V0 doesn't have explicit versioning, but entaudit provides audit trail. Can add versioning in future if needed.

**Q: How do we prevent tag misuse?**
A: Three-layer governance:
1. Tag creation requires PR review
2. Tag application restricted to authorized group
3. Documentation and runbook for tag usage

**Q: Performance impact?**
A: Minimal. Tags loaded once every 10s into memory. Lookups are O(1) map access. No network calls.

---

## Presentation Tips

1. **Start broad, go deep**: Begin with problem (before/after), then show solution (flow), then architecture
2. **Use animations**: If presenting live, reveal phases sequentially
3. **Connect to real examples**: Reference actual rebasing tokens (stETH, rETH)
4. **Emphasize no-code**: Highlight that tagging new assets doesn't require deployment
5. **Show progress**: Transition from diagrams to Linear tickets to show work in flight
6. **Invite feedback**: Ask if teams have other tag use cases

---

## Diagram File Summary

| File | Purpose | When to Use |
|------|---------|-------------|
| dynamic_asset_simple_flow.puml | Current asset system | Opening - set context |
| before_after_tagging.puml | Problem + solution | Motivation - show value |
| asset_tagging_flow.puml | Tagging flow detail | Core explanation |
| asset_tagging_architecture.puml | Technical architecture | Deep dive (technical audience) |
| integrated_asset_and_tagging_flow.puml | Integration story | Show how it fits together |

All diagrams available in `/home/user/playground/linear/`
