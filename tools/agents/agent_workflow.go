package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

// ---------------------------------------------------------
// AGENT SYSTEM PROMPTS
// ---------------------------------------------------------

const backendPrompt = `You are an Expert Go Backend Developer.
You strictly follow Clean Architecture.
Your rules:
1. Always separate HTTP handlers, Usecases, Domain entities, and Repositories.
2. Use sqlc for database queries and gin-gonic for routing.
3. Provide ONLY the raw Go code. Do not use markdown blocks like ` + "```go" + `. Just the code.`

const logicCheckPrompt = `You are a Strict QA Engineer and Security Auditor. 
Analyze the provided Go code.
1. Look for edge cases, missing validations, and double-booking vulnerabilities.
2. Ensure it strictly matches the business rules.
3. Provide ONLY the raw Go code with your fixes applied. Do not use markdown blocks like ` + "```go" + `. Just the code.`

const reviewerPrompt = `You are a Strict Code Reviewer.
Review the provided Go code.
1. Check for SOLID principles, race conditions, and unhandled errors.
2. Provide the refactored, clean version of the code.
3. Provide ONLY the raw Go code and your comments as standard Go comments (//). Do not use markdown blocks like ` + "```go" + `. Just the code.`

// Helper to extract text from the new SDK response format
func extractText(resp *genai.GenerateContentResponse) string {
	text := ""
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			text += part.Text
		}
	}
	return text
}

func main() {
	// 1. Ensure API key is set
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("Please set the GEMINI_API_KEY environment variable")
	}

	// 2. Read the input file (e.g., a Use Case Markdown file)
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run agent_workflow.go <input_file.md>")
	}
	inputFile := os.Args[1]
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	ctx := context.Background()

	// The new SDK automatically picks up the GEMINI_API_KEY environment variable
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// ---------------------------------------------------------
	// AGENT 1: The Backend Developer (Drafting)
	// ---------------------------------------------------------
	fmt.Println("🤖 Backend Agent is writing the initial code...")

	// New SDK syntax for passing System Instructions
	backendConfig := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: backendPrompt}},
		},
		Temperature: genai.Ptr[float32](0.2), // Low temperature for code generation
	}

	// New SDK syntax for generating content
	backendResp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(string(inputData)),
		backendConfig,
	)
	if err != nil {
		log.Fatalf("Backend Agent failed: %v", err)
	}

	draftCode := extractText(backendResp)

	os.WriteFile("draft_output.go", []byte(draftCode), 0o644)
	fmt.Println("   ✅ Draft code saved to draft_output.go\n")

	// ---------------------------------------------------------
	// AGENT 2: The Logic Check / QA (Bug Hunting)
	// ---------------------------------------------------------
	fmt.Println("🕵️ Logic Check Agent is hunting for bugs and edge cases...")

	logicConfig := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: logicCheckPrompt}},
		},
		Temperature: genai.Ptr[float32](0.2),
	}

	logicResp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-pro",
		genai.Text(draftCode),
		logicConfig,
	)
	if err != nil {
		log.Fatalf("Logic Check Agent failed: %v", err)
	}

	qaCode := extractText(logicResp)

	os.WriteFile("qa_output.go", []byte(qaCode), 0o644)
	fmt.Println("   ✅ QA checked code saved to qa_output.go\n")

	// ---------------------------------------------------------
	// AGENT 3: The Code Reviewer (Final Polish)
	// ---------------------------------------------------------
	fmt.Println("🧐 Reviewer Agent is polishing and refactoring...")

	reviewerConfig := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: reviewerPrompt}},
		},
		Temperature: genai.Ptr[float32](0.2),
	}

	reviewerResp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-pro",
		genai.Text(qaCode),
		reviewerConfig,
	)
	if err != nil {
		log.Fatalf("Reviewer Agent failed: %v", err)
	}

	finalCode := extractText(reviewerResp)

	os.WriteFile("final_reviewed_output.go", []byte(finalCode), 0o644)
	fmt.Println("   ✅ Final reviewed code saved to final_reviewed_output.go\n")

	fmt.Println("🎉 3-Agent Workflow Complete! Check final_reviewed_output.go")
}
