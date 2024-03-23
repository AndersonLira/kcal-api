package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetItems() (Items, error) {
	path := getPath()
	jsonFile, err := os.Open(path)

	if err != nil {
		return Items{}, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var items Items
	json.Unmarshal(byteValue, &items)

	return items, nil
}

func UpdateItems(items Items) (Items, error) {
	bytes, err := json.Marshal(items)
	if err != nil {
		return items, nil
	}
	err = ioutil.WriteFile(getPath(), bytes, 0644)
	return items, err
}

func getPath() string {

	folder := os.Getenv("KCAL_API_DB_FOLDER")

	return fmt.Sprintf("%s/items.json", folder)
}
