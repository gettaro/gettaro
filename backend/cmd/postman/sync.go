package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

const (
	postmanAPIBaseURL = "https://api.getpostman.com"
)

type PostmanCollection struct {
	Collection struct {
		Info struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Schema      string `json:"schema"`
		} `json:"info"`
		Item []PostmanItem `json:"item"`
	} `json:"collection"`
}

type PostmanItem struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Request     PostmanRequest `json:"request,omitempty"`
	Item        []PostmanItem  `json:"item,omitempty"`
}

type PostmanRequest struct {
	Method      string          `json:"method"`
	Header      []PostmanHeader `json:"header"`
	URL         PostmanURL      `json:"url"`
	Description string          `json:"description"`
	Body        *PostmanBody    `json:"body,omitempty"`
}

type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanURL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Path     []string `json:"path"`
}

type PostmanBody struct {
	Mode    string                 `json:"mode"`
	Raw     string                 `json:"raw"`
	Options map[string]interface{} `json:"options,omitempty"`
}

func main() {
	// Get Postman API key from environment
	apiKey := os.Getenv("POSTMAN_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: POSTMAN_API_KEY environment variable is required")
		os.Exit(1)
	}

	// Get collection ID from environment
	collectionID := os.Getenv("POSTMAN_COLLECTION_ID")
	if collectionID == "" {
		fmt.Println("Error: POSTMAN_COLLECTION_ID environment variable is required")
		os.Exit(1)
	}

	// Read OpenAPI spec
	specPath := filepath.Join("..", "..", "http", "openapi.yaml")
	specData, err := ioutil.ReadFile(specPath)
	if err != nil {
		fmt.Printf("Error reading OpenAPI spec: %v\n", err)
		os.Exit(1)
	}

	// Parse OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(specData)
	if err != nil {
		fmt.Printf("Error parsing OpenAPI spec: %v\n", err)
		os.Exit(1)
	}

	// Convert OpenAPI spec to Postman collection
	collection := convertToPostmanCollection(doc)

	// Update Postman collection
	err = updatePostmanCollection(apiKey, collectionID, collection)
	if err != nil {
		fmt.Printf("Error updating Postman collection: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully updated Postman collection")
}

func convertToPostmanCollection(doc *openapi3.T) PostmanCollection {
	var collection PostmanCollection
	collection.Collection.Info.Name = doc.Info.Title
	collection.Collection.Info.Description = doc.Info.Description
	collection.Collection.Info.Schema = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

	// Group paths by their first segment
	folders := make(map[string][]PostmanItem)
	for path, pathItem := range doc.Paths.Map() {
		segments := splitPath(path)
		if len(segments) == 0 {
			continue
		}
		folderName := segments[0]

		for method, operation := range pathItem.Operations() {
			item := PostmanItem{
				Name:        operation.Summary,
				Description: operation.Description,
				Request: PostmanRequest{
					Method: method,
					URL: PostmanURL{
						Raw:      "{{baseUrl}}" + path,
						Protocol: "http",
						Host:     []string{"{{baseUrl}}"},
						Path:     segments,
					},
				},
			}

			// Add headers
			if operation.Security != nil {
				item.Request.Header = append(item.Request.Header, PostmanHeader{
					Key:   "Authorization",
					Value: "Bearer {{token}}",
				})
			}

			// Add request body if present
			if operation.RequestBody != nil {
				content := operation.RequestBody.Value.Content.Get("application/json")
				if content != nil && content.Schema != nil {
					example := generateExampleFromSchema(content.Schema.Value)
					item.Request.Body = &PostmanBody{
						Mode: "raw",
						Raw:  example,
						Options: map[string]interface{}{
							"raw": map[string]string{
								"language": "json",
							},
						},
					}
				}
			}

			folders[folderName] = append(folders[folderName], item)
		}
	}

	// Convert folders to Postman collection items
	for name, items := range folders {
		collection.Collection.Item = append(collection.Collection.Item, PostmanItem{
			Name: name,
			Item: items,
		})
	}

	return collection
}

func generateExampleFromSchema(schema *openapi3.Schema) string {
	if schema == nil {
		return "{}"
	}

	example := make(map[string]interface{})
	for name, prop := range schema.Properties {
		if prop.Value.Type == "string" {
			example[name] = "string"
		} else if prop.Value.Type == "integer" {
			example[name] = 0
		} else if prop.Value.Type == "boolean" {
			example[name] = true
		} else if prop.Value.Type == "array" {
			if prop.Value.Items != nil {
				example[name] = []interface{}{generateExampleFromSchema(prop.Value.Items.Value)}
			} else {
				example[name] = []interface{}{}
			}
		} else if prop.Value.Type == "object" {
			example[name] = generateExampleFromSchema(prop.Value)
		}
	}

	jsonBytes, _ := json.MarshalIndent(example, "", "  ")
	return string(jsonBytes)
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func updatePostmanCollection(apiKey, collectionID string, collection PostmanCollection) error {
	// Convert collection to JSON
	collectionJSON, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("error marshaling collection: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/collections/%s", postmanAPIBaseURL, collectionID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(collectionJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error updating collection: %s", string(body))
	}

	return nil
}
