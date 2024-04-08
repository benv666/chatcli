package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"benvbin/pkg/api"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	prompt         string
	inputTokens    uint64
	outputTokens   uint64
	responseTimeMS uint64
)

func init() {
	inputTokens = 0
	outputTokens = 0
	responseTimeMS = 0
}

func main() {
	model := flag.String("model", "anthropic.claude-v2", "The model to use")
	flag.Parse()

	quit := false
	lastres := ""

	for ok := true; ok; ok = !quit {
		form := huh.NewForm(
			huh.NewGroup(huh.NewNote().
				Title("Claude Chat").
				Description("Welcome to _Claude™_.\n\nGive your prompt?"),
				huh.NewNote().
					Title("Claude Says").
					Description(lastres),
				huh.NewText().
					Title("Prompt please, empty to quit").
					Value(&prompt),
			),
		)
		err := form.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if len(prompt) == 0 {
			quit = true
		} else {
			askClaude := func() {
				output, err := api.MakeRequest(prompt, *model)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				lastres = output.Response
				inputTokens += output.PromptTokenCount
				outputTokens += output.TokensConsumed
				responseTimeMS += output.ResponseTimeMS
			}
			_ = spinner.New().Title("Claude is thinking...").Action(askClaude).Run()
			prompt = ""
		}
	}

	// Usage summary.
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb, "%s\n\nYou've used up %s input and %s output tokens. \nIt took claude %s ms to handle this.\n\nHope it was worth it!\n",
		lipgloss.NewStyle().Bold(true).Render("CLAUDE RECEIPT"),
		keyword(strconv.FormatUint(inputTokens, 10)),
		keyword(strconv.FormatUint(outputTokens, 10)),
		keyword(strconv.FormatUint(responseTimeMS, 10)),
	)

	fmt.Println(
		lipgloss.NewStyle().
			Width(40).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(sb.String()),
	)
}
