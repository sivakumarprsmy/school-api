package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/periasamy/school-api/internal/mcptools"
)

func main() {
	baseURL := os.Getenv("SCHOOL_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	client := mcptools.NewSchoolClient(baseURL)

	s := server.NewMCPServer("school-api", "1.0.0")

	// --- list_students ---
	s.AddTool(
		mcp.NewTool("list_students",
			mcp.WithDescription("List all students"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			students, err := client.ListStudents(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list students: %v", err)), nil
			}
			data, _ := json.MarshalIndent(students, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// --- get_student ---
	s.AddTool(
		mcp.NewTool("get_student",
			mcp.WithDescription("Get a student by their ID"),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description("The student's numeric ID"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := int(req.Params.Arguments["id"].(float64))
			student, err := client.GetStudent(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get student: %v", err)), nil
			}
			data, _ := json.MarshalIndent(student, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// --- add_student ---
	s.AddTool(
		mcp.NewTool("add_student",
			mcp.WithDescription("Add a new student to the school"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Student's full name"),
			),
			mcp.WithString("email",
				mcp.Required(),
				mcp.Description("Student's email address"),
			),
			mcp.WithNumber("age",
				mcp.Required(),
				mcp.Description("Student's age (1–120)"),
			),
			mcp.WithString("grade",
				mcp.Required(),
				mcp.Description("Student's grade, e.g. '10th'"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, _ := req.Params.Arguments["name"].(string)
			email, _ := req.Params.Arguments["email"].(string)
			age := int(req.Params.Arguments["age"].(float64))
			grade, _ := req.Params.Arguments["grade"].(string)

			student, err := client.AddStudent(ctx, mcptools.AddStudentRequest{
				Name: name, Email: email, Age: age, Grade: grade,
			})
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to add student: %v", err)), nil
			}
			data, _ := json.MarshalIndent(student, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// --- delete_student ---
	s.AddTool(
		mcp.NewTool("delete_student",
			mcp.WithDescription("Delete a student by their ID"),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description("The student's numeric ID to delete"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := int(req.Params.Arguments["id"].(float64))
			if err := client.DeleteStudent(ctx, id); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to delete student: %v", err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Student %d deleted successfully", id)), nil
		},
	)

	log.Printf("school-api MCP server starting (API: %s)\n", baseURL)
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
