package database_init

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type IDBService interface {
	DeleteTables()
	CreateTable(tableName, fileName string, parser func(line []string) (interface{}, error))
	CreateTableByChan(tableName, fileName string, parser func(line []string) (interface{}, error))
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

func ParseClashes(line []string) (interface{}, error) {
	HexId, err := strconv.Atoi(line[0])
	if err != nil {
		return nil, err
	}

	ClashesShare, err := strconv.ParseFloat(line[1], 64)
	if err != nil {
		return nil, err
	}

	ClashPowerMedian, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return nil, err
	}

	clash := Clash{HexId: HexId, ClashesShare: ClashesShare, ClashPowerMedian: ClashPowerMedian}
	return clash, nil
}

func ParseScooters(line []string) (interface{}, error) {
	TsUtc, err := time.Parse("2006-01-02 15:04:05", line[0])
	if err != nil {
		return nil, err
	}

	ParkingId, err := strconv.Atoi(line[1])
	if err != nil {
		return nil, err
	}

	ScootersAtParking, err := strconv.Atoi(line[2])
	if err != nil {
		return nil, err
	}

	scooters := ScootersAtParkings{TsUtc: TsUtc, ParkingId: ParkingId, ScootersAtParking: ScootersAtParking}
	return scooters, nil
}

func (s *DBService) CreateTable(tableName, fileName string, parser func(line []string) (interface{}, error)) {
	s.repository.CreateTable(tableName)

	dirLink := fmt.Sprintf("C:\\Users\\insha\\OneDrive\\Документы\\whoosh\\Datasets\\%s", fileName) // path in computer
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

func (s *DBService) CreateTableByChan(tableName, fileName string, parser func(line []string) (interface{}, error)) {
	s.repository.CreateTable(tableName)

	dirLink := fmt.Sprintf("C:\\Users\\insha\\OneDrive\\Документы\\whoosh\\Datasets\\%s", fileName) // path in computer
	csvFile, err := os.Open(dirLink)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	//chSize := 10
	i := 0
	ch := make(chan interface{})

	go s.repository.AddRowsByChan(ch, tableName)

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
			log.Fatal(err)
		}

		ch <- value
	}
}
