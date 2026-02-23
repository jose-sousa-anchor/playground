# Dynamic Asset Addition - Presentation Guide

## Overview
This document explains the current system for dynamically adding assets to Anchorage's platform without code changes or deployments.

---

## Diagrams Created

1. **dynamic_asset_simple_flow.puml** - Clean sequence diagram showing the end-to-end flow (RECOMMENDED FOR PRESENTATION)
2. **dynamic_asset_architecture.puml** - Component diagram showing system architecture
3. **dynamic_asset_flow.puml** - Detailed sequence diagram with all technical details

---

## Key Concepts to Present

### 1. Special Organization (ASSET_SUPPORT_V2)
- **Purpose**: Dedicated organization that stores asset configuration metadata
- **Structure**: Contains `AnchorageData.AssetConfigs[]` array
- **Security**: HSM-signed for authenticity and integrity
- **Location**: Retrieved via `GetOrganizationForPurpose(ASSET_SUPPORT_V2)`

### 2. The Four-Phase Flow

#### Phase 1: Asset Addition & Approval
**What happens:**
- Admin submits `requestAddAssetSupport` GraphQL mutation
- System creates `AssetSupportAddChange` operation
- Change modifies the special organization's AssetConfigs array
- Requires biometric approval (Bio challenge)
- HSM signs the updated organization state

**Key files:**
- `mutations.go:643-690` - GraphQL mutation handler
- `asset_support_add_changes.go` - Change implementation
- `operations_handlers.go` - HSM signing workflow

**Duration:** Minutes (depends on approval)

---

#### Phase 2: Asset Registry Synchronization
**What happens:**
- Cron job runs every 60 seconds
- Retrieves the signed special organization
- Extracts all AssetConfigs from the organization
- Upserts each asset into local `SupportedAsset` database table
- Creates asset enablement records

**Key files:**
- `cron_update_assets.go:45-193` - Cron implementation
- Database table: `supported_asset`

**Duration:** Up to 60 seconds from approval

**Why this approach:**
- Single source of truth (special org)
- HSM-signed data ensures authenticity
- Transactional database storage for reliability
- Enables asset querying and filtering

---

#### Phase 3: Client Cache Population
**What happens:**
- AssetTypes client runs refresh job every 10 seconds
- Calls `ListAssetTypes` gRPC endpoint (OnlyXpress: true)
- Registers each asset via `RegisterToken()`
- Stores in in-memory cache (`sync.Map`)
- Triggers registered callbacks to notify consumers

**Key files:**
- `asset_types_client.go` - Client implementation
- Cache structure: `assetTypesMap sync.Map`

**Duration:** Up to 10 seconds from database storage

**Why this approach:**
- Fast in-memory access for services
- Automatic updates without service restarts
- Callback mechanism for reactive updates
- Supports both static and dynamic assets

---

#### Phase 4: Service Usage
**What happens:**
- Services call `Get(assetTypeID)` or `List()` on client
- Client returns cached asset information
- Services use asset data for:
  - Transaction validation
  - Balance calculations
  - Display formatting (decimals, names)
  - Network routing

**Key files:**
- `asset_types_client.go:Get()`, `List()` methods

**Duration:** Nanoseconds (in-memory cache lookup)

---

## System Properties

### Consistency Model
- **Special Org**: Strong consistency via HSM signing
- **Database**: Eventually consistent (up to 60s lag)
- **Cache**: Eventually consistent (up to 10s lag from DB)
- **Total latency**: Up to 70 seconds from approval to cache

### Reliability
- **Single Source of Truth**: Special organization (HSM-signed)
- **Idempotent Operations**: Cron upserts are safe to retry
- **Graceful Degradation**: Services continue with cached assets if updates fail
- **Audit Trail**: All changes tracked via organization history

### Security
- **Authorization**: Biometric approval required
- **Integrity**: HSM signature verification
- **Audit**: Complete operation history in organization
- **Access Control**: Special org access controlled by service accounts

### Scalability
- **Read-heavy**: Cache serves unlimited read requests
- **Write-light**: Asset additions are infrequent
- **Horizontal scaling**: Each service maintains its own cache
- **No coordination**: Services update independently

---

## Timing Breakdown

| Phase | Latency | Notes |
|-------|---------|-------|
| Mutation submitted | 0s | User action |
| Approval | ~minutes | Human approval with biometric |
| HSM signing | ~seconds | HSM processing |
| Cron retrieval | up to 60s | Next cron cycle |
| Client refresh | up to 10s | Next client refresh |
| **Total** | **~70s + approval** | From approval to cache |

---

## Key Technical Details

### Special Organization Structure
```protobuf
message Organization {
  AnchorageData anchorage_data = X;
}

message AnchorageData {
  repeated AssetConfig asset_configs = 1;
}

message AssetConfig {
  string network_id = 1;           // e.g., "ethereum-mainnet"
  string asset_type_id = 2;        // e.g., "ethereum-mainnet:0x..."
  string token_type = 3;           // e.g., "ERC20"
  string name = 4;                 // e.g., "USD Coin"
  string abbrev = 5;               // e.g., "USDC"
  string on_chain_identifier = 6;  // e.g., "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
  int64 decimals = 7;              // e.g., 6
  string symbol = 8;               // e.g., "USDC"
}
```

### Database Schema
```sql
CREATE TABLE supported_asset (
  id VARCHAR PRIMARY KEY,          -- asset_type_id
  name VARCHAR NOT NULL,
  abbreviation VARCHAR NOT NULL,
  decimals INTEGER NOT NULL,
  on_chain_identifier VARCHAR,
  network_id VARCHAR NOT NULL,
  token_type VARCHAR NOT NULL,
  is_xpress BOOLEAN DEFAULT true,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  created_by VARCHAR,
  updated_by VARCHAR
);
```

### Cache Structure
```go
type assetTypesClient struct {
  assetTypesMap sync.Map // map[assettype.ID]RegisteredAsset
  callbacks []func(context.Context, RegisteredAsset) error
}

type RegisteredAsset struct {
  AssetInfo assettype.Info
  RegistrationSource RegistrationSource // Static or Dynamic
}
```

---

## Comparison: Before vs After Tagging

### Current System (Dynamic Assets)
- Assets stored in special organization
- No categorization or grouping
- Services query all assets or by specific ID
- No way to discover "rebasing assets" or "stablecoins"

### Future System (With Tagging)
- Assets still in special organization
- **NEW**: Tags stored in asset_tag table
- **NEW**: Services can query by tag: `GetByTags("REBASING")`
- **NEW**: Dynamic discovery of asset categories
- **NEW**: Tags loaded into client cache alongside assets

---

## Discussion Points

### Why Special Organization?
- ✅ Single source of truth
- ✅ HSM-signed for security
- ✅ Built-in audit trail
- ✅ Existing approval workflow
- ✅ Can be retrieved by any service

### Why Two-Stage Sync (Cron → Client)?
- ✅ **Database persistence**: Reliable storage, queryable
- ✅ **Client cache**: Fast access, no network calls
- ✅ **Separation of concerns**: Registry service vs clients
- ✅ **Resilience**: Services continue with cache if registry is down

### Why 60s Cron + 10s Client Refresh?
- ✅ **60s cron**: Special org changes are infrequent, no need for aggressive polling
- ✅ **10s client**: Faster propagation to services, minimal overhead
- ✅ **Staggered updates**: Reduces thundering herd on special org service

### Why Not Event-Driven?
- ❌ More complex implementation
- ❌ Requires message broker infrastructure
- ❌ More failure modes to handle
- ✅ Polling is simple and reliable for low-frequency updates
- 💡 **Could consider for future optimization**

---

## Presentation Tips

1. **Start with the simplified diagram** (`dynamic_asset_simple_flow.puml`)
   - Walk through each phase
   - Emphasize the "single source of truth" concept

2. **Show the architecture diagram** (`dynamic_asset_architecture.puml`)
   - Highlight the component interactions
   - Point out the separation of concerns

3. **Discuss timing and consistency**
   - Up to 70 seconds propagation time
   - Why this is acceptable for asset additions
   - How it compares to code deployment cycles (minutes/hours)

4. **Transition to tagging**
   - Current system has no asset categorization
   - Tagging builds on this foundation
   - Same sync mechanism, additional metadata

5. **Highlight benefits**
   - No code changes required
   - No service restarts needed
   - Complete audit trail
   - HSM security guarantees

---

## Questions to Anticipate

**Q: Why not real-time updates?**
A: Asset additions are infrequent (weekly/monthly). 70-second latency is acceptable and simpler than event-driven architecture.

**Q: What if cron fails?**
A: Next cycle (60s) will retry. Services continue with current cache. No data loss.

**Q: What if client refresh fails?**
A: Next cycle (10s) will retry. Service continues with current cache.

**Q: What about asset updates (not additions)?**
A: Same flow. Cron detects changes and updates database. Client refresh picks up changes.

**Q: Can we make it faster?**
A: Yes, but adds complexity. Could reduce cron interval or add event-driven updates. Current speed sufficient for use case.

**Q: What about asset removal?**
A: Not currently supported. Would require deletion/soft-delete logic in cron and cache invalidation.

**Q: How does tagging fit in?**
A: Tags stored separately in asset_tag table. Retrieved alongside assets. Same caching mechanism.

---

## Next Slide: Asset Registry Tagging

This sets up the perfect transition to your tagging presentation:

*"Now that we understand how dynamic assets work, let's talk about a critical missing piece: **Asset Categorization and Discovery**. That's where Asset Registry Tagging comes in..."*

Then show the Linear tickets and implementation progress.
