package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"

// 	"gopkg.in/yaml.v2"
// 	//"github.com/gofiber/fiber/v2"
// 	//"github.com/gofiber/fiber/v2/middleware/proxy"
// )

// var Config config
// var Addresses []string
// var queue map[string][]string
// var activeAddress map[string]bool
// var targetUrl string
// var requests []http.Request
// var toSend map[string]chan *http.Request

// func main() {
// 	err := Config.GetConfig("settings.yml")
// 	if err != nil {
// 		log.Println(err)
// 		os.Exit(1)
// 	}

// 	Addresses = Config.Servers.Url
// 	queue = make(map[string][]string)
// 	targetUrl = "192.168.111.111:1010"

// 	for _, addr := range Addresses {
// 		activeAddress[addr] = true
// 		toSend[addr] = make(chan *http.Request)
// 		go SendProxy(addr, toSend[addr])
// 	}

// 	http.HandleFunc("/send", Proxy)
// 	err = http.ListenAndServe(":3333", nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	// app := fiber.New(fiber.Config{
// 	// 	BodyLimit: 1 * 1024 * 1024,
// 	// })
// 	// app.Use(proxy.Balancer(proxy.Config{
// 	// 	Servers: Addresses,
// 	// }))
// 	// app.Get("/", Index)
// 	// app.Get("/something", proxy.Forward(targetUrl))
// 	//app.Get("/Req", Proxy)
// 	//app.Server().Name = "Proxy balancer"

// 	log.Println()
// }

// func Index(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "OK")
// }

// func Proxy(w http.ResponseWriter, r *http.Request) {
// 	requests = append(requests, *r)
// 	//SendProxy(b, r.Header)
// 	io.WriteString(w, "ok")
// }

// func GetMin() string {
// 	min := -1
// 	var result string
// 	for _, addr := range Addresses {
// 		if activeAddress[addr] && min == -1 {
// 			min = len(queue[addr])
// 		}
// 		if min > len(queue[addr]) {
// 			min = len(queue[addr])
// 			result = addr
// 		}
// 	}
// 	return result
// }

// func Balance(w http.ResponseWriter, r *http.Request) {

// }

// func IsAlive(url, method string) bool {
// 	r, err := http.Get(url)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return r.StatusCode == 200
// }

// func SendProxy(url string, in chan *http.Request) error {
// 	client := http.Client{}
// 	var reqIn *http.Request
// 	for {
// 		reqIn = <-in
// 		req, err := http.NewRequest(reqIn.Method, url, nil)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 		req.Header = reqIn.Header
// 		req.Body = reqIn.Body

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 		if resp.StatusCode != 200 {
// 			activeAddress[url] = false
// 			return err
// 		} else {
// 			activeAddress[url] = true
// 		}
// 	}
// }

// func (c *config) GetConfig(configPath string) error {
// 	file, err := os.Open(configPath)
// 	if err != nil {
// 		return err
// 	}

// 	defer func(file *os.File) {
// 		err := file.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}(file)

// 	d := yaml.NewDecoder(file)

// 	if err := d.Decode(&c); err != nil {
// 		return err
// 	}
// 	return nil
// }

// type config struct {
// 	Servers struct {
// 		Url []string `yaml:"url"`
// 	} `yaml:"servers"`
// }

// type requestC struct {
// }

// //
// //

// // func Index(c *fiber.Ctx) error {
// // 	return c.Status(200).SendString("Got it " + c.IP())
// // }

// // func Proxy(c *fiber.Ctx) error {
// // 	//c.Response().Header = c.Request().Header
// // 	//http.NewRequest("GET", url, nil)
// // 	return c.Status(200).SendString("Ok")
// // }

// // type Queue interface {
// // 	Front() *list.Element
// // 	Len() int
// // 	Add(interface{})
// // 	Remove()
// // }

// // type queueImpl struct {
// // 	*list.List
// // }

// // func (q *queueImpl) Add(v interface{}) {
// // 	q.PushBack(v)
// // }

// // func (q *queueImpl) Remove() {
// // 	e := q.Front()
// // 	q.List.Remove(e)
// // }

// // // New is a new instance of a Queue
// // func New() Queue {
// // 	return &queueImpl{list.New()}
// // }
