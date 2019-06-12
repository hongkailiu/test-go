package server

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"gopkg.in/resty.v1"
)

const (
	defaultTargetUrl = "https://web-hongkliu-run.b542.starter-us-east-2a.openshiftapps.com/"
)

var (
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{} // use default options
	ss       *ServiceStatus
)

// Check ...
type Check struct {
	Time   time.Time `json:"time"`
	Status int       `json:"status"`
}

// ServiceStatus ...
type ServiceStatus struct {
	Url                 string `json:"url"`
	LastSuccessfulCheck *Check `json:"lastSuccessfulCheck,omitempty"`
	LastFailedCheck     *Check `json:"lastFailedCheck,omitempty"`
}

func status(w http.ResponseWriter, r *http.Request) {
	log.Println("echo: a")
	c, err := upgrader.Upgrade(w, r, nil)
	log.Println("echo: b")
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		err := c.Close()
		if err != nil {
			log.Print("close:", err)
		}
	}()

	for {
		log.Println("for loop:")
		//mt, message, err := c.ReadMessage()
		//if err != nil {
		//	log.Println("read:", err)
		//	break
		//}
		//log.Printf("recv: %s", message)
		if ss != nil {
			log.Println("sending:")
			b, err := json.Marshal(ss)
			if err != nil {
				log.Println("marshal:", err)
				break
			}
			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	err := homeTemplate.Execute(w, "ws://"+r.Host+"/status")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

// Start status http server
func Start() {
	flag.Parse()
	log.SetFlags(0)

	url := os.Getenv("target_url")
	if url == "" {
		url = defaultTargetUrl
	}

	ss = &ServiceStatus{Url: url}

	log.Println("url: " + ss.Url)

	log.Println("configure cron jobs ...")
	c := cron.New()
	//log.Println("000")
	err := c.AddFunc("*/10 * * * * *", func() {
		now := time.Now()
		log.Println("Every 10 seconds" + now.Format(time.RFC3339))
		updateStatus(now)
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()

	http.HandleFunc("/status", status)
	http.HandleFunc("/", home)
	log.Println("starting status http server ...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func updateStatus(now time.Time) {
	resp, err := resty.R().Get(ss.Url)
	if err != nil {
		log.Print("resty: ", err)
	}
	sc := resp.StatusCode()
	log.Printf("sc is %d\n", sc)
	check := &Check{now, sc}
	if sc >= http.StatusOK && sc < http.StatusMultipleChoices {
		ss.LastSuccessfulCheck = check
	} else {
		ss.LastFailedCheck = check
	}
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
