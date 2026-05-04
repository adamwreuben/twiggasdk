package main

import (
	"context"
	"fmt"
	"log"

	"github.com/adamwreuben/twiggasdk/twigga"
)

func main() {
	ctx := context.Background()

	fmt.Println("Initializing SDK...")
	opt := twigga.WithCredentialsFile("./bongo.json")
	env := twigga.WithEnvironment(twigga.TEST)

	app, err := twigga.InitializeApp(ctx, opt, env)
	if err != nil {
		log.Fatalf("Error initializing Twigga App: %v\n", err)
	}
	fmt.Println("SDK Initialized Successfully!")

	clientId := "my_mobile_app_123"

	// 1. Sign up a new user
	tokens, err := app.Auth().Signup(ctx, twigga.CreateUserReq{
		Email:     "xx@xx.xx.xx",
		Password:  "SuperSecret123!",
		FirstName: "xxx",
		LastName:  "Mtata",
		ClientId:  clientId,
	})

	// Refresh a Token quietly in the background
	newTokens, _ := app.Auth().RefreshToken(ctx, tokens.RefreshToken)

	// User forgot their password
	app.Auth().ForgotPassword(ctx, "email")

	// User clicks email link and sets new password
	app.Auth().ResetPassword(ctx, "token_from_email", "newPassword")

	// User logs out
	app.Auth().Logout(ctx, newTokens.RefreshToken)

	// DATABASE
	fmt.Println("\n--- Testing Database ---")

	// Add a new document (Auto-ID)
	app.Database().Collection("users").Add(ctx, map[string]any{
		"name": "William Mtata",
		"role": "CTO",
	})

	// Set a document with a specific ID
	app.Database().Collection("users").Doc("dev_123").Set(ctx, map[string]any{
		"name":   "Steve",
		"status": "active",
	})

	// Get a document
	userDoc, err := app.Database().Collection("users").Doc("dev_123").Get(ctx)
	if err == nil {
		fmt.Printf("User Retrieved: %s\n", string(userDoc))
	}

	// ZANZIBAR AUTH (Permissions Engine)
	fmt.Println("\n--- Testing Zanzibar Auth ---")

	// Write a new permission tuple
	app.Auth().Zanzibar().Write(ctx, twigga.AuthorizationTuple{
		ObjectType:  "document",
		ObjectID:    "doc_999",
		Relation:    "owner",
		SubjectType: "user",
		SubjectID:   "dev_123",
	})

	// Check a permission in real-time
	allowed, err := app.Auth().Zanzibar().Check(ctx, "dev_123", "user", "owner", "document", "doc_999")
	if err == nil {
		fmt.Printf("Is 'dev_123' the owner of 'doc_999'? %v\n", allowed)
	}

	// SERVERLESS FUNCTIONS
	fmt.Println("\n--- Testing Functions ---")

	payload := map[string]any{"to": "user@xx.xx.xx", "subject": "Welcome"}
	result, err := app.Functions().Invoke(ctx, "sendWelcomeEmail", payload)
	if err == nil {
		fmt.Printf("Function Invoked! Result: %s\n", string(result))
	}

	// SECURITY RULES & BACKUPS
	fmt.Println("\n--- Testing Rules & DR ---")

	// Trigger an incremental disaster recovery backup
	app.Backups().Create(ctx, "incremental")
	fmt.Println("Incremental Backup Triggered!")

	// STORAGE
	fmt.Println("\n--- Testing Storage ---")

	bucketName := "profile-pictures"

	// Create a Bucket
	app.Storage().CreateBucket(ctx, bucketName)

	// Make it public
	app.Storage().Bucket(bucketName).SetPolicy(ctx, "public")

	// Upload a file (Make sure you have a dummy test.txt file in the directory!)
	// app.Storage().Bucket(bucketName).Upload(ctx, "./test.txt")

	// Get a specific file's URL
	fileData, err := app.Storage().Bucket(bucketName).File("test.txt").Get(ctx)
	if err == nil {
		fmt.Printf("File Download URL: %v\n", fileData["url"])
	}

	// Check Project Quotas/Usage
	usage, err := app.Storage().Usage(ctx)
	if err == nil {
		fmt.Printf("Total Storage Used: %v\n", usage["totalStorageUsed"])
	}

}
