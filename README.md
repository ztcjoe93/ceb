# Telegram bot to record Clover's expenses
My :cat: has great expenses and the hoomans need to record them down quickly on Google sheets, without using Google sheets.

## Setup
To run the server:  
```shell
$ go run .
```

Ensure that the following variables are configured:  
| Environment variable | Description |
| --- | --- | 
| `SHEETS_ID` | ID of the Google Spreadsheet to append expense rows |
| `SVC_ACC_EMAIL` | Email of the Google service account to access the spreadsheet |
| `SVC_ACC_KEY` | Private key of the Google service account |
| `SVC_ACC_KEY_ID` | ID of the private key of the Google service account |
| `TGBOT_TOKEN` | Token of the telegram bot |

| Command | Description | Syntax | 
| --- | --- | --- |
| `/record` | Insert an expense record into Google sheets | `item, price, comment` |