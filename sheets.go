package main

import (
	"context"
	"os"
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

func createValues(item string, price string, comment string) sheets.ValueRange {
	sgtTz, _ := time.LoadLocation("Asia/Singapore")
	dateTime := time.Now().In(sgtTz).Format("2006/01/02 15:04:05")

	return sheets.ValueRange{Values: [][]interface{}{{dateTime, item, price, comment}}}
}

func insertValuesToExpensesSheet(values sheets.ValueRange) bool {
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

	resp, err := service.Spreadsheets.Values.Append(SHEETS_ID, "01_expenses!A1:D1", &values).
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
