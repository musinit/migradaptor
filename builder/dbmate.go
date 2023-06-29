package builder

type DbmateCmd string

var (
	DbmateCmdMigrationDown DbmateCmd = "migrate:down"
	DbmateCmdMigrationUp   DbmateCmd = "migrate:up"
	DbmateCmdNoTransaction DbmateCmd = "transaction:false"
)
