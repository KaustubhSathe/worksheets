package model

import (
	"encoding/json"
	"time"
)

// This struct represents the SpreadSheet object stored in DB
type SpreadSheet struct {
	Base
	UserName         string
	UserID           int64
	SpreadSheetTitle string
	Favorited        bool
	Versions         []Version // Will contain pointer to S3 objects, index indicates the sheet number
	LastOpened       time.Time
}

type Version struct {
	VersionName string
	VersionID   string
	CreatedAt   time.Time
	Sheets      []Sheet
}

type Sheet struct {
	SheetName  string
	SheetIndex int32
	State      map[string]State
}

type State struct {
	FontWeight      string
	FontSize        int32
	FontStyle       string
	TextDecoration  string
	FontColor       string
	BackGroundColor string
	BackGroundImage string
	FontFamily      string
	TextContent     string
	TextAlign       string
}

func (in *SpreadSheet) Stringify() string {
	b, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(b)
}

func StringifySpreadSheets(in []*SpreadSheet) string {
	b, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(b)
}
