# Asset Registry Tagging Project - Linear Tickets

**Project ID**: d4caeba1-8192-48a7-ac21-c07b53ecf60a
**Total Tickets**: 22
**Last Updated**: 2026-02-18

---

## Executive Summary

The Asset Registry Tagging project enables flexible categorization and discovery of assets through a tagging system. The project includes:
- Tag definition and management infrastructure
- gRPC APIs for tag queries and mutations
- Client library integration
- Retool UI for operations team
- Database schema and indexing

---

## Current Status

### Completed (8 tickets)

#### PROC-4520: Apply Go-idiomatic interface patterns
- **Status**: Done
- **Assignee**: Bruno Vilar
- **Description**: Refactored to separate database access from domain logic using consumer-defined interfaces
- **Key Changes**:
  - Moved database/ent logic out of logic.go into dbprovider layer
  - Adopted patterns from subaccount and custody services
  - No behavior changes, only architectural improvements

#### PROC-4466: Add ListTags gRPC endpoint
- **Status**: Done
- **Assignee**: Bruno Vilar
- **Description**: Public API to fetch all available tag definitions
- **Deliverables**:
  - Proto definition in assetsupport.proto
  - gRPC handler in assetsupport_server.go
  - Returns sorted list of TagDefinition with descriptions
  - Testable with grpcurl

#### PROC-4458: Add tag mutation proto definitions
- **Status**: Done
- **Assignee**: Bruno Vilar
- **Description**: Proto definitions for adding/removing tags
- **Deliverables**:
  - AddTagsToAsset RPC definition
  - RemoveTagsFromAsset RPC definition
  - Request/Response messages
  - Response includes current tags after mutation

#### PROC-4454: Add tag methods to AssetTypesReader interface
- **Status**: Done
- **Assignee**: José Sousa
- **Description**: Extended client interface with tag query methods
- **New Methods**:
  - `GetTagsByAssetTypeIDs(ctx, assetIDs ...assettype.ID)`
  - `GetAssetsTypeIdsByTags(ctx, tag ...string)` (OR logic)

#### PROC-4450: Implement ListTagsWithAssetTypeIDs gRPC handler
- **Status**: Done
- **Assignee**: Bruno Vilar
- **Description**: Handler to retrieve all asset type IDs grouped by tags
- **Location**: internal/grpcserver/assettagging_server.go
- **Returns**: Map of tag to asset type ID list

#### PROC-4448: Implement ListTagsWithAssetTypeIDs service logic
- **Status**: Done
- **Assignee**: Bruno Vilar
- **Description**: Service layer method for efficient tag/asset retrieval
- **Returns**: `map[string][]assettype.ID` with tags grouped by asset
- **Coverage**: >80% unit test coverage

#### PROC-4447: Setup Wire dependency injection
- **Status**: Done
- **Assignee**: Bruno Vilar
- **Description**: Configured dependency injection for tag-related services
- **Location**: internal/cmd/assetregistry/wire.go
- **Verification**: Service compiles and starts without errors

#### PROC-4446: Implement repository layer queries
- **Status**: Done
- **Assignee**: Bruno Vilar
- **Description**: Database query helpers for asset_tag table
- **Deliverables**:
  - ListTagsWithAssetTypeIDs() helper
  - Unit tests with in-memory SQLite
  - Query optimization with indexes

---

### In Progress (3 tickets)

#### PROC-4465: Create Retool UI for tag management
- **Status**: In Progress
- **Assignee**: Bruno Vilar
- **Priority**: 2 (High)
- **URL**: https://retoolv3.anchorage-development.com/apps/58d84db4-acc0-11ef-8293-8bca805579cf/Asset%20types
- **Features**:
  - Searchable asset list with current tags
  - Add tags via dropdown/multi-select
  - Remove tags with confirmation
  - Success/error messaging
  - Audit trail logging
- **Dependencies**: None

#### PROC-4460: Implement tag mutation gRPC handlers
- **Status**: In Progress
- **Assignee**: Bruno Vilar
- **Priority**: 2 (High)
- **Location**: internal/grpcserver/assetsupport_server.go
- **Handlers**:
  - AddTagsToAsset()
  - RemoveTagsFromAsset()
- **Error Handling**:
  - Invalid tags → InvalidArgument status
  - Non-existent assets → NotFound status
  - Database errors → Internal status
- **Authorization**: Only EM google group allowed
- **Dependencies**: None

#### PROC-4456: Implement GetByTags methods in client
- **Status**: In Progress
- **Assignee**: José Sousa
- **Priority**: 2 (High)
- **Description**: Client method to query assets by tags
- **Method**: `GetByTags(ctx, tags)` with OR logic
- **Error Handling**:
  - Empty tags → clear error
  - Invalid tags → server error
  - No matches → empty array (not error)
- **Dependencies**: None

---

### In Review (2 tickets)

#### PROC-4459: Implement tag mutation service logic
- **Status**: In Review
- **Assignee**: Bruno Vilar
- **Priority**: 2 (High)
- **URL**: https://linear.app/anchorlabs/issue/PROC-4459
- **Methods**:
  - `AddTags(ctx, assetID, tags)` - idempotent, ignores duplicates
  - `RemoveTags(ctx, assetID, tags)` - deletes specified tags
- **Validation**:
  - Tags validated before mutation
  - Asset existence verified
  - Clear error messages
- **Audit**: Via entaudit
- **Question**: Should RemoveTags use soft delete?

#### PROC-4455: Load tags in client
- **Status**: In Review
- **Assignee**: José Sousa
- **Priority**: 2 (High)
- **URL**: https://linear.app/anchorlabs/issue/PROC-4455
- **Implementation**:
  - `tagsCache sync.Map` field in assetTypesClient
  - Key: tag, Value: []assetTypeIDs
  - Loaded on client bootstrap via ListTagsWithAssetTypeIDs RPC
  - Thread-safe concurrent access
- **Compatibility**: Ensure existing assettype lookups not broken (feature flag suggested)

---

### Triage/Low Priority (8 tickets)

#### PROC-4463: Create documentation and deployment plan
- **Priority**: 3 (Low)
- **Status**: Triage
- **Scope**:
  - Service README with tag system overview
  - Tag addition process (PR workflow)
  - Client usage guide with examples
  - Operations runbook
  - Database schema documentation
  - Migration guide for service owners
  - API documentation for new RPCs
  - Troubleshooting guide
  - Deployment plan (staging → production)

#### PROC-4462: Create end-to-end smoke tests
- **Priority**: 3 (Low)
- **Status**: Triage
- **Location**: internal/integration/smoketests/
- **Coverage**:
  - Complete tag lifecycle (add → query → mutate → verify)
  - Error scenarios (invalid tags, auth failures)
  - Performance with 100+ assets
  - Tags persist across service restart
  - No memory leaks or performance regressions
- **Environment**: make run-lite

#### PROC-4461: Create integration tests for mutations
- **Priority**: 3 (Low)
- **Status**: Triage
- **Tests**:
  - AddTagsToAsset RPC
  - RemoveTagsFromAsset RPC
  - SetAssetTags RPC
  - Full mutation lifecycle
- **Database**: Real PostgreSQL
- **Verification**: Tags persist, audit trail exists
- **CI**: Added to pipeline

#### PROC-4457: Add client tests and documentation
- **Priority**: 3 (Low)
- **Status**: Triage
- **Scope**:
  - Tests in asset_types_client_test.go
  - GetTags lazy loading behavior
  - GetAssetsByTag methods
  - Example usage in code comments
  - Backward compatibility verification
  - Caching behavior documentation

#### PROC-4453: Create integration tests for read path
- **Priority**: 3 (Low)
- **Status**: Triage
- **Location**: internal/integration/
- **Coverage**:
  - GetAssetTags RPC end-to-end
  - GetAssetsByTags with OR logic (ANY mode)
  - GetAssetsByTags with AND logic (ALL mode)
  - Error cases (invalid tags, missing assets)
  - Performance with 100+ assets
- **Database**: Real PostgreSQL

#### PROC-4452: Add AssetTag to UnitOfWork repository
- **Priority**: 4 (Very Low)
- **Status**: Triage
- **Description**: Enable transactional tag operations
- **Scope**:
  - AssetTag accessible via UnitOfWorkRepository
  - usable within uow.Do() transactions
  - Supports commits and rollbacks
  - Unit tests verify transaction behavior

#### PROC-4445: Add database indexes
- **Priority**: 2 (High)
- **Status**: Triage
- **Scope**:
  - Index on asset_tag.asset_id for fast lookups
  - Index on asset_tag.tag for tag-based queries
  - Composite index on (asset_id, tag) for uniqueness
  - Migration script with index creation
  - EXPLAIN ANALYZE verification
  - No performance regression

#### PROC-4444: Create asset_tag table with Ent
- **Priority**: 2 (High)
- **Status**: Triage
- **Schema**:
  - id (UUID, primary key)
  - asset_id (UUID, references asset_types.id)
  - tag (string)
  - created_at, updated_at (timestamps)
- **Features**:
  - Ent schema definition in schema/assettag.go
  - Foreign key constraint to asset_types
  - Unique constraint on (asset_id, tag)
  - Migration script
  - Up/down migration tested

---

### Canceled (2 tickets)

#### PROC-4451: Implement GetAssetsTypeIDsByTag gRPC handler
- **Status**: Canceled
- **Reason**: Likely consolidated into other handlers

#### PROC-4449: Implement GetAssetsTypeIDsByTag service logic
- **Status**: Canceled
- **Reason**: Likely consolidated into other service methods

---

## Technical Architecture

### Database Layer
- **Table**: asset_tag (id, asset_id, tag, timestamps)
- **Indexes**: asset_id, tag, (asset_id, tag) composite
- **Framework**: Ent ORM
- **Constraints**: Foreign key to asset_types, unique (asset_id, tag)

### Service Layer
- **Package**: internal/tags/
- **Interfaces**: AssetTagging
- **Methods**:
  - ListTagsWithAssetTypeIDs()
  - AddTags() - idempotent
  - RemoveTags()
- **Audit**: Via entaudit

### gRPC Layer
- **Proto Files**: assetsupport.proto, assettagging.proto
- **Handlers**: internal/grpcserver/
- **RPCs**:
  - ListTags
  - ListTagsWithAssetTypeIDs
  - AddTagsToAsset
  - RemoveTagsFromAsset

### Client Layer
- **Package**: lib/clients/assettypes/
- **Interface**: AssetTypesReader
- **Methods**:
  - GetTagsByAssetTypeIDs()
  - GetAssetsTypeIdsByTags()
  - GetByTags()
- **Caching**: Tags loaded on bootstrap via sync.Map

### UI Layer
- **Platform**: Retool
- **Features**: Search, add tags, remove tags, audit trail
- **Authorization**: EM google group (TBD)

---

## Team Assignments

### Bruno Vilar
- Core infrastructure (proto, handlers, service logic)
- Wire dependency injection
- Repository layer
- Retool UI
- Interface patterns refactoring

### José Sousa
- Client library integration
- Tag caching mechanism
- GetByTags implementation

### Unassigned
- Documentation
- Testing (integration, smoke tests)
- Database optimization
- UnitOfWork integration

---

## Key Decisions & Open Questions

1. **Soft Delete**: Should RemoveTags use soft delete? (PROC-4459)
2. **Authorization Group**: What is the exact EM google group for mutations? (PROC-4460)
3. **Feature Flag**: Should tag loading in client use feature flag for compatibility? (PROC-4455)
4. **Performance**: Index verification needed for production scale (PROC-4445)

---

## Next Steps for Presentation

1. **Highlight Completed Work**:
   - Core infrastructure done (proto, handlers, service layer)
   - Client interface extended
   - Dependency injection configured

2. **Show Active Development**:
   - Retool UI in progress
   - Mutation handlers being implemented
   - Client methods under review

3. **Discuss Future Work**:
   - Documentation and deployment plan
   - Comprehensive test coverage
   - Performance optimization

4. **Address Questions**:
   - Soft delete strategy
   - Authorization model
   - Backward compatibility approach
