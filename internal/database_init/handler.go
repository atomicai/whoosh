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
	h.dbService.CreateTable("graph", "new_graph.csv", ParseGraph)
	//h.dbService.CreateTable("clashes", "clashes (1).csv", ParseClashes)
	//h.dbService.CreateTable("road_index", "road_index (1).csv", ParseRoads)
	//h.dbService.CreateTable("routers_hex20m", "routes_hex20m.csv", ParseRoutes)
	//h.dbService.CreateTable("scooters_at_parkings", "scooters_at_parkings (1).csv", ParseScooters)
	//h.dbService.CreateTable("speed_median_hex20m_hackaton", "speed_median", ParseSpeedMedian)

}
