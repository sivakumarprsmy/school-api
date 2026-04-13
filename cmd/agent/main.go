package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/periasamy/school-api/internal/mcptools"
)

// ── Tool input shapes ──────────────────────────────────────────────────────

type validateEmailInput struct {
	Email string `json:"email"`
}

type lookupStudentInput struct {
	Email string `json:"email"`
}

type enrollStudentInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
}

// ── Tool definitions ───────────────────────────────────────────────────────

func buildTools() []anthropic.ToolUnionParam {
	makeSchema := func(properties map[string]interface{}, required []string) anthropic.ToolInputSchemaParam {
		return anthropic.ToolInputSchemaParam{
			Properties: properties,
			Required:   required,
		}
	}

	validateEmail := anthropic.ToolParam{
		Name:        "validate_email",
		Description: anthropic.String("Check whether an email address has a valid format and belongs to an allowed domain. Returns {\"valid\": true} or {\"valid\": false, \"error\": \"reason\"}."),
		InputSchema: makeSchema(map[string]interface{}{
			"email": map[string]interface{}{
				"type":        "string",
				"description": "The email address to validate",
			},
		}, []string{"email"}),
	}

	lookupStudent := anthropic.ToolParam{
		Name:        "lookup_student_by_email",
		Description: anthropic.String("Check whether a student with the given email already exists in the system. Returns {\"found\": true, \"student\": {...}} or {\"found\": false}."),
		InputSchema: makeSchema(map[string]interface{}{
			"email": map[string]interface{}{
				"type":        "string",
				"description": "The email address to look up",
			},
		}, []string{"email"}),
	}

	enrollStudent := anthropic.ToolParam{
		Name:        "enroll_student",
		Description: anthropic.String("Enroll a new student in the school system. All four fields are required. Only call this after validating the email and confirming there is no duplicate."),
		InputSchema: makeSchema(map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Student's full name",
			},
			"email": map[string]interface{}{
				"type":        "string",
				"description": "Student's email address",
			},
			"age": map[string]interface{}{
				"type":        "integer",
				"description": "Student's age (1–120)",
			},
			"grade": map[string]interface{}{
				"type":        "string",
				"description": "Student's grade, e.g. '10th'",
			},
		}, []string{"name", "email", "age", "grade"}),
	}

	return []anthropic.ToolUnionParam{
		{OfTool: &validateEmail},
		{OfTool: &lookupStudent},
		{OfTool: &enrollStudent},
	}
}

// ── Tool executor ──────────────────────────────────────────────────────────

func executeTool(
	name string,
	rawInput json.RawMessage,
	validator *mcptools.EmailValidator,
	schoolClient *mcptools.SchoolClient,
) (string, bool) {
	switch name {
	case "validate_email":
		var in validateEmailInput
		if err := json.Unmarshal(rawInput, &in); err != nil {
			return fmt.Sprintf(`{"valid":false,"error":"bad input: %v"}`, err), true
		}
		if err := validator.Validate(in.Email); err != nil {
			out, _ := json.Marshal(map[string]interface{}{"valid": false, "error": err.Error()})
			return string(out), false
		}
		return `{"valid":true}`, false

	case "lookup_student_by_email":
		var in lookupStudentInput
		if err := json.Unmarshal(rawInput, &in); err != nil {
			return fmt.Sprintf(`{"found":false,"error":"bad input: %v"}`, err), true
		}
		students, err := schoolClient.ListStudents(context.Background())
		if err != nil {
			return fmt.Sprintf(`{"found":false,"error":"failed to list students: %v"}`, err), true
		}
		needle := strings.ToLower(strings.TrimSpace(in.Email))
		for _, s := range students {
			if strings.ToLower(s.Email) == needle {
				data, _ := json.Marshal(map[string]interface{}{"found": true, "student": s})
				return string(data), false
			}
		}
		return `{"found":false}`, false

	case "enroll_student":
		var in enrollStudentInput
		if err := json.Unmarshal(rawInput, &in); err != nil {
			return fmt.Sprintf(`{"enrolled":false,"error":"bad input: %v"}`, err), true
		}
		student, err := schoolClient.AddStudent(context.Background(), mcptools.AddStudentRequest{
			Name:  in.Name,
			Email: in.Email,
			Age:   in.Age,
			Grade: in.Grade,
		})
		if err != nil {
			out, _ := json.Marshal(map[string]interface{}{"enrolled": false, "error": err.Error()})
			return string(out), false
		}
		data, _ := json.Marshal(map[string]interface{}{"enrolled": true, "student": student})
		return string(data), false

	default:
		return fmt.Sprintf(`{"error":"unknown tool %q"}`, name), true
	}
}

// ── Main ───────────────────────────────────────────────────────────────────

func main() {
	// Read request from stdin
	stat, _ := os.Stdin.Stat()
	var request string
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// piped input
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read stdin: %v", err)
		}
		request = strings.TrimSpace(string(data))
	} else {
		// interactive — prompt the user
		fmt.Print("Enrollment request: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			request = strings.TrimSpace(scanner.Text())
		}
	}

	if request == "" {
		fmt.Fprintln(os.Stderr, "error: no enrollment request provided")
		os.Exit(1)
	}

	// Set up dependencies
	baseURL := os.Getenv("SCHOOL_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	schoolClient := mcptools.NewSchoolClient(baseURL)
	validator := mcptools.NewEmailValidator()

	// Set up Anthropic client
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	systemPrompt := `You are an enrollment agent for a school system. Your job is to enroll students based on the user's request.

Follow these steps in order:
1. Extract the student's name, email, age, and grade from the request.
2. If any required field (name, email, age, grade) is missing, ask the user to provide it.
3. Call validate_email to check the email address.
   - If invalid format: report the format error and ask for a corrected email.
   - If wrong domain: tell the user which domains are accepted and ask for a corrected email.
4. Call lookup_student_by_email to check for duplicates.
   - If a student already exists with that email: inform the user and do NOT enroll again.
5. If all checks pass, call enroll_student to complete the enrollment.
6. Confirm the result clearly to the user.

Be concise. Do not invent data — only use what the user provides.`

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(request)),
	}

	ctx := context.Background()
	tools := buildTools()

	// Agentic loop
	for {
		resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     anthropic.ModelClaudeOpus4_6,
			MaxTokens: 4096,
			System: []anthropic.TextBlockParam{
				{Text: systemPrompt},
			},
			Tools:    tools,
			Messages: messages,
		})
		if err != nil {
			log.Fatalf("API error: %v", err)
		}

		// Append assistant turn
		messages = append(messages, resp.ToParam())

		// Process response content
		var toolResults []anthropic.ContentBlockParamUnion
		for _, block := range resp.Content {
			switch variant := block.AsAny().(type) {
			case anthropic.TextBlock:
				if variant.Text != "" {
					fmt.Println(variant.Text)
				}
			case anthropic.ToolUseBlock:
				rawInput := json.RawMessage(variant.JSON.Input.Raw())
				result, isError := executeTool(variant.Name, rawInput, validator, schoolClient)
				toolResults = append(toolResults, anthropic.NewToolResultBlock(variant.ID, result, isError))
			}
		}

		// If there were tool calls, feed results back and continue
		if len(toolResults) > 0 {
			messages = append(messages, anthropic.NewUserMessage(toolResults...))
			continue
		}

		// No more tool calls — we're done
		break
	}
}
