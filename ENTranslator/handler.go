package function

import (
	//"net/http"
	"cloud.google.com/go/translate"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"os"
	"strings"
)

type dockerHubStatsType struct {
	Count int `json:"count"`
}

func sanitizeInput(input string) string {
	parts := strings.Split(input, "\n")
	return strings.Trim(parts[0], " ")
}

func translateTextToEnglish(text string) (string, error) {
	ctx := context.Background()

	targetLanguage := "En"

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", err
	}
	authOpt := option.WithAPIKey(os.Getenv("APIKEY"))
	client, err := translate.NewClient(ctx, authOpt)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", err
	}
	return resp[0].Text, nil
}

// Handle a serverless request
func Handle(req []byte) string {

	if len(req) == 0 {
		fmt.Fprintf(os.Stderr, "A word or a sentance required to detect language")
		return fmt.Sprintf("A word or a sentance required to detect language")
	}

	arg := string(req)

	arg = sanitizeInput(arg)
	translated, terr := translateTextToEnglish(arg)
	if terr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", terr)
		return fmt.Sprintf("Error: %v", terr)
	}

	return translated
}
