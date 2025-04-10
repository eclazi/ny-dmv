package main

import (
	"dmv-ny/pkg/dmvapi"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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

var watchState *WatchState

type WatchState struct {
	serviceId      int
	locationIds    []int
	cancel         chan struct{}
	timerChan      chan struct{}
	checkPeriod    time.Duration
	withinDuration time.Duration

	chatId int64
	bot    *tgbotapi.BotAPI
}

func NewWatchState(serviceId int, locationIds []int, withinDuration time.Duration, chatId int64, bot *tgbotapi.BotAPI, checkPeriod time.Duration) *WatchState {
	return &WatchState{
		serviceId:      serviceId,
		locationIds:    locationIds,
		cancel:         make(chan struct{}),
		timerChan:      make(chan struct{}),
		withinDuration: withinDuration,
		checkPeriod:    checkPeriod,
		chatId:         chatId,
		bot:            bot,
	}
}

// Get sorted appointments for a specific service and set of locations
func getAppointments(serviceId int, locationIds []int) ([]dmvapi.Appointment, error) {
	appointments := []dmvapi.Appointment{}
	for _, locationId := range locationIds {
		locAppointments, err := dmvClient.GetAppointments(locationId, serviceId)
		if err != nil {
			return nil, fmt.Errorf("error getting appointments for location %d: %w", locationId, err)
		}
		appointments = append(appointments, locAppointments...)
	}

	// Sort appointments by date (soonest first)
	sort.Slice(appointments, func(i, j int) bool {
		return appointments[i].DateTime.Before(appointments[j].DateTime)
	})

	return appointments, nil
}

// Filter appointments within the specified duration
func filterAppointments(appointments []dmvapi.Appointment, now time.Time, withinDuration time.Duration) []dmvapi.Appointment {
	var filtered []dmvapi.Appointment
	for _, appointment := range appointments {
		if appointment.DateTime.Sub(now) <= withinDuration && appointment.DateTime.After(now) {
			filtered = append(filtered, appointment)
		}
	}
	return filtered
}

func (ws *WatchState) Start() {
	go func() {
		ticker := time.NewTicker(ws.checkPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Check for new appointments
				appointments, err := getAppointments(ws.serviceId, ws.locationIds)
				if err != nil {
					log.Printf("Error getting appointments: %v", err)
					continue
				}
				if len(appointments) == 0 {
					log.Println("No appointments found.")
					continue
				}
				// Filter appointments within the specified duration
				filteredAppointments := filterAppointments(appointments, time.Now(), ws.withinDuration)

				// Send message to Telegram bot
				if len(filteredAppointments) > 0 {
					result := strings.Builder{}
					for index, appointment := range filteredAppointments {

						if index < 5 {
							result.WriteString(fmt.Sprintf("Appointment ID: %d\tDate: %s\n", appointment.SlotId, appointment.DateTime.Format(time.RFC3339)))
						} else {
							break
						}
					}
					msg := tgbotapi.NewMessage(ws.chatId, result.String())
					if _, err := ws.bot.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
				} else {
					log.Println("No appointments found within the specified duration.")
				}
			case <-ws.cancel:
				log.Println("Stopping watch...")
				return
			}
		}
	}()
}

func (ws *WatchState) Stop() {
	// Stop the watch
	ws.cancel <- struct{}{}
	log.Println("Watch stopped.")
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
	case "watch":
		// Get the service ID and location ID from the command arguments
		args := update.Message.CommandArguments()
		if args == "" {
			msg.Text = "Please provide a service ID and location ID to watch for appointments."
		} else {
			ids := strings.Split(args, " ")
			if len(ids) != 2 {
				msg.Text = "Please provide both a service ID and a location ID."
			} else {
				serviceId, err := strconv.Atoi(ids[0])
				if err != nil {
					msg.Text = "Invalid service ID. Please provide a valid number."
				} else {
					locationId, err := strconv.Atoi(ids[1])
					if err != nil {
						msg.Text = "Invalid location ID. Please provide a valid number."
					} else {

						if watchState != nil {
							watchState.Stop()
							watchState = nil
						}

						// Create a new watch state
						// Start watching for appointments
						withinDuration := 24 * 100 * time.Hour // Set the duration to watch for appointments
						watchState = NewWatchState(serviceId, []int{locationId}, withinDuration, update.Message.Chat.ID, bot, 20*time.Second)
						watchState.Start()
						msg.Text = fmt.Sprintf("Watching for appointments for service ID %d at location ID %d.", serviceId, locationId)
					}
				}
			}
		}
	case "stop":
		if watchState != nil {
			watchState.Stop()
			watchState = nil
		}
		// Stop watching for appointments

	default:
		msg.Text = "Unknown command. Type /help to see available commands."
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
