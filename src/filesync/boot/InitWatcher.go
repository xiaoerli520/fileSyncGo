package boot

import (
	"net/http"
	"fmt"
	"sinago/healthKeeper"
)

var hkw *healthKeeper.HealthKeeper

func HandleIndex(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<h1>File Sync Watcher</h1>")

	arrReady := hkw.List(Viper.GetInt("HealthKeeper.status.ready"))
	arrDead := hkw.List(Viper.GetInt("HealthKeeper.status.dead"))
	arrRec := hkw.List(Viper.GetInt("HealthKeeper.status.recovery"))
	arrUnhealth := hkw.List(Viper.GetInt("HealthKeeper.status.unhealth"))
	fmt.Fprintf(w, "<h2>Ready:</h2><br>")
	for k := range arrReady {
		fmt.Fprintf(w, k+"  ")
	}
	fmt.Fprintf(w, "<h2>Dead:</h2><br>")
	for k := range arrDead {
		fmt.Fprintf(w, k+"  ")
	}
	fmt.Fprintf(w, "<h2>Recovery:</h2><br>")
	for k := range arrRec {
		fmt.Fprintf(w, k+"  ")
	}
	fmt.Fprintf(w, "<h2>Unhealth:</h2><br>")
	for k := range arrUnhealth {
		fmt.Fprintf(w, k+"  ")
	}

}

func InitWatcher() {
	hkw = Hk
	fmt.Println("Web Watcher Inited")
	go func() {
		http.HandleFunc("/", HandleIndex)
		http.ListenAndServe(":7777", nil)
	}()
}
