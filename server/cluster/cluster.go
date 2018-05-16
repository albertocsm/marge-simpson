package cluster

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type PodDoc struct {
	Name          string    `json:"Name"`
	CanonicalName string    `json:"CanonicalName"`
	Timestamp     time.Time `json:"Timestamp"`
	Version       string    `json:"Version"`
	Namespace     string    `json:"Namespace"`
	//podLabels        map[string]string
}

type PodDocCollection struct {
	pods []PodDoc `json:"pod"`
}

func Fetch(clientset kubernetes.Clientset) *[]PodDoc {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	allPods := new(PodDocCollection)
	for i := 0; i < len(pods.Items); i++ {
		allPods.pods = append(
			allPods.pods,
			PodDoc{
				CanonicalName: cleanPodName(pods.Items[i].GenerateName),
				Name:          pods.Items[i].Name,
				Timestamp:     time.Unix(time.Now().Unix(), 0),
				Version:       pods.Items[i].Spec.Containers[0].Image,
				Namespace:     pods.Items[i].Namespace,
			})
	}
	return &allPods.pods
}

func Index(client elastic.Client, pods []PodDoc) {

	exists, err := client.IndexExists("pods_idx").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err = client.
			CreateIndex("pods_idx").
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
	}
	//else {
	//
	//	_, err := client.DeleteIndex("pods_idx").Do(context.Background())
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

	// Add a document to the index
	for _, pod := range pods {

		_, err = client.Index().
			Index("pods_idx").
			Type("pods").
			BodyJson(pod).
			Do(context.Background())

		if err != nil {
			// Handle error
			panic(err)
		}

	}
	_, err = client.Flush().Index("pods_idx").Do(context.TODO())
	if err != nil {
		panic(err)
	}
}

func cleanPodName(str string) string {

	fmt.Println("1-" + str)
	str = strings.TrimRight(str, "-")
	fmt.Println("2-" + str)
	stop := strings.LastIndex(str, "-")
	fmt.Println("3-" + string(stop))
	var val int
	if stop == -1 {
		val = len(str)
	} else {
		val = stop
	}
	str = str[0:val]
	fmt.Println("4-" + str)
	return str
}
