package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"gopkg.in/resty.v1"

	"github.com/hongkailiu/test-go/pkg/quay"
)

const (
	defaultTargetUrl = "https://web-hongkliu-run.apps.ca-central-1.starter.openshift-online.com/"
	addr             = ":8080"
)

var (
	upgrader = websocket.Upgrader{} // use default options
	ss       *ServiceStatus
	log      *logrus.Logger
	helper   k8sHelper
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

func webhook(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost != r.Method {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	//log.WithField("string(body", string(body)).Debug("post with body")
	go handle(body)
	content := "OK"
	n, err := fmt.Fprint(w, content)
	if err != nil {
		log.WithError(err).Error("cannot write response")
	}
	if n != len(content) {
		log.WithField("n", n).Errorf("cannot write '%s' properly", content)
	}
}

func handle(bytes []byte) {
	event := quay.RepositoryEvent{}
	if err := json.Unmarshal(bytes, &event); err != nil {
		log.WithError(err).WithField("string(bytes)", string(bytes)).Error("cannot unmarshal with json")
	}
	log.WithField("event", event).Debug("received an event")
	if err := applyDeployment(event); err != nil {
		log.WithError(err).WithField("event", event).Errorf("error occurred when applying deployment")
	}
}

func applyDeployment(event quay.RepositoryEvent) error {
	if len(event.UpdatedTags) == 0 {
		return fmt.Errorf("invalid event: empty tags")
	}
	tag := event.GetTheMostRecentTag()
	if !strings.Contains(tag, "testctl") {
		log.WithField("tag", tag).Warn("not testctl tag, ignoring ...")
		return nil
	}
	if !helper.inCluster {
		return fmt.Errorf("not in cluster")
	}
	client := helper.k8sClientSet.AppsV1().Deployments(helper.project)
	d, err := client.Get(helper.deployment, metav1.GetOptions{})
	if err != nil {
		return err
	}
	log.WithField("d.Name", d.Namespace).WithField("d.Name", d.Name).Debug("found deployment")
	found := false
	containers := d.Spec.Template.Spec.Containers
	for i, c := range containers {
		log.WithField("c.Name", c.Name).Debug("listing container ...")
		if c.Name == helper.container {
			log.WithField("c.Name", c.Name).Debug("found the matching container")
			found = true
			targetImage := fmt.Sprintf("%s:%s", event.DockerURL, tag)
			log.WithField("c.Image", c.Image).WithField("targetImage", targetImage).Debug("get the current image ...")
			if c.Image != targetImage {
				d.Spec.Template.Spec.Containers[i].Image = targetImage
				log.WithField("d.Name", d.Name).Debug("updating deployment ...")
				if _, err := client.Update(d); err != nil {
					log.WithError(err).WithField("d.Name", d.Name).Debug("error occurred when updating deployment")
					return err
				}
			}
		}
	}
	if !found {
		return fmt.Errorf("cannot find the container with name '%s'", helper.container)
	}
	return nil
}

func status(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("error occurred when getting web socket connection")
		return
	}
	defer func() {
		err := c.Close()
		if err != nil {
			log.WithError(err).Error("error occurred when closing web socket connection")
		}
	}()

	for {
		if ss != nil {
			b, err := json.Marshal(ss)
			if err != nil {
				log.WithError(err).Error("error occurred when marshaling service status")
				break
			}
			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.WithError(err).Error("error occurred when writing message")
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

type k8sHelper struct {
	inCluster    bool
	config       *rest.Config
	k8sClientSet *kubernetes.Clientset
	project      string
	deployment   string
	container    string
}

// Start status http server
func Start(logger *logrus.Logger) {
	log = logger

	helper = getK8SHelper()

	url := os.Getenv("target_url")
	if url == "" {
		url = defaultTargetUrl
	}

	ss = &ServiceStatus{Url: url}

	log.WithField("ss.url", ss.Url).Info("contacting web server ...")
	// fresh web server status
	go func() {
		log.Info("configure cron jobs ...")
		c := cron.New()
		err := c.AddFunc("*/10 * * * * *", func() {
			now := time.Now()
			log.WithField("now", now.Format(time.RFC3339)).Debug("Every 10 seconds ... ")
			updateStatus(now)
		})
		if err != nil {
			log.WithError(err).Error("error occurred when the add cron job")
		}
		c.Start()
	}()

	http.HandleFunc("/status", status)
	http.HandleFunc("/", home)
	http.HandleFunc("/webhook", webhook)
	log.Println("starting status http server ...")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getK8SHelper() k8sHelper {
	if os.Getenv("IN_CLUSTER") == "true" {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.WithError(err).Error("cannot get in cluster config")
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.WithError(err).Error("cannot get in cluster client set")
		}
		return k8sHelper{
			inCluster:    true,
			config:       config,
			k8sClientSet: clientset,
			project:      os.Getenv("PROJECT_NAME"),
			deployment:   os.Getenv("DEPLOYMENT_NAME"),
			container:    os.Getenv("CONTAINER_NAME"),
		}
	}
	return k8sHelper{}
}

func updateStatus(now time.Time) {
	resp, err := resty.R().Get(ss.Url)
	if err != nil {
		log.WithError(err).Error("error occurred when getting status from web server")
	}
	sc := resp.StatusCode()
	log.WithField("statusCode", sc).Debug("web server returned")
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
