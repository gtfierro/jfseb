package main

import (
	"flag"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	user := os.Getenv("ARONNAXTESTUSER")
	pass := os.Getenv("ARONNAXTESTPASS")
	dbname := os.Getenv("ARONNAXTESTDB")
	backend := newBackend(user, pass, dbname)
	uuid1, _ := uuid.FromString("2b365d6a-8cbd-11e5-8bb3-0cc47a0f7eea")
	uuid2, _ := uuid.FromString("370dd17c-8cbd-11e5-8bb3-0cc47a0f7eea")
	uuid3, _ := uuid.FromString("3a77a0e0-8cbd-11e5-8bb3-0cc47a0f7eea")
	uuid4, _ := uuid.FromString("3da1cafc-8cbd-11e5-8bb3-0cc47a0f7eea")
	uuid5, _ := uuid.FromString("411ce89c-8cbd-11e5-8bb3-0cc47a0f7eea")
	backend.RemoveData()

	initTags := map[string]string{
		"Location/City":            "Berkeley",
		"Location/Building":        "Soda",
		"Location/Floor":           "4",
		"Location/Room":            "410",
		"Properties/Timezone":      "America/Los_Angeles",
		"Properties/ReadingType":   "double",
		"Properties/UnitofMeasure": "F",
		"Properties/UnitofTime":    "ms",
		"Properties/StreamType":    "numeric",
	}

	temperatureTags := map[string]string{
		"Metadata/Point/Type":   "Sensor",
		"Metadata/Point/Sensor": "Temperature",
	}

	for _, doc := range []Document{ // initial documents
		Document{UUID: uuid1, Tags: initTags},
		Document{UUID: uuid2, Tags: initTags},
		Document{UUID: uuid3, Tags: initTags},
		Document{UUID: uuid4, Tags: initTags},
		Document{UUID: uuid5, Tags: initTags},

		// change location on uuid 1, 3, 5
		Document{UUID: uuid1, Tags: map[string]string{"Location/Room": "411"}},
		Document{UUID: uuid3, Tags: map[string]string{"Location/Room": "420"}},
		Document{UUID: uuid5, Tags: map[string]string{"Location/Room": "405"}},

		// add new tags describing temperature
		Document{UUID: uuid1, Tags: temperatureTags},
		Document{UUID: uuid2, Tags: temperatureTags},
		Document{UUID: uuid3, Tags: temperatureTags},
		Document{UUID: uuid4, Tags: temperatureTags},
		Document{UUID: uuid5, Tags: temperatureTags},

		// add exposure
		Document{UUID: uuid1, Tags: map[string]string{"Metadata/Exposure": "South"}},
		Document{UUID: uuid2, Tags: map[string]string{"Metadata/Exposure": "West"}},
		Document{UUID: uuid3, Tags: map[string]string{"Metadata/Exposure": "North"}},
		Document{UUID: uuid4, Tags: map[string]string{"Metadata/Exposure": "East"}},
		Document{UUID: uuid5, Tags: map[string]string{"Metadata/Exposure": "South"}},
	} {
		if err := backend.Insert(&doc); err != nil {
			log.Fatal("Error inserting: %v", err)
		}
	}

	flag.Parse()
	os.Exit(m.Run())
}

func TestInsert(t *testing.T) {
	user := os.Getenv("ARONNAXTESTUSER")
	pass := os.Getenv("ARONNAXTESTPASS")
	dbname := os.Getenv("ARONNAXTESTDB")
	backend := newBackend(user, pass, dbname)
	uuid, _ := uuid.FromString("aa45f708-8be8-11e5-86ae-5cc5d4ded1ae")

	for _, test := range []struct {
		doc Document
		ok  bool
	}{
		{
			Document{UUID: uuid, Tags: map[string]string{"key1": "val1", "key2": "val2"}},
			true,
		},
		{
			Document{UUID: uuid, Tags: map[string]string{"key1": ""}},
			true,
		},
	} {
		if err := backend.Insert(&test.doc); test.ok != (err == nil) {
			t.Errorf("Insert test failed: Expected err? %v Err: %v", test.ok, err)
		}
	}
}
