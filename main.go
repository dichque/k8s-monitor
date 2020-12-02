package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {

	sysNamespace := "kube-system"
	dnsTarget := "www.google.com"
	endPointTarget := "http://www.google.com/robots.txt"
	crashLoop := false
	notRunning := false
	dnsFailure := false
	endPointFailure := false

	log.Println("Program to monitor pods in ", sysNamespace, " namespace")
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		pods, err := clientset.CoreV1().Pods(sysNamespace).List(metav1.ListOptions{})
		if err != nil {
			notify("API service is not available: " + err.Error())
			panic(err.Error())
		}

		for _, item := range pods.Items {
			pod, err := clientset.CoreV1().Pods(sysNamespace).Get(item.ObjectMeta.Name, metav1.GetOptions{})
			if err != nil {
				notify("API service is not available: " + err.Error())
			}

			// log.Println("Pod: ", pod.Name, "Status: ", pod.Status.Phase, "RestartCount: ", pod.Status.ContainerStatuses[0].RestartCount)

			if !isDNSUp(dnsTarget) {
				dnsFailure = true
			}

			if pod.Status.Phase != "Running" {
				notRunning = true
			}

			if !isendpointUP(endPointTarget) {
				endPointFailure = true
			}

			// if pod.Status.ContainerStatuses[0].RestartCount > 0 {
			//   crashLoop = true
			//}

			if notRunning || crashLoop {
				notify(pod.Name)
				crashLoop = false
				notRunning = false
				break
			} else if dnsFailure {
				notify("DNS service is down")
				dnsFailure = false
				break
			} else if endPointFailure {
				notify("Endpoint " + endPointTarget + " is not reachable")
				endPointFailure = false
				break
			} else {
				log.Debug("All is well!!")
			}
		}

		time.Sleep(4 * time.Second)
	}
}

func notify(msg string) {
	log.Error("scream!! critical component is down: ", msg)
}

func isDNSUp(host string) bool {
	_, err := net.LookupHost(host)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func isendpointUP(url string) bool {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
