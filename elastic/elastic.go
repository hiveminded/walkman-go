package elastic

import (
	"log"

	"fmt"
	"gopkg.in/olivere/elastic.v3"
)

var (
	client *elastic.Client
)

type DocumentInterface interface {
	GetID() string
	GetType() string
}

func init() {
	var err error
	// Create a client
	client, err = elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateIndex(index string) error {
	// Create an index
	_, err := client.CreateIndex(index).Do()
	return err

	//if !createIndex.Acknowledged {
	//	// Not acknowledged
	//}
}

func Add(index string, document DocumentInterface) error {
	_, err := client.Index().
		Index(index).
		Type(document.GetType()).
		Id(document.GetID()).
		BodyJson(document).
		Refresh(true).
		Do()
	return err
}

func Get(index, typ, id string) (DocumentInterface, error) {
	// Get tweet with specified ID
	get1, err := client.Get().
		Index(index).
		Type(typ).
		Id(id).
		Do()
	if err != nil {
		return nil, err
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}
}

func IndexExists(index string) bool {
	exists, err := client.IndexExists(index).Do()
	if err != nil {
		log.Fatal(err)
	}

	return exists
}
