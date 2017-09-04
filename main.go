package main

import (
    "encoding/json"
    "log"
    "fmt"
    "net/http"
    "github.com/codeskyblue/go-sh"
    "github.com/gorilla/mux"
    "os"
    "io/ioutil"
)

type Diff struct {
    Subject  string   `json:"subject,omitempty"`
    Origindiff   *DiffTarget `json:"origin,omitempty"`
    Targetdiff   *DiffTarget `json:"target,omitempty"`
}

type DiffTarget struct {
    Body  string `json:"body,omitempty"`
    Version string `json:"version,omitempty"`
}

type DiffResponse struct {
    Body string `json:"body,omitempty"`
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
    var diff Diff
    var response DiffResponse
    _ = json.NewDecoder(req.Body).Decode(&diff)
    response.Body, _ = getDiff(diff.Origindiff.Body, diff.Targetdiff.Body)
    json.NewEncoder(w).Encode(response)
}

func getDiff(origin, target string) (string, error) {

    originTempfile, err := ioutil.TempFile("", "example")
    if err != nil {
        log.Fatal("wef",err)
    }

    targetTempfile, err := ioutil.TempFile("", "example")
    if err != nil {
        log.Fatal("1",err)
    }

    defer os.Remove(originTempfile.Name()) // clean up
    defer os.Remove(targetTempfile.Name()) // clean up

    if _, err := originTempfile.Write([]byte(origin)); err != nil {
        log.Fatal("3",err)
    }

    if _, err := targetTempfile.Write([]byte(target)); err != nil {
        log.Fatal("4",err)
    }
    output, _ := sh.Command("diff", "-u", originTempfile.Name(), targetTempfile.Name()).Output()
    fmt.Println(string(output))
    return string(output), nil
}


func main() {
    router := mux.NewRouter()
    router.HandleFunc("/diff", HandleRequest).Methods("POST")
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
    log.Fatal(http.ListenAndServe(":8085", router))
}

