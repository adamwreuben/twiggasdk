# Twigga SDK

This SDK provides a seamless interface to interact with Twigga authentication, NoSQL database, serverless functions, cloud storage, and granular permissions engine (Zanzibar).



## Installation

```bash
go get github.com/adamwreuben/twiggasdk
```

## Prerequisites
Ensure you have your service account credentials saved as bongo.json in your project root or set securely in your environment variables.
```
{
  "projectId": "your-project-id",
  "appId": "your-client-id",
  "appSecret": "your-client-secret"
}
```

## Initialization
To start using the SDK, initialize the App client using your credentials. The SDK supports environments (e.g., TEST or PROD) and handles token fetching automatically.

```
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/adamwreuben/twiggasdk"
)

func main() {
	ctx := context.Background()

	// Initialize with credentials and environment configuration
	opt := twigga.WithCredentialsFile("./bongo.json")
	env := twigga.WithEnvironment(twigga.TEST)

	app, err := twigga.InitializeApp(ctx, opt, env)
	if err != nil {
		log.Fatalf("Error initializing Twigga App: %v\n", err)
	}
	fmt.Println("SDK Initialized Successfully!")
}

```

## Authentication
Manage user identities, sessions, and password recovery natively.

clientId := "my_mobile_app_123"

1. Sign Up a New User
```
tokens, err := app.Auth().Signup(ctx, twigga.CreateUserReq{
    Email:     "user@example.com",
    Password:  "SuperSecret123!",
    FirstName: "John",
    LastName:  "Doe",
    ClientId:  clientId,
})

```

2. Refresh Token in the background
```
newTokens, _ := app.Auth().RefreshToken(ctx, tokens.RefreshToken)
```

3. Password Recovery Flow
```
app.Auth().ForgotPassword(ctx, "user@example.com")
app.Auth().ResetPassword(ctx, "token_from_email", "newPassword")
```

4. Logout (Revoke Session)
```
app.Auth().Logout(ctx, newTokens.RefreshToken)
```


## Database
Interact with your NoSQL document database. You can query, add, set, and delete documents easily.

Add a new document (Database Auto-Generates ID)
```
app.Database().Collection("users").Add(ctx, map[string]any{
    "name": "mtata",
    "role": "chair",
})

```

Set a document with a specific custom ID

```
app.Database().Collection("users").Doc("dev_123").Set(ctx, map[string]any{
    "name":   "mtata",
    "status": "active",
})

```

Retrieve a document
```
userDoc, err := app.Database().Collection("users").Doc("dev_123").Get(ctx)
if err == nil {
    fmt.Printf("User Retrieved: %s\n", string(userDoc))
}

```

## Zanzibar Authorization (ReBAC)
Manage granular, relationship-based access control (ReBAC) using the Zanzibar engine. Define complex ownership and permission structures.

Write a new permission tuple: "dev_123 is the owner of doc_999"

```
app.Auth().Zanzibar().Write(ctx, twigga.AuthorizationTuple{
    ObjectType:  "document",
    ObjectID:    "doc_999",
    Relation:    "owner",
    SubjectType: "user",
    SubjectID:   "dev_123",
})

```
// Check permission in real-time

```
allowed, err := app.Auth().Zanzibar().Check(ctx, "dev_123", "user", "owner", "document", "doc_999")
if err == nil {
    fmt.Printf("Is 'dev_123' the owner of 'doc_999'? %v\n", allowed)
}

```

## Serverless Functions
Invoke your deployed backend serverless functions directly from the SDK.

```
payload := map[string]any{
    "to": "user@example.com", 
    "subject": "Welcome to BongoCloud",
}

result, err := app.Functions().Invoke(ctx, "sendWelcomeEmail", payload)
if err == nil {
    fmt.Printf("Function Invoked! Result: %s\n", string(result))
}

```

## Cloud Storage
Manage buckets, set access policies, upload files, and retrieve secure download URLs.

bucketName := "profile-pictures"

1. Create a Bucket & Make it Public
```
app.Storage().CreateBucket(ctx, bucketName)
app.Storage().Bucket(bucketName).SetPolicy(ctx, "public")
```

2. Upload a File
```
app.Storage().Bucket(bucketName).Upload(ctx, "./test.txt")
```
3. Get File Metadata & Download URL
```
fileData, err := app.Storage().Bucket(bucketName).File("test.txt").Get(ctx)
if err == nil {
    fmt.Printf("File Download URL: %v\n", fileData["url"])
}

```

## Backups & Disaster Recovery
Trigger manual snapshots of your Twigga infrastructure for disaster recovery.

Trigger an incremental disaster recovery backup
```
app.Backups().Create(ctx, "incremental")
fmt.Println("Incremental Backup Triggered!")
```
