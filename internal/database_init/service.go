package database_init

import (
	"encoding/csv"
	"fmt"
	"github.com/atomicai/whoosh/internal/models"
	"io"
	"log"
	"os"
	"strconv"
)

type IDBService interface {
	DeleteTables()
	CreateTable(tableName, fileName string, parser func(line []string) (interface{}, error))
}

type DBService struct {
	repository IDBRepository
}

func NewDBService(dbname string) *DBService {
	return &DBService{repository: NewDBRepository(dbname)}
}

func (s *DBService) DeleteTables() {
	s.repository.DeleteTables()
}

func (s *DBService) CreateTable(tableName, fileName string, parser func(line []string) (interface{}, error)) {
	s.repository.CreateTable(tableName)

	dirLink := fmt.Sprintf("%s", fileName) // path in computer
	csvFile, err := os.Open(dirLink)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	i := 0
	arrSize := 10_000
	values := make([]interface{}, 0, arrSize)
	for ; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if i == 0 {
			continue // skip [ts_utc parking_id scooters_at_parkings]
		}
		if err != nil {
			log.Fatal(err)
		}

		value, err := parser(line)
		if err != nil {
			fmt.Println("error on parsing")
			log.Fatal(err)
		}

		values = append(values, value)
		if i%arrSize == 0 {
			s.repository.AddRows(&values, tableName)
			values = make([]interface{}, 0, arrSize)
		}
	}
	s.repository.AddRows(&values, tableName)
}

func ParseGraph(line []string) (interface{}, error) {
	graph := models.NewNode()

	id, err := strconv.Atoi(line[0])
	if err != nil {
		return nil, err
	}
	graph.NodeId = id

	lon, err := strconv.ParseFloat(line[1], 64)
	if err != nil {
		return nil, err
	}
	graph.Lon = lon

	lat, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return nil, err
	}
	graph.Lat = lat

	for i := 0; i < 6; i++ {
		idTo, err := strconv.ParseFloat(line[3+i*2], 32)
		if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(line[4+i*2], 64)
		if err != nil {
			return nil, err
		}

		if idTo != -1.0 && weight != 1.0 {
			graph.Neighbors = append(graph.Neighbors, &models.Edge{To: int(idTo), Weight: weight})
		}
	}

	isParking, err := strconv.ParseFloat(line[15], 32)
	if err != nil {
		return nil, err
	}
	graph.IsParking = isParking

	nodeRoadIndex, err := strconv.ParseFloat(line[16], 64)
	if err != nil {
		return nil, err
	}
	graph.NodeRoadIndex = nodeRoadIndex

	nodeStar, err := strconv.ParseFloat(line[17], 64)
	if err != nil {
		return nil, err
	}
	graph.NodeStar = nodeStar

	nodeSpeed, err := strconv.ParseFloat(line[18], 64)
	if err != nil {
		return nil, err
	}
	graph.NodeSpeed = nodeSpeed

	return graph, nil
}
