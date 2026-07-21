package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Billy19191/Telegram-Morpho-Bot/model"
	"github.com/Billy19191/Telegram-Morpho-Bot/service"
	"github.com/Billy19191/Telegram-Morpho-Bot/util"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	// "github.com/joho/godotenv"
)

var accountableService *service.AccountableService
var morphoService *service.MorphoService

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	panic("ENV not found")
	// }

	// _, err := config.LoadConfig()
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to load config: %v", err))
	// }

	tgBotToken := getEnvKey("TG_BOT_TOKEN", "")
	accountableService = service.NewAccountableService(
		getEnvKey("BASE_URL", ""),
		getEnvKey("WALLET_ADDRESS", ""),
		getEnvKey("CHAIN_ID", ""),
	)
	morphoService = service.NewMorphoService(
		getEnvKey("BASE_URL", ""),
		getEnvKey("WALLET_ADDRESS", ""),
		getEnvKey("CHAIN_ID", ""),
	)

	chatID, err := strconv.ParseInt(getEnvKey("TG_CHAT_ID", "0"), 10, 64)
	if err != nil || chatID == 0 {
		panic("TG_CHAT_ID is required and must be a valid number")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	tgBot, err := bot.New(tgBotToken, opts...)

	if err != nil {
		panic(err)
	}

	go startCronMonitor(ctx, tgBot, chatID)

	log.Println("🤖 Bot started. Cron monitor running every 5 minutes.")
	tgBot.Start(ctx)
}

func getEnvKey(key string, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallbackValue
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("User message: " + update.Message.Text)

	accountableResult, err := accountableService.GetBorrowPositions()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}

	morphoResult, err := morphoService.GetBorrowPosition()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}

	if len(accountableResult.ResponseData.VaultAllocations) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "No accountable positions were returned by the API.",
		})
		return
	}

	if len(morphoResult.ResponseData) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "No vault positions were returned by the API.",
		})
		return
	}

	accountableAllocation := accountableResult.ResponseData.VaultAllocations[0]
	morphoVault := morphoResult.ResponseData[0]
	riskReport := service.EvaluateVaultRisk(accountableAllocation, morphoVault)
	msg := formatVaultMessage(accountableAllocation, morphoVault, riskReport)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

// startCronMonitor checks vault positions every 5 minutes.
// - Critical status → sends alert immediately
// - Normal status → sends routine report every 8 hours
func startCronMonitor(ctx context.Context, b *bot.Bot, chatID int64) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	lastRoutineReport := time.Now()

	checkAndNotify(ctx, b, chatID, &lastRoutineReport)

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Cron monitor stopped.")
			return
		case <-ticker.C:
			checkAndNotify(ctx, b, chatID, &lastRoutineReport)
		}
	}
}

func checkAndNotify(ctx context.Context, b *bot.Bot, chatID int64, lastRoutineReport *time.Time) {
	accountableResult, err := accountableService.GetBorrowPositions()
	if err != nil {
		log.Printf("❌ Cron check failed: %v", err)
		return
	}

	morphoResult, err := morphoService.GetBorrowPosition()
	if err != nil {
		log.Printf("❌ Cron check failed: %v", err)
		return
	}

	if len(accountableResult.ResponseData.VaultAllocations) == 0 {
		log.Println("❌ Cron check returned no accountable positions.")
		return
	}

	if len(morphoResult.ResponseData) == 0 {
		log.Println("❌ Cron check returned no vault positions.")
		return
	}

	accountableAllocation := accountableResult.ResponseData.VaultAllocations[0]
	morphoVault := morphoResult.ResponseData[0]
	riskReport := service.EvaluateVaultRisk(accountableAllocation, morphoVault)
	log.Printf("📋 Cron check — Status: %s", riskReport.OverallStatus)

	if riskReport.OverallStatus == model.StatusCritical {
		msg := "🚨 CRITICAL ALERT 🚨\n" + formatVaultMessage(accountableAllocation, morphoVault, riskReport)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
		})
		*lastRoutineReport = time.Now()
		log.Println("🚨 Critical alert sent!")
		return
	}

	if time.Since(*lastRoutineReport) >= 8*time.Hour {
		msg := "📊 Routine Monitor Report\n\n" + formatVaultMessage(accountableAllocation, morphoVault, riskReport)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
		})
		*lastRoutineReport = time.Now()
		log.Println("📊 Routine report sent!")
	}
}

func formatVaultMessage(accountable model.AccountableVaultAllocationEntity, morpho model.VaultEntity, riskReport model.RiskReport) string {
	netApy := calculateNetApy(accountable, morpho)
	netPnl := calculateNetPnl(accountable, morpho)

	return fmt.Sprintf(
		"⚡ Portfolio Summary\n\n"+
			"📊 Status: %s\n"+
			"📈 Net APY: %s%%\n"+
			"💰 Net PNL (USD): $%s\n"+
			"💧 Net Asset (USD): $%s\n"+
			"----------------------\n"+
			"🏦 Accountable\n\n"+
			"📝 Name: %s\n"+
			"📈 Deposit APY: %s%%\n"+
			"💰 Deposit PNL (USD): $%s\n"+
			"----------------------\n"+
			"🏛️ Morpho\n\n"+
			"📝 Name: %s\n"+
			"❤️ Health Factor: %s\n"+
			"📉 Borrow APY: %s%% (Net)\n"+
			"💰 Borrow PNL (USD): $%s\n",
		riskReport.OverallStatus,
		util.FormatNumberWithSeparator(netApy),
		util.FormatNumberWithSeparator(netPnl),
		util.FormatNumberWithSeparator(accountable.Value-morpho.BorrowAssetsUsd),
		accountable.VaultName,
		util.FormatNumberWithSeparator(accountable.Apy),
		util.FormatNumberWithSeparator(accountable.UnrealizedPnl),
		morpho.Name,
		util.FormatNumberWithSeparator(morpho.HealthFactor),
		util.FormatNumberWithSeparator(morpho.NetBorrowApy),
		util.FormatNumberWithSeparator(morpho.BorrowPnlUsd),
	)
}

func calculateNetApy(accountable model.AccountableVaultAllocationEntity, morpho model.VaultEntity) float64 {
	denominator := accountable.MyDepositUsd - morpho.BorrowAssetsUsd
	if denominator == 0 {
		return 0
	}

	return (accountable.Apy*accountable.Value - morpho.NetBorrowApy*morpho.BorrowAssetsUsd) / denominator
}

func calculateNetPnl(accountable model.AccountableVaultAllocationEntity, morpho model.VaultEntity) float64 {
	return accountable.UnrealizedPnl + morpho.BorrowPnlUsd
}
