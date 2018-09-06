package main

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "/root/.kube/config", "absolute path to the kubeconfig file")
	client     *kubernetes.Clientset
)

func GetPods(namespace, svcId string) []string {
	pods, err := client.Core().Pods(namespace).List(v1.ListOptions{
		LabelSelector: "svc_id=" + svcId,
	})
	if err != nil {
		panic(err.Error())
	}
	res := []string{}
	for _, pod := range pods.Items {
		res = append(res, pod.GetName())
	}
	fmt.Println("DEBUG: get pods ", res)
	return res
}

//cmd := exec.Command("kubectl", "exec", "--namespace=zcontainer", "php-web-30tcd", "grep", " ls", "/tmp/")
func runCommand(commandWaitGroup *sync.WaitGroup, commandOutputChan chan string, namespace, name string, command ...string) string {
	defer commandWaitGroup.Done()
	commands := make([]string, 0)
	commands = append(commands, namespace, name)
	//commands = append(commands, "exec", "--namespace="+namespace, name)
	for _, c := range command {
		commands = append(commands, c)
	}

	fmt.Println("./kubectl", commands)
	cmd := exec.Command("./kubectl.sh", commands...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	commandOutputChan <- string(output)

	return string(output)
}

func doRequest(w http.ResponseWriter, r *http.Request) {
	var commandOutputChan = make(chan string, 1024)
	var commandWaitGroup sync.WaitGroup

	if len(r.URL.Query()["namespace"]) == 0 {
		return
	}

	namespace := r.URL.Query()["namespace"][0]
	svcId := r.URL.Query()["svc_id"][0]
	commands := r.URL.Query()["commands"]

	pods := GetPods(namespace, svcId)

	res := ""
	for _, pod := range pods {
		commandWaitGroup.Add(1)
		go runCommand(&commandWaitGroup, commandOutputChan, namespace, pod, commands...)
	}
	commandWaitGroup.Wait()
	close(commandOutputChan)

	for s := range commandOutputChan {
		res += s
	}

	w.Write([]byte(res))
}

func main() {
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", http.HandlerFunc(doRequest))
	http.ListenAndServe(":8888", nil)
}
