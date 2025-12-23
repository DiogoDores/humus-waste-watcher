package handlers

import (
	"context"
	"fmt"
	"log"

	"src/formatters"
	repo "src/repository"

	tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler func(ctx context.Context, bot *tg_bot.BotAPI, repo repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) error

// HandleMyPoopLog handles the /my_poop_log command
func HandleMyPoopLog(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) error {
	globalPoopCount, errGlobal := r.GetGlobalPoopCount(ctx, userId)
	monthlyPoopCounts, errMonthly := r.GetMonthlyPoopStats(ctx, userId)
	daysWithoutPoop, errNoPoop := r.GetDaysWithoutPoop(ctx, userId)
	maxStreak, errStreak := r.GetMaxPoopStreak(ctx, userId)
	day, poops, mostPoopsErr := r.GetDayWithMostPoops(ctx, userId)

	if errGlobal != nil || errMonthly != nil || errNoPoop != nil || errStreak != nil || mostPoopsErr != nil {
		msg.Text = "Sorry, I couldn't retrieve your poop log\\. Please try again later\\!"
		_, err := bot.Send(msg)
		return err
	}

	username := update.Message.From.UserName
	msg.Text = formatters.FormatPoopLog(username, globalPoopCount, monthlyPoopCounts, daysWithoutPoop, maxStreak, day, poops)
	_, err := bot.Send(msg)
	return err
}

// HandleLeaderboard handles the /leaderboard command
func HandleLeaderboard(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) error {
	monthlyLeaderboard, err := r.GetMonthlyLeaderboard(ctx)
	if err != nil {
		msg.Text = "Sorry, I couldn't retrieve the monthly leaderboard\\. Please try again later\\!"
		_, sendErr := bot.Send(msg)
		if sendErr != nil {
			return fmt.Errorf("failed to send error message: %w", sendErr)
		}
		return err
	}

	msg.Text = formatters.FormatLeaderboard(monthlyLeaderboard)
	_, err = bot.Send(msg)
	return err
}

// HandleBottomPoopers handles the /bottom_poopers command
func HandleBottomPoopers(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) error {
	bottomPoopers, err := r.GetBottomPoopers(ctx)
	if err != nil {
		msg.Text = "Sorry, I couldn't retrieve the bottom poopers\\. Please try again later\\!"
		_, sendErr := bot.Send(msg)
		if sendErr != nil {
			return fmt.Errorf("failed to send error message: %w", sendErr)
		}
		return err
	}

	msg.Text = "This month's bottom poopers are:\n" + formatters.BuildPoodiumMessage(bottomPoopers)
	_, err = bot.Send(msg)
	return err
}

// HandlePoodium handles the /poodium command
func HandlePoodium(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) error {
	monthlyPoodium, err := r.GetMonthlyPoodium(ctx)
	if err != nil {
		msg.Text = "Sorry, I couldn't retrieve the monthly poodium\\. Please try again later\\!"
		_, sendErr := bot.Send(msg)
		if sendErr != nil {
			return fmt.Errorf("failed to send error message: %w", sendErr)
		}
		return err
	}

	msg.Text = "This month's top poopers are:\n" + formatters.BuildPoodiumMessage(monthlyPoodium)
	_, err = bot.Send(msg)
	return err
}

// HandleYearlyPoodium handles the /poodium_year command
func HandleYearlyPoodium(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) error {
	yearlyPoodium, err := r.GetYearlyPoodium(ctx)
	if err != nil {
		msg.Text = "Sorry, I couldn't retrieve the yearly poodium\\. Please try again later\\!"
		_, sendErr := bot.Send(msg)
		if sendErr != nil {
			return fmt.Errorf("failed to send error message: %w", sendErr)
		}
		return err
	}

	msg.Text = "This year's top poopers are:\n" + formatters.BuildPoodiumMessage(yearlyPoodium)
	_, err = bot.Send(msg)
	return err
}

// HandleHelp handles the /help command and unknown commands
func HandleHelp(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) error {
	isUnknownCommand := update.Message.Command() != "help"
	msg.Text = formatters.FormatHelpMessage(isUnknownCommand)
	_, err := bot.Send(msg)
	return err
}

func GetCommandHandlers() map[string]CommandHandler {
	return map[string]CommandHandler{
		"my_poop_log":    HandleMyPoopLog,
		"leaderboard":    HandleLeaderboard,
		"bottom_poopers": HandleBottomPoopers,
		"poodium":        HandlePoodium,
		"poodium_year":   HandleYearlyPoodium,
		"help":           HandleHelp,
	}
}

// HandleCommand routes commands to their respective handlers
func HandleCommand(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) {
	log.Println("Command received:", update.Message.Command())

	handlers := GetCommandHandlers()
	command := update.Message.Command()

	handler, exists := handlers[command]
	if !exists {
		handler = HandleHelp
	}

	if err := handler(ctx, bot, r, update, userId, msg); err != nil {
		log.Printf("Error handling command %s: %v", command, err)
	}
}
