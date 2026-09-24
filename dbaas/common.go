package dbaas

// SnapshotType and SnapshotStatus are shared between instance snapshots and
// cluster snapshots.
type (
	SnapshotType   string
	SnapshotStatus string
)

const (
	SnapshotTypeOnDemand  SnapshotType = "ON_DEMAND"
	SnapshotTypeAutomated SnapshotType = "AUTOMATED"
)

const (
	SnapshotStatusPending   SnapshotStatus = "PENDING"
	SnapshotStatusCreating  SnapshotStatus = "CREATING"
	SnapshotStatusAvailable SnapshotStatus = "AVAILABLE"
	SnapshotStatusRestoring SnapshotStatus = "RESTORING"
	SnapshotStatusError     SnapshotStatus = "ERROR"
	SnapshotStatusDeleting  SnapshotStatus = "DELETING"
	SnapshotStatusDeleted   SnapshotStatus = "DELETED"
)
