package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// Handler represents the HTTP request handlers
type Handler struct {
    MinIOOps  *MinIOOperations
    DockerOps *DockerOperations
}

func (h *Handler) ExecuteHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        BucketName string `json:"bucket_name"`
        ScriptName string `json:"script_name"`
    }
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    resp, err := sendRequest("POST", "http://localhost:8000/execute", payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var result struct {
        Message string `json:"message"`
        Result  string `json:"result"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Message: %s\nResult: %s", result.Message, result.Result)
}

func (h *Handler) LangchainExecuteHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        InputText string `json:"input_text"`
    }
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    resp, err := sendRequest("POST", "http://localhost:8000/langchain-execute", payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var result struct {
        Message string `json:"message"`
        Result  string `json:"result"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Message: %s\nResult: %s", result.Message, result.Result)
}

func (h *Handler) HydrateDataHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        URL        string `json:"url"`
        BucketName string `json:"bucket_name"`
    }
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    resp, err := sendRequest("POST", "http://localhost:8000/hydrate-data", payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var result struct {
        Message string `json:"message"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Message: %s", result.Message)
}

func (h *Handler) MinioWebhookHandler(w http.ResponseWriter, r *http.Request) {
    var payload map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    resp, err := sendRequest("POST", "http://localhost:8000/minio-webhook", payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var result struct {
        Message string `json:"message"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Message: %s", result.Message)
}

func sendRequest(method, url string, payload interface{}) (*http.Response, error) {
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    return resp, nil
}