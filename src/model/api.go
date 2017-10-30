package model

import (
	"fmt"
	"rent-notifier/src/db"
	"strings"
	"bytes"
)

func FormatMessage(db *dbal.DBAL, note dbal.Note) string {

	var b bytes.Buffer

	b.WriteString("<b>")
	b.WriteString(formatType(note.Type))
	price := formatPrice(note.Price)
	if price != "" {
		b.WriteString(" за ")
		b.WriteString(price)
		b.WriteString(" руб/мес")
	}

	text_subways := formatSubways(db, note)
	if text_subways != "" {
		b.WriteString(" около метро ")
		b.WriteString(text_subways)
	}
	b.WriteString("</b>")
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("%s <a href='%s'>Перейти к объявлению</a>", note.Contact, note.Link))

	b.WriteString("\n\n")

	b.WriteString(note.Description)

	if len(note.Photos) != 0 {
		b.WriteString("\n")
		b.WriteString(note.Photos[0])
	} else {
		b.WriteString("https://socrent.ru/img/logo.png")
	}

	return b.String()
}

func formatType(note_type int) string {
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

func formatPrice(price int) string {
	return fmt.Sprintf("%d", price)
}

func formatSubways(db *dbal.DBAL, note dbal.Note) string {

	subways := make([]string, 0)
	for _, subway := range db.FindSubwaysByIds(note.Subways) {
		subways = append(subways, subway.Name)
	}

	if len(subways) == 0 {
		return ""
	}

	return strings.Join(subways, ", ")
}
