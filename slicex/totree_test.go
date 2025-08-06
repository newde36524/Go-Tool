package slicex

import (
	"encoding/json"
	"fmt"
	"testing"
)

type NodeInfo struct {
	Data     string
	Id       string
	ParentId string
}

func TestToTree(t *testing.T) {
	data := make([]NodeInfo, 0)
	data = append(data, NodeInfo{"a", "1", ""})
	data = append(data, NodeInfo{"b", "2", "1"})
	data = append(data, NodeInfo{"c", "3", "2"})
	data = append(data, NodeInfo{"d", "4", "2"})

	data = append(data, NodeInfo{"a2", "11", ""})
	data = append(data, NodeInfo{"b2", "21", "11"})
	data = append(data, NodeInfo{"c2", "31", "21"})
	data = append(data, NodeInfo{"d2", "41", "21"})

	root, _ := ToTree(data, func(e NodeInfo) (string, string) {
		return e.Id, e.ParentId
	})
	for _, node := range root {
		bs, _ := json.Marshal(node)
		fmt.Println(string(bs))
	}
	/*
		{
		    "Value": {
		        "Data": "a",
		        "Id": "1",
		        "ParentId": ""
		    },
		    "Children": {
		        "2": {
		            "Value": {
		                "Data": "b",
		                "Id": "2",
		                "ParentId": "1"
		            },
		            "Children": {
		                "3": {
		                    "Value": {
		                        "Data": "c",
		                        "Id": "3",
		                        "ParentId": "2"
		                    },
		                    "Children": {}
		                },
		                "4": {
		                    "Value": {
		                        "Data": "d",
		                        "Id": "4",
		                        "ParentId": "2"
		                    },
		                    "Children": {}
		                }
		            }
		        }
		    }
		}
	*/
}
