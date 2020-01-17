package main

import (
    "fmt"
    "time"
    "os"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

func main() {
    dateFmt := "2006-01-02"
    url := "https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson"
    fullUrl := url +
        "&starttime=" + time.Now().Add(-24 * time.Hour).Format(dateFmt) +
        "&endtime=" + time.Now().Format(dateFmt)

    resp, e := http.Get(fullUrl)
    if e != nil {
        fmt.Println("Error requesting data: " + e.Error())
        os.Exit(1)
    }
    body, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    
    var data interface{}
    json.Unmarshal(body, &data)

    for _, feature := range arr(obj(data)["features"]) {
        props := obj(obj(feature)["properties"])
        fmt.Println(props["place"], props["mag"])
    }
}

func obj(val interface{}) (map[string]interface{}) {
    return val.(map[string]interface{})
}

func arr(val interface{}) ([]interface{}) {
    return val.([]interface{})
}
