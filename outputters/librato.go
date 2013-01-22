package outputters

import (
  "io/ioutil"
  "encoding/json"
  "net/http"
  "time"
  "math/rand"
  "log"
  "fmt"
  "bytes"
  "os"
)

type Metric struct {
  Name string `json:"name"`
  Value string `json:"value"`
  When int64  `json:"measure_time"`
  Source string `json:"source,omitempty"`
}

type PostBody struct {
  Gauges []Metric `json:"gauges,omitempty"`
  Counters []Metric `json:"counters,omitempty"`
}

const (
  LIBRATO_URL = "https://metrics-api.librato.com/v1/metrics"
)

var (
  user string = os.Getenv("SHH_LIBRATO_USER")
  token string = os.Getenv("SHH_LIBRATO_TOKEN")
  batchLength string = os.Getenv("SHH_LIBRATO_BATCH")
)

func init {
}

type Librato {}

func (out Librato) Output (measurements <-chan *mm.Measurement) {
  ticker := time.Tick(time.Duration(batchTimeout) * time.Millisecond)
  batch := makeBatch() //make([]LogMessage, 0, batchSize)
  for {
    select {
    case <-ticker:
      if len(batch) > 0 {
        batches <- batch
        batch = makeBatch() //make([]LogMessage, 0, batchSize)
      }
    case line := <-lines:
      batch = append(batch, parseLogMessage(line))
      if len(batch) == cap(batch) {
        batches <- batch
        batch = makeBatch() //make([]LogMessage, 0, batchSize)
      }
    }
  }
}

func foo() {
  rand.Seed(time.Now().Unix())
  ticker := time.Tick(10 * time.Second)
  for t := range ticker {
    g1 := Metric{"gauge.test.no_source", fmt.Sprintf("%f", rand.Float64() * 100), t.Unix(), ""}
    g2 := Metric{"gauge.test.with_source", fmt.Sprintf("%f", rand.Float64() * 100), t.Unix(),"a"}
    g3 := Metric{"gauge.test.with_source", fmt.Sprintf("%f", rand.Float64() * 100), t.Unix(),"b"}
    /*c1 := Metric{"counter.test.no_source", fmt.Sprintf("%d", rand.Int()), t.Unix(), ""}
    c2 := Metric{"counter.test.with_source", fmt.Sprintf("%d", rand.Int()), t.Unix(),"a"}
    c3 := Metric{"counter.test.with_source", fmt.Sprintf("%d", rand.Int()), t.Unix(),"b"}*/

    payload := new(PostBody)
    payload.Gauges = []Metric{g1,g2,g3}

    //payload := Post{[]Metric{c1,c2,c3}}
    fmt.Println("Payload")
    fmt.Println(payload)
    j, err := json.Marshal(payload)
    fmt.Println("JSON")
    fmt.Printf("%s\n",j)
    if err != nil { log.Fatal(err) }
    body := bytes.NewBuffer(j)
    req, err := http.NewRequest("POST", LIBRATO_URL, body)
    if err != nil { log.Fatal(err) }
    req.Header.Add("Content-Type", "application/json")
    req.SetBasicAuth(user,token)
    fmt.Println("Request")
    fmt.Println(req)
    resp, err := http.DefaultClient.Do(req)
    if err != nil { log.Fatal(err) }
    fmt.Println("Response")
    fmt.Println(resp)
    if resp.StatusCode/100 != 2 {
      fmt.Println("Response Body")
      b, _ := ioutil.ReadAll(resp.Body)
      fmt.Printf("%s\n",b)
    }
    resp.Body.Close()
  }
}



