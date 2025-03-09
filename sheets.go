package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	SHEETS_ID                 string          = os.Getenv("SHEETS_ID")
	SHEETS_INSERT_DATA_OPTION string          = "INSERT_ROWS"
	SHEETS_VALUE_INPUT_OPTION string          = "RAW"
	SVC_ACC_EMAIL             string          = os.Getenv("SVC_ACC_EMAIL")
	SVC_ACC_KEY               string          = os.Getenv("SVC_ACC_KEY")
	SVC_ACC_KEY_ID            string          = os.Getenv("SVC_ACC_KEY_ID")
	ctx                       context.Context = context.Background()
)

func createExpenseValue(item string, price string, comment ...string) sheets.ValueRange {
	sgtTz, _ := time.LoadLocation("Asia/Singapore")
	dateTime := time.Now().In(sgtTz)

	price64, _ := strconv.ParseFloat(price, 64)

	if len(comment) > 0 {
		return sheets.ValueRange{Values: [][]interface{}{{dateTime, item, price64, comment[0]}}}
	} else {
		return sheets.ValueRange{Values: [][]interface{}{{dateTime, item, price64}}}
	}
}

func createWeightValue(weight string, comment ...string) sheets.ValueRange {
	sgtTz, _ := time.LoadLocation("Asia/Singapore")
	dateTime := time.Now().In(sgtTz)

	weight64, _ := strconv.ParseFloat(weight, 64)

	if len(comment) > 0 {
		return sheets.ValueRange{Values: [][]interface{}{{dateTime, weight64, comment[0]}}}
	} else {
		return sheets.ValueRange{Values: [][]interface{}{{dateTime, weight64}}}
	}
}

func insertValuesToSheet(sheetType string, values sheets.ValueRange) bool {
	config := &jwt.Config{
		Email:        SVC_ACC_EMAIL,
		PrivateKey:   []byte(strings.ReplaceAll(SVC_ACC_KEY, `\n`, "\n")),
		PrivateKeyID: SVC_ACC_KEY_ID,
		TokenURL:     "https://oauth2.googleapis.com/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	service, err := sheets.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		logger.Errorf("Unable to retrieve Sheets client: %v", err)
		return false
	}

	var sheetRange string
	switch sheetType {
	case "/expenses":
		sheetRange = "expenses!A1:D1"
	case "/weight":
		sheetRange = "weight!A1:B1"
	default:
		logger.Errorf("Invalid sheetType: %v", sheetType)
	}

	resp, err := service.Spreadsheets.Values.Append(SHEETS_ID, sheetRange, &values).
		ValueInputOption(SHEETS_VALUE_INPUT_OPTION).
		InsertDataOption(SHEETS_INSERT_DATA_OPTION).
		Context(ctx).Do()
	if err != nil {
		logger.Errorf("Unable to insert data: %v", err)
		return false
	}

	logger.Info(resp)
	return true
}
