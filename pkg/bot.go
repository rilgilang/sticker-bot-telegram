package pkg

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"strings"
	"tele-sticker-finder/config"
	"tele-sticker-finder/pkg/connector"
	"time"
)

// var resty =
// This bot demonstrates some example interactions with commands on telegram.
// It has a basic start command with a bot intro.
// It also has a source command, which sends the bot sourcecode, as a file.
func StartBot(cfg *config.Config) {
	// Get token from the environment variable
	token := cfg.BotToken
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	// /find command to find specific sticker
	dispatcher.AddHandler(handlers.NewCommand("find", find))
	// /update command to update sticker collection to bot server
	dispatcher.AddHandler(handlers.NewCommand("update", update))
	// /source command to send the bot source code
	dispatcher.AddHandler(handlers.NewCommand("source", source))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func source(b *gotgbot.Bot, ctx *ext.Context) error {
	// Sending a file by file handle
	f, err := os.Open("samples/commandBot/main.go")
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}

	m, err := b.SendDocument(ctx.EffectiveChat.Id,
		gotgbot.InputFileByReader("source.go", f),
		&gotgbot.SendDocumentOpts{
			Caption: "Here is my source code, by file handle.",
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId: ctx.EffectiveMessage.MessageId,
			},
		})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}

	// Or sending a file by file ID
	_, err = b.SendDocument(ctx.EffectiveChat.Id,
		gotgbot.InputFileByID(m.Document.FileId),
		&gotgbot.SendDocumentOpts{
			Caption: "Here is my source code, sent by file id.",
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId: ctx.EffectiveMessage.MessageId,
			},
		})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}

	return nil
}

// start introduces the bot.
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("PERINGATAN DARURAT, penghytaman segera dimulai!"), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

// start introduces the bot.
func find(b *gotgbot.Bot, ctx *ext.Context) error {
	set, err := b.GetStickerSet("DUBxVZp16HGJeocv44Bjv4crsfhRvSVA", nil)
	if err != nil {
		return err
	}

	sticker := gotgbot.InputFileByID(set.Stickers[1].FileId)
	_, err = b.SendSticker(ctx.EffectiveChat.Id, sticker, nil)
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	file, _ := b.GetFile(set.Stickers[1].FileId, nil)

	fmt.Println("file urll ", file.URL(b, nil))

	// Try to get the underlying struct
	//if fr, ok := sticker.(*gotgbot.FileReader); ok {
	//	fr.
	//} else {
	//	fmt.Println("Not a FileReader")
	//}

	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

// update sticker
func update(b *gotgbot.Bot, ctx *ext.Context) error {
	set, err := b.GetStickerSet("DUBxVZp16HGJeocv44Bjv4crsfhRvSVA", nil)
	if err != nil {
		return err
	}

	restyClient := resty.New()

	for _, sticker := range set.Stickers {
		file, err := b.GetFile(sticker.FileId, nil)
		if err != nil {
			panic(err)
		}

		fileUrl := file.URL(b, nil)

		response, err := restyClient.R().Get(fileUrl)
		if err != nil {
			panic(err)
		}

		filePath, err := saveSticker(response, sticker.FileId, strings.Split(fileUrl, ".")[len(strings.Split(fileUrl, "."))-1])
		if err != nil {
			panic(err)
		}

		apiKeys := connector.APIKeys("")
		ocrResult, err := apiKeys.ProcessOCR(restyClient, filePath, strings.Split(fileUrl, ".")[1])
		if err != nil {
			panic(err)
		}

		if err = removeFile(*filePath); err != nil {
			panic(err)
		}

	}
	return nil
}

func saveSticker(response *resty.Response, filename, fileExt string) (*string, error) {
	filePath := fmt.Sprintf(`./pkg/temp_files/%s.%s`, filename, fileExt)

	// Create or open the file. If the file exists, it will be truncated.
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
		return nil, err
	}
	// Ensure the file is closed when the function exits
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Fatalf("Failed to close file: %v", closeErr)
		}
	}()

	// Write the byte slice to the file
	n, err := f.Write(response.Body())
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
		return nil, err
	}

	fmt.Printf("Successfully wrote %d bytes to %s.%s \n", n, filename, fileExt)

	return &filePath, nil
}

func removeFile(filePath string) error {
	// Attempt to remove the file
	return os.Remove(filePath)
}
