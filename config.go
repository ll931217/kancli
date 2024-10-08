package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	DbPath string `json:"dbPath"`
}

func readConfig() Config {
	// look for $XDG_CONFIG_HOME/kancli/config.json or $HOME/.config/kancli/config.json
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// cant find home or config just give up
			return Config{}
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	// Create ~/.config/kancli/ if it does not exist
	configPath := filepath.Join(configDir, "kancli")

	mkdirErr := os.MkdirAll(configPath, 0755)
	if mkdirErr != nil {
		log.Fatal(mkdirErr)
	}

	configFile := filepath.Join(configPath, "config.json")

	var config Config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		csvFile := configDir + "/tasks.csv"
		file.WriteString(fmt.Sprintf("{\"dbPath\": \"%s\"}", csvFile))
	}

	data, readJSONFileErr := os.ReadFile(configFile)
	if readJSONFileErr != nil {
		log.Fatal(readJSONFileErr)
	}

	if readJSONerr := json.Unmarshal(data, &config); readJSONerr != nil {
		log.Fatal(readJSONerr)
	}

	return config
}

func readCSV() [][]string {
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		file, err := os.Create(csvFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString("title,description,status\n")
	}

	content, err := os.ReadFile(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(content)))
	r.Comma = ','

	// Skip the header
	_, readErr := r.Read()
	if readErr != nil {
		log.Fatal(readErr)
	}

	records, readAllErr := r.ReadAll()
	if readAllErr != nil {
		log.Fatal(readAllErr)
	}

	return records
}

func updateCSV() {
	file, err := os.Create(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("title,description,status\n")
	w := csv.NewWriter(file)
	for _, col := range board.cols {
		for _, item := range col.list.Items() {
			task := item.(Task)
			w.Write([]string{task.Title(), task.Description(), task.status.String()})
		}
	}
	w.Flush()
}
