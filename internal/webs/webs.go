package webs

import (
	"courseproject/internal/storages"
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"time"
)

func SendNewLots(client *websocket.Conn, db storages.INTT, tpl map[string]*template.Template){
	ticker := time.NewTicker(3 * time.Second)
	for {
		w, err:= client.NextWriter(websocket.TextMessage)
		if err != nil{
			ticker.Stop()
			break
		}

		lts, _ := db.GetLots("")
		//tmpl.RenderTemplate(w, "index", "base", lts , tpl)
		msg, err:= json.Marshal(lts)

		w.Write(msg)
		w.Close()

		<-ticker.C
	}
}