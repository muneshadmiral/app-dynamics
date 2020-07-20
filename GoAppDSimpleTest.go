package main

import (
	appd "appdynamics"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func main() {

	fmt.Println("Arch Check")

	// Exec ldd test
	ldd_cmd := exec.Command("ldd", "--version")

	var ldd_out bytes.Buffer
	ldd_cmd.Stdout = &ldd_out

	ldd_err := ldd_cmd.Run()
	if ldd_err != nil {
		log.Fatal(ldd_err)
	}

	fmt.Printf("ldd --version:\n")
	fmt.Printf("%s\n", ldd_out.String())

	// Exec uname test
	uname_cmd := exec.Command("uname", "-a")

	var uname_out bytes.Buffer
	uname_cmd.Stdout = &uname_out

	uname_err := uname_cmd.Run()
	if uname_err != nil {
		log.Fatal(uname_err)
	}

	fmt.Printf("uname -a:\n")
	fmt.Printf("%s\n", uname_out.String())

	//Configure AppD
	cfg := appd.Config{}

	// Controller
	cfg.Controller.Host = "turing2020071222040213.saas.appdynamics.com"
	cfg.Controller.Port = 443
	cfg.Controller.UseSSL = true
	cfg.Controller.Account = "turing2020071222040213"
	cfg.Controller.AccessKey = "ovz4a1b6qhy5"

	// App Context
	cfg.AppName = "ms-login"
	cfg.TierName = ""
	cfg.NodeName = ""

	// misc
	cfg.InitTimeoutMs = 1000

	// init the SDK
	if err := appd.InitSDK(&cfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}

	// Run some BTs
	maxBtCount := 200000
	btCount := 0

	fmt.Print("Doing something")
	for btCount < maxBtCount {
		// start the "Checkout" transaction
		btHandle := appd.StartBT("MyTestGolangBT", "")

		// do something....
		fmt.Print("A")
		milliseconds := 1000
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		setSnapshotAttributes(btHandle, "1", "1")
		post()
		// end the transaction
		appd.EndBT(btHandle)
		fmt.Print("Call")
		btCount++
		fmt.Println(btCount)
	}
	fmt.Print("\n")

	// Stop/Clean up the AppD SDK.
	appd.TerminateSDK()

	fmt.Println("Terminated")

}

func setSnapshotAttributes(bt appd.BtHandle, key, value string) {
	if appd.IsBTSnapshotting(bt) {
		fmt.Println("snapshot")
		appd.AddUserDataToBT(bt, key, value)
		appd.SetBTURL(bt, "user/login")
	}
}
func post() {
	fmt.Println("2. Performing Http Post...")
	todo := Todo{1, 2, "lorem ipsum dolor sit amet", true}
	jsonReq, err := json.Marshal(todo)
	resp, err := http.Post("https://jsonplaceholder.typicode.com/todos", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

}

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
