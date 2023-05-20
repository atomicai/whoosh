package database_init

type DBHanlder struct {
	dbService IDBService
}

func NewDBHandler(dbname string) *DBHanlder {
	return &DBHanlder{
		dbService: NewDBService(dbname),
	}
}

func (h *DBHanlder) DeleteTables() {
	h.dbService.DeleteTables()
}

func (h *DBHanlder) CreateTables() {
	h.dbService.CreateTable("scooters_at_parkings", "scooters_at_parkings (1).csv", ParseScooters)
	h.dbService.CreateTable("clashes", "clashes (1).csv", ParseClashes)
}
