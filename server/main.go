package main

import (
	"fmt"
	"github.com/albertocsm/ta-hackday-k8s-scrapper/server/cluster"
	"github.com/olivere/elastic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"time"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a ES client
	client, err := elastic.NewClient(
		elastic.SetURL("http://elasticsearch:9200"),
		elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		// Handle error
	}

	for {
		pods := cluster.Fetch(*clientset)
		fmt.Printf("There are %d pods in the cluster\n", len(*pods))

		cluster.Index(*client, *pods)

		time.Sleep(1 * time.Minute)
	}

}
