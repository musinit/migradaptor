package builder

type GooseCmd string

var (
	GooseCmdMigrationDown  GooseCmd = "+goose Down"
	GooseCmdMigrationUp    GooseCmd = "+goose Up"
	GooseCmdStatementBegin GooseCmd = "+goose StatementBegin"
	GooseCmdStatementEnd   GooseCmd = "+goose StatementEnd"
	GooseCmdNoTransaction  GooseCmd = "NO TRANSACTION"
)
