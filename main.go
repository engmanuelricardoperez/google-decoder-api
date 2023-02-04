package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ResponseFile struct {
	ID  int     `json:"id"`
	Dir string  `json:"dir"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Response struct {
	Results []Results `json:"results"`
}

type Results struct {
	FormattedAddress string `json:"formatted_address"`
}

func main() {
	// Open File of Apps
	f, err := os.Open("coordenadas.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	fmt.Println(len(data))
	if err != nil {
		log.Fatal(err)
	}

	// Create File of Apps
	csvFile, err := os.Create("./data.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	for i, line := range data {

		if i > 0 { // omit header line
			// project := "fury_" + line[0]
			coordinates := strings.Replace(line[0], ";", ",", 1)
			// fmt.Println("Processing project>", project)
			index := strings.Index(line[0], ";")
			lat := line[0][:index]
			lon := line[0][index+1 : len(line[0])]
			///////////////////////////////////////////
			response, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?latlng=" + coordinates + "&sensor=true&key=API_KEY_GOOGLE_MAPS_PLATFORM_HERE")

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(string(responseData))

			var responseObject Response
			json.Unmarshal(responseData, &responseObject)

			// fmt.Println(len(responseObject.Results))
			country := ""
			departament := ""
			city := ""
			address := ""

			split := strings.Split(responseObject.Results[0].FormattedAddress, ",")
			for i := range split {
				split[i] = strings.TrimSpace(split[i])
			}
			if len(split) == 2 {
				if split[1] == "Bogotá" {
					country = split[2]
					departament = split[1]
					city = split[1]
					address = split[0]
				} else {
					country = split[2]
					departament = split[1]
					city = split[0]
					address = split[0]
				}
			}

			if len(split) == 3 {
				if split[1] == "Bogotá" {
					country = split[2]
					departament = split[1]
					city = split[1]
					address = split[0]
				} else {
					country = split[2]
					departament = split[1]
					city = split[0]
					address = split[0]
				}
			}

			if len(split) == 4 {
				country = split[3]
				departament = split[2]
				city = split[1]
				address = split[0]
			}

			if len(split) == 5 {
				country = split[4]
				departament = split[3]
				city = split[2]
				address = split[1]
			}

			fmt.Println(responseObject.Results[0].FormattedAddress, " - Lat: ", lat, " - Lon: ", lon)
			///////////////////////////////////////////
			var row []string
			row = append(row, strconv.Itoa(i))
			row = append(row, lat)
			row = append(row, lon)
			row = append(row, strings.Trim(address, " "))
			row = append(row, strings.Trim(city, " "))
			row = append(row, strings.Trim(departament, " "))
			row = append(row, strings.Trim(country, " "))
			writer.Write(row)
			///////////////////////////////////////////
		} else {

			var rowHeader []string
			rowHeader = append(rowHeader, string("id"))
			rowHeader = append(rowHeader, string("lat"))
			rowHeader = append(rowHeader, string("lon"))
			rowHeader = append(rowHeader, string("address"))
			rowHeader = append(rowHeader, string("city"))
			rowHeader = append(rowHeader, string("departament"))
			rowHeader = append(rowHeader, string("country"))
			writer.Write(rowHeader)

		}
	}
}
