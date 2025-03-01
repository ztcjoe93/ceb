package main

import (
	"bufio"
	"context"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

var (
	logger     = logrus.New()
	bot        *tgbotapi.BotAPI
	tgbotToken = os.Getenv("TGBOT_TOKEN")
)

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI(tgbotToken)
	if err != nil {
		logger.Panic(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates := bot.GetUpdatesChan(u)
	go receiveUpdates(ctx, updates)

	logger.Println("Start listening for updates. Press enter to stop")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		handleMessage(update.Message)
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	logger.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(message.Chat.ID, text)
	} else {
		copyMsg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
		_, err = bot.CopyMessage(copyMsg)
	}

	if err != nil {
		logger.Printf("An error occured: %s", err.Error())
	}
}

func handleCommand(chatId int64, command string) error {
	var err error

	re := regexp.MustCompile(`^(/\S+)\s+(.*)$`)
	matches := re.FindStringSubmatch(command)

	if len(matches) < 3 {
		logger.Infof("Invalid command + args: %v", command)
		msg := tgbotapi.NewMessage(chatId, "Invalid command")
		_, err := bot.Send(msg)
		if err != nil {
			logger.Errorf("Failed to send message: %v", err)
		}

		return err
	}

	cmd := matches[1]
	args := strings.Split(matches[2], ",")

	switch cmd {
	case "/record":
		logger.Info("Received /record invocation")
		err = saveRecord(chatId, args)
	}

	return err
}

func saveRecord(chatId int64, args []string) error {
	var values sheets.ValueRange

	if len(args) > 2 {
		values = createValues(args[0], args[1], args[2])
	} else {
		values = createValues(args[0], args[1], "")
	}
	success := insertValuesToExpensesSheet(values)

	var msg tgbotapi.MessageConfig
	if success {
		msg = tgbotapi.NewMessage(chatId, "Successfully saved record to google sheets!")
	} else {
		msg = tgbotapi.NewMessage(chatId, "Record saving failed, please check logs...")
	}
	_, err := bot.Send(msg)

	return err
}
