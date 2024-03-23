package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func init() {
	fmt.Println("Initialized")
	var err error
	items, err = GetItems()
	if err != nil {
		items = make(map[string][]Info)
	}
}

type NewFood struct {
	Measure string  `json:"measure"`
	Food    string  `json:"food"`
	Qtd     float64 `json:"qtd"`
	Kcal    float64 `json:"kcal"`
}

func (nf *NewFood) Valid() bool {
	return nf.Food != "" && nf.Kcal > 0 && nf.Measure != "" && nf.Qtd > 0
}

type Items map[string][]Info

var items Items

type Info struct {
	Measure   string  `json:"measure"`
	QtdByUnit float64 `json:"qtdByUnit"`
}

func getCache(measure, food string) float64 {
	infos := items[food]
	for _, v := range infos {
		if v.Measure == measure {
			fmt.Println("Returning from cache")
			return v.QtdByUnit
		}
	}
	return 0
}

func SaveFood(newFood NewFood) error {
	info := Info{
		Measure:   newFood.Measure,
		QtdByUnit: newFood.Kcal / newFood.Qtd,
	}

	food := newFood.Food
	if items[food] != nil {
		for _, item := range items[food] {
			if item.Measure == newFood.Measure {
				return errors.New("already exist")
			}

		}
	}

	items[food] = append(items[food], info)
	save()
	return nil
}

func GetInfo(measure, food string) float64 {
	cacheValue := getCache(measure, food)

	if cacheValue > 0 {
		return cacheValue
	}

	client := openai.NewClient("API_KEY")

	qtd := 1.0

	if strings.Contains(measure, "grama") {
		qtd = 100
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf("calorias %v %s de %s em uma frase curta", qtd, measure, food),
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return 0
	}

	content := resp.Choices[0].Message.Content
	fmt.Println(content)
	exp := regexp.MustCompile(`(([0-9]*[,])?[0-9]+) caloria`)
	matches := exp.FindAllString(content, -1)

	var qtdByUnit float64
	if len(matches) > 0 {
		target := matches[0]
		qtdByUnit, _ = strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(target, " caloria", ""), ",", "."), 64)
	}
	finalQtdByUnit := qtdByUnit / qtd
	info := Info{
		Measure:   measure,
		QtdByUnit: finalQtdByUnit,
	}

	if finalQtdByUnit > 0 {
		items[food] = append(items[food], info)
		save()
	}

	//fmt.Println(qtdByUnit)

	//fmt.Println(resp.Choices[0].Message.Content)
	return finalQtdByUnit
}

func save() {
	UpdateItems(items)
}
