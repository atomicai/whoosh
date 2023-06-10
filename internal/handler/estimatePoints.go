package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func EstimatePoints() {
	csvFile, err := os.Open("parkings.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	data := [][]string{
		{"id", "distance"},
	}
	// create a file
	file, err := os.Create("estimatePointsDistance.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)

	defer writer.Flush()

	// write all rows at once
	writer.WriteAll(data)

	// write single row

	i := 0

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
		xstr := line[1]
		ystr := line[2]
		xFloat, err := strconv.ParseFloat(xstr, 64)
		yFloat, err := strconv.ParseFloat(ystr, 64)
		yDelta := 0.0089832
		xDelta := 0.00923248
		httpPostUrl := "https://routing.api.2gis.com/carrouting/6.0.0/global?key=17fb5617-40a7-496c-9cc0-6bf1b70f55fa"

		arr := CreatePaths(xFloat, yFloat, xDelta, yDelta)
		pathSum := 0
		for j := 0; j < len(arr); j++ {
			stringData := arr[j]
			bytesData := []byte(stringData)

			req, err := http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(bytesData))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			if err != nil {
				panic(err)
			}

			client := &http.Client{}
			res, err := client.Do(req)

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			dist := takeDistance(string(body))
			pathSum += dist
		}
		averagePath := pathSum / 4
		fmt.Println(averagePath)
		extraData := []string{line[0], strconv.Itoa(averagePath)}
		writer.Write(extraData)
	}

}

func CreatePaths(xFloat, yFloat, xDelta, yDelta float64) []string {
	arr := make([]string, 0, 4)
	x := [4]float64{xDelta, 0, -xDelta, 0}
	y := [4]float64{0, yDelta, 0, -yDelta}
	for i := 0; i < 4; i++ {
		xNow := xFloat + x[i]
		yNow := yFloat + y[i]
		stringData := fmt.Sprintf(`{
        "points": [
         {
             "type": "walking",
             "x": %f,
             "y": %f
         },
         {
             "type": "walking",
             "x": %f,
             "y": %f
         }
       ]
      }`, xNow, yNow, xFloat, yFloat)
		arr = append(arr, stringData)
	}
	return arr
}

func takeDistance(body string) int {
	index := strings.Index(body, `"total_distance":`)
	if index == -1 {
		fmt.Println(body)
		log.Fatal("no total distance in string")
		return -1
	}

	needIndex := index + len(`"total_distance":`)
	var numString string
	for i := needIndex; i < len(body); i++ {
		if body[i] == ',' {
			numString = body[needIndex:i]
			break
		}
	}

	num, err := strconv.Atoi(numString)
	if err != nil {
		fmt.Println("error on converting string to int")
	}

	return num
}
