package config

import (
	"log"
	"os"
	"sync"

	"github.com/Atgoogat/openmensarobot/domain"
	"github.com/Atgoogat/openmensarobot/emojifier"
)

var (
	formatter     []domain.TextFormatter
	formatterOnce sync.Once
)

func GetTextFormatter() []domain.TextFormatter {
	formatterOnce.Do(func() {
		formatter = make([]domain.TextFormatter, 0)

		url, ok := os.LookupEnv(EMOJIFIER_URL)
		if ok {
			log.Printf("usage of emojifier at %s", url)
			formatter = append(formatter, emojifier.NewEmojifierTextFormatter(url))
		} else {
			log.Println("no emojifier setup")
		}
	})

	return formatter
}
