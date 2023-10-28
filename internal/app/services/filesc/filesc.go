package filesc

import (
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/config/files"
	"go.uber.org/zap"
)

func AddNote(sugar *zap.SugaredLogger, event Event.ShURLFile, filename string) {
	if filename != "" {
		p, er := files.NewProducer(filename, sugar)
		if er != nil {
			sugar.Infow("In MWPostAPIURL AddNote func under event var. Invalid path to file.")
		}
		er = p.WriteEvent(&event)
		if er != nil {
			sugar.Infow("In MWPostAPIURL AddNote func under WriteEvent. Cant Write to file.")
		}
	}
}
