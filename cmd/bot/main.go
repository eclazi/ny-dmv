package main

import (
	"dmv-ny/pkg/dmvapi"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var dmvClient *dmvapi.Client = dmvapi.NewClient()

func main() {
	// Get bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	// Create a new bot instance
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set up an update configuration
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60 // Timeout in seconds

	// Get updates channel
	updates := bot.GetUpdatesChan(updateConfig)

	// Process incoming updates
	for update := range updates {
		// Ignore non-message updates
		if update.Message == nil {
			continue
		}

		// Log received message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Handle commands (messages starting with '/')
		if update.Message.IsCommand() {
			handleCommand(bot, update)
			continue
		}
	}
}

func listServices() string {
	// List available services
	services, err := dmvClient.GetServices()
	if err != nil {
		log.Printf("Error listing services: %v", err)
		return "Error retrieving services."
	}

	result := strings.Builder{}
	for _, service := range services {
		result.WriteString(fmt.Sprintf("Service ID: %d\tName: %s\n", service.Id, service.Name))
	}
	return result.String()
}

func listLocations(serviceId int) string {
	// List available locations for a specific service
	locations, err := dmvClient.GetLocations(serviceId)
	if err != nil {
		log.Printf("Error listing locations: %v", err)
		return "Error retrieving locations."
	}

	result := strings.Builder{}
	for _, location := range locations {
		result.WriteString(fmt.Sprintf("Location ID: %d\tName: %s\n", location.Id, location.Name))
	}
	return result.String()
}

// handleCommand processes bot commands
func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	
	// Get the command without the '/'
	command := update.Message.Command()
	
	switch strings.ToLower(command) {
	case "start":
		msg.Text = "Welcome! I'm a Go Telegram bot. Use /help to see available commands."
	case "help":
		msg.Text = "Available commands:\n" +
			"/services - List available services\n" +
			"/locations <service_id> - List locations for a specific service\n" +
			"/help - Show this help message\n"
	case "echo":
		// Get arguments after the command
		args := update.Message.CommandArguments()
		if args == "" {
			msg.Text = "Please provide some text to echo."
		} else {
			msg.Text = args
		}
	case "services":
		// List available services
		msg.Text = listServices()
	case "locations":
		// Get the service ID from the command arguments
		args := update.Message.CommandArguments()
		if args == "" {
			msg.Text = "Please provide a service ID to list locations."
		} else {
			serviceId, err := strconv.Atoi(args)
			if err != nil {
				msg.Text = "Invalid service ID. Please provide a valid number."
			} else {
				msg.Text = listLocations(serviceId)
			}
		}
	default:
		msg.Text = "Unknown command. Type /help to see available commands."
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
