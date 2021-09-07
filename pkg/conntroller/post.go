package conntroller

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func PromPost(query string) (body []byte) {
	url := "http://prometheus.kube-system:9090/api/v1/query"
	//query :="sum(container_memory_working_set_bytes{container_name!=\"\",image!=\"google_containers/pause-amd64:3.0\",namespace=\"daas\",pod_name=\"ys-bi-app-api-extract-6c9ff66fc7-xqrhd\"})by(namespace,pod_name) /1024/1024 "
	request, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("query="+query))
	if err != nil {
		log.Println(err)
	}
	defer request.Body.Close()
	body, err = ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
	}
	return
}
