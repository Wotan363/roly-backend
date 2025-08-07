package ai

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// This will make the prompts to ai
// TODO: do a proper implementation

var availableModels = map[string]string{
	"1": openai.ChatModelGPT4_1Nano, // 0.1$ - 0-4$ cost
	"2": openai.ChatModelGPT4_1Mini, // 0.4$ - 1.6$ cost
	"3": openai.ChatModelO3Mini,     // 1.1$ - 4.4$ cost
}

func selectModel(reader *bufio.Reader) string {
	fmt.Println("Wähle ein Modell:")
	fmt.Println("1. gpt-4.1-nano")
	fmt.Println("2. gpt-4.1-mini")
	fmt.Println("3. o3-mini")
	fmt.Print("Modellnummer > ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	model, ok := availableModels[choice]
	if !ok {
		fmt.Println("Ungültige Auswahl, Standard: gpt-4.1-nano")
		model = openai.ChatModelGPT4_1Nano
	}
	return model
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Fehler: OPENAI_API_KEY ist nicht gesetzt.")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🧠 ChatGPT CLI – mit Modellwahl und Konversation")
	fmt.Println("------------------------------------------------")

	// Modell wählen
	model := selectModel(reader)

	// System-Rolle setzen
	fmt.Print("Welche Rolle soll der Assistent übernehmen? > ")
	systemPrompt, _ := reader.ReadString('\n')
	systemPrompt = strings.TrimSpace(systemPrompt)

	// Nachrichtenverlauf vorbereiten
	var messages []openai.ChatCompletionMessageParamUnion
	messages = append(messages, openai.SystemMessage(systemPrompt))

	fmt.Println("\nStarte die Konversation! (Tippe 'exit' zum Beenden)")
	fmt.Println("-----------------------------------------------------")

	for {
		fmt.Print("\nDu > ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if strings.EqualFold(userInput, "exit") {
			fmt.Println("🛑 Konversation beendet.")
			break
		}

		// Nutzerfrage hinzufügen
		messages = append(messages, openai.UserMessage(userInput))

		// Anfrage an OpenAI senden
		resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
			Model:    model,
			Messages: messages,
		})
		if err != nil {
			log.Fatalf("API-Fehler: %v", err)
		}

		if len(resp.Choices) > 0 {
			reply := resp.Choices[0].Message.Content
			fmt.Printf("\n🧠 Modell > %s\n", reply)

			// Antwort auch zum Verlauf hinzufügen
			messages = append(messages, openai.AssistantMessage(resp.Choices[0].Message.Content))
		} else {
			fmt.Println("⚠️ Keine Antwort erhalten.")
		}
	}
}
