package builder

type MigrationType string

var (
	MigrationTypeUp   = "up"
	MigrationTypeDown = "down"
)

type Command struct {
	MigrationType
	Statements []string
}

func GetConcurrentIndexesRows(lineJoined string) {

}
