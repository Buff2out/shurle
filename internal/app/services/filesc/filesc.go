package filesc

import (
	Event "github.com/Buff2out/shurle/internal/app/api/shortener"
	"github.com/Buff2out/shurle/internal/app/config/files"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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

func FillEvents(sugar *zap.SugaredLogger, file string, links map[string]string) map[string]string {
	//var events = make([]Event.ShURLFile, 0, 5) // мда, теперь это events атавизм
	if file == "" {
		file = filepath.Join(os.TempDir(), "short-url-db.json")
	}
	c, err := files.NewConsumer(file, sugar)
	if err != nil {
		sugar.Infow(
			"in fillEvents failed",
		)
	} else {
		for {
			sugar.Infow(
				"info about path of file", "file", file,
			)
			el, er := c.ReadEvent()
			if er != nil {
				sugar.Infow(
					"END OF FILE", "element", el,
				)
				break
			}
			//events = append(events, *el) // мда, теперь это events атавизм
			links[el.ShortURL] = el.OriginalURL
		}
	}
	return links
}
