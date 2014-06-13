package splunk

import (
        "fmt"
        "net/url"
	"encoding/json"
	"bytes"
	"strings"
	"bufio"
)

type Value interface{}

type Row struct {
	Preview bool `json:"preview"`
	Offset int `json:"offset"`
	Result map[string]Value `json:"result"`
}

type Rows []Row

func NewRow() (ret Row) {
	ret = Row{}
	ret.Result = make(map[string]Value)
	return
}

func parseLine(line string) (r Row,err error) {
	fbuf := bytes.NewBufferString(line)
	dec := json.NewDecoder(fbuf)
	r = NewRow()
	return r, dec.Decode(&r)
}

func (conn SplunkConnection) Search(searchString string) (rows Rows,err error) {
	data := make(url.Values)
	data.Add("search",searchString)
	data.Add("output_mode","json")

	/* TODO: return stream in order to read reponses that do not fit in memory. */
	response, err := conn.httpPost(fmt.Sprintf("%s/servicesNS/admin/search/search/jobs/export",conn.BaseURL),&data)

	if err != nil {
		return nil,err
	}

	lines := strings.Split(response,"\n")
	rows = make(Rows,len(lines))
	var ni int
	ni = 0

	for _,v := range lines {
		if r,err := parseLine(v); err != nil {
			fmt.Printf("Could not decode line: '%s' %s\n",v,err)
		} else {
			rows[ni] = r
			ni++
		}
	}
	return rows[:ni], nil;
}


func (conn SplunkConnection) SearchStream(searchString string, params ...map[string]string) (events chan *Row,err error) {
	data := make(url.Values)
	data.Add("search",searchString)
	data.Add("output_mode","json")

	for _,m := range params {
		for k,v := range m {
			data.Add(k,v)
		}
	}

	response, err := conn.httpCall(fmt.Sprintf("%s/servicesNS/admin/search/search/jobs/export",conn.BaseURL),"POST",&data)

	if err != nil {
		return nil,err
	}

	events = make(chan *Row,50)

	go func() { // Using closures here (events,response)
		scanner := bufio.NewScanner(response.Body)
		defer response.Body.Close()

		for scanner.Scan() {
			line := scanner.Text()
			if row,err := parseLine(line); err != nil {
				fmt.Printf("Could not decode line: '%s' %s\n",line,err)
			} else {
				events <- &row
			}
		}
		events <- nil //Signal EOF
	}()

	return events,err
}