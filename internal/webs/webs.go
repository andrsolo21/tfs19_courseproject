package webs

import (
	"courseproject/internal/storages"
	"encoding/json"
	"html/template"
	"time"

	"github.com/gorilla/websocket"
)

func SendNewLots(client *websocket.Conn, db storages.INTT, tpl map[string]*template.Template) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		w, err := client.NextWriter(websocket.TextMessage)
		if err != nil {
			ticker.Stop()
			break
		}

		lts, _ := db.GetLots("")
		//tmpl.RenderTemplate(w, "index", "base", lts , tpl)
		msg, _ := json.Marshal(lts)

		_, _ = w.Write(msg)
		w.Close()

		<-ticker.C
	}
}
