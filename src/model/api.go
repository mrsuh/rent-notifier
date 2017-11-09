package model

import (
	"fmt"
	"rent-notifier/src/db"
	"strings"
	"strconv"
)

func FormatType(note_type int) string {
	type_string := "";
	if note_type == 0 {
		type_string = "Комната";
	} else if note_type == 1 {
		type_string = "1 комнатная квартира";
	} else if note_type == 2 {
		type_string = "2 комнатная квартира";
	} else if note_type == 3 {
		type_string = "3 комнатная квартира";
	} else if note_type == 4 {
		type_string = "4+ комнатная квартира";
	} else if note_type == 5 {
		type_string = "Студия";
	}

	return type_string;
}

func FormatTypes(types []int) string {
	types_string := make([]string, 0)
	for _, rent_type := range types {
		types_string = append(types_string, FormatType(rent_type))
	}

	return strings.Join(types_string, ", ")
}

func FormatPrice(price int) string {

	priceStr := strconv.Itoa(price)
	if len(priceStr) > 3 {
		fmt.Sprintf("%s %s", priceStr[0:(len(priceStr)-3)], priceStr[(len(priceStr)-3):])
	}

	return priceStr
}

func FormatSubways(db *dbal.DBAL, subway_ids []int) string {

	subways := make([]string, 0)
	for _, subway := range db.FindSubwaysByIds(subway_ids) {
		subways = append(subways, subway.Name)
	}

	if len(subways) == 0 {
		return ""
	}

	return strings.Join(subways, ", ")
}
