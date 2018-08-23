package main

import (
	"fmt"
	//"bytes"
	"net/http"
	//"os"
	"os/exec"
	"strings"
)

func GetPods(namespace, svcId string) []string {
	cmd := exec.Command("./get_pod.sh", namespace, svcId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("DEBUG: get output", string(output))
	res := strings.Split(string(output), " ")

	trimed := make([]string, 0)
	for _, r := range res {
		if strings.TrimSpace(r) != "" {
			trimed = append(trimed, strings.TrimSpace(r))
		}
	}
	fmt.Println("DEBUG: get pods ", trimed)
	return trimed
}

//cmd := exec.Command("kubectl", "exec", "--namespace=zcontainer", "php-web-30tcd", "grep", " ls", "/tmp/")
func runCommand(namespace, name string, command ...string) string {
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

	return string(output)
}

func doRequest(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query()["namespace"][0]
	svcId := r.URL.Query()["svc_id"][0]
	commands := r.URL.Query()["commands"]

	pods := GetPods(namespace, svcId)

	res := ""
	for _, pod := range pods {
		res += runCommand(namespace, pod, commands...)
	}

	w.Write([]byte(res))
}

func main() {

	cmd := exec.Command("kubectl", []string{"exec", "--namespace=fortest", "test-ok--jetty-2604245236-30wh8", "grep", "/tmp", "/tmp/tmp"}...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output), err)
	http.HandleFunc("/", http.HandlerFunc(doRequest))
	http.ListenAndServe(":8888", nil)
}
