package outputters

import (
  "github.com/freeformz/shh/utils"
  "github.com/freeformz/shh/mm"
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
  batchLength int = utils.GetEnvWithDefaultInt("SHH_LIBRATO_BATCH_SIZE",50)
  batchTimeout time.Duration = utils.GetEnvWithDefaultDuration("SHH_LIBRATO_BATCH_TIMEOUT", "1s")
  batches chan[]*mm.Measurement = make(chan[]*mm.Measurement, 4)
)

type Librato struct {}

func init() {
  go deliver()
}

func deliver() {
  for batch := range batches {
    fmt.Println(batch)
  }
}

func (out Librato) Output (measurements <-chan *mm.Measurement) {
  ticker := time.Tick(batchTimeout)
  batch := makeBatch()
  for {
    select {
    case <-ticker:
      if len(batch) > 0 {
        batches <- batch
        batch = makeBatch()
      }
    case measurement := <-measurements:
      batch = append(batch, measurement)
      if len(batch) == cap(batch) {
        batches <- batch
        batch = makeBatch()
      }
    }
  }
}

func makeBatch() []*mm.Measurement {
  return make([]*mm.Measurement, 0, batchLength)
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



