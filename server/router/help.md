**Path Parameters:**
```bash
curl http://localhost:8080/users/123
curl http://localhost:8080/users/123/posts/456
```

**Query Parameters:**
```bash
curl "http://localhost:8080/search?q=golang&page=2&sort=name"
curl "http://localhost:8080/advanced-search?tags=go&tags=web&filters[status]=active&min_price=10.5"
```

**Headers:**
```bash
curl -H "Authorization: Bearer token123" -H "X-API-Key: secret" http://localhost:8080/headers
curl -X POST -H "Authorization: Bearer token123" http://localhost:8080/auth
```

**JSON Body:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","age":25,"password":"secret123"}' \
  http://localhost:8080/users
```

**Form Data:**
```bash
curl -X POST -d "username=john&password=secret123" http://localhost:8080/login
curl -X POST -d "username=john&hobbies=coding&hobbies=gaming&preferences[theme]=dark" http://localhost:8080/register
```

**File Upload:**
```bash
curl -X POST -F "file=@example.txt" http://localhost:8080/upload
curl -X POST -F "name=John" -F "avatar=@avatar.jpg" http://localhost:8080/form-with-file
```

**Raw Body:**
```bash
curl -X POST -H "Content-Type: text/plain" -d "raw text data" http://localhost:8080/webhook
```

**Cookies:**
```bash
curl http://localhost:8080/set-cookie
curl -b "session_id=abc123" http://localhost:8080/get-cookie
```


```go
// 1. PATH PARAMETERS
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		
		// Convert to int if needed
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"message": fmt.Sprintf("Getting user with ID: %d", userID),
		})
	})

	r.GET("/users/:id/posts/:postId", func(c *gin.Context) {
		userID := c.Param("id")
		postID := c.Param("postId")

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"post_id": postID,
			"message": fmt.Sprintf("User %s, Post %s", userID, postID),
		})
	})

	// 2. QUERY PARAMETERS
	r.GET("/search", func(c *gin.Context) {
		// Single query parameter
		query := c.Query("q")                    // Returns empty string if not found
		page := c.DefaultQuery("page", "1")      // Returns default if not found
		limit := c.DefaultQuery("limit", "10")

		// Check if query parameter exists
		sort, exists := c.GetQuery("sort")
		if !exists {
			sort = "created_at"
		}

		c.JSON(http.StatusOK, gin.H{
			"query": query,
			"page":  page,
			"limit": limit,
			"sort":  sort,
		})
	})

	r.GET("/advanced-search", func(c *gin.Context) {
		// Query arrays: ?tags=go&tags=web&tags=api
		tags := c.QueryArray("tags")
		
		// Query map: ?filters[status]=active&filters[type]=premium
		filters := c.QueryMap("filters")

		// Convert query params to specific types
		minPrice, _ := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
		maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "1000"), 64)
		
		isActive, _ := strconv.ParseBool(c.DefaultQuery("active", "true"))

		c.JSON(http.StatusOK, gin.H{
			"tags":      tags,
			"filters":   filters,
			"min_price": minPrice,
			"max_price": maxPrice,
			"is_active": isActive,
		})
	})

	// 3. HEADERS
	r.GET("/headers", func(c *gin.Context) {
		// Get specific headers
		userAgent := c.GetHeader("User-Agent")
		contentType := c.GetHeader("Content-Type")
		authorization := c.GetHeader("Authorization")
		
		// Get custom headers
		apiKey := c.GetHeader("X-API-Key")
		requestID := c.GetHeader("X-Request-ID")

		// Get all headers
		headers := c.Request.Header

		c.JSON(http.StatusOK, gin.H{
			"user_agent":    userAgent,
			"content_type":  contentType,
			"authorization": authorization,
			"api_key":       apiKey,
			"request_id":    requestID,
			"all_headers":   headers,
		})
	})

	r.POST("/auth", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authenticated",
			"token":   token,
		})
	})

	// 4. JSON BODY
	r.POST("/users", func(c *gin.Context) {
		var user User
		
		// ShouldBindJSON for JSON
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user":    user,
		})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("User %s updated", id),
			"user":    user,
		})
	})

	r.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updates UpdateUser
		
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("User %s partially updated", id),
			"updates": updates,
		})
	})

	// 5. FORM DATA (application/x-www-form-urlencoded)
	r.POST("/login", func(c *gin.Context) {
		var form LoginForm
		
		// ShouldBind automatically detects content type
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Login successful",
			"username": form.Username,
		})
	})

	r.POST("/register", func(c *gin.Context) {
		// Manual form parsing
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.DefaultPostForm("email", "")
		
		// Check if form field exists
		newsletter, exists := c.GetPostForm("newsletter")
		subscribeNewsletter := exists && newsletter == "on"

		// Form arrays: name="hobbies" multiple values
		hobbies := c.PostFormArray("hobbies")
		
		// Form map: preferences[theme]=dark&preferences[lang]=en
		preferences := c.PostFormMap("preferences")

		c.JSON(http.StatusOK, gin.H{
			"username":            username,
			"email":               email,
			"subscribe_newsletter": subscribeNewsletter,
			"hobbies":             hobbies,
			"preferences":         preferences,
		})
	})

	// 6. MULTIPART FORM DATA
	r.POST("/upload", func(c *gin.Context) {
		// Single file upload
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Save file
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"filename": file.Filename,
			"size":     file.Size,
			"saved_as": filename,
		})
	})

	r.POST("/upload-multiple", func(c *gin.Context) {
		// Multiple file upload
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		files := form.File["files"]
		var uploadedFiles []string

		for _, file := range files {
			filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + file.Filename})
				return
			}
			uploadedFiles = append(uploadedFiles, filename)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Files uploaded successfully",
			"uploaded_files": uploadedFiles,
			"count":          len(uploadedFiles),
		})
	})

	r.POST("/form-with-file", func(c *gin.Context) {
		// Mix of form data and file upload
		name := c.PostForm("name")
		description := c.PostForm("description")
		
		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file required"})
			return
		}

		filename := fmt.Sprintf("avatars/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name":        name,
			"description": description,
			"avatar":      filename,
		})
	})

	// 7. RAW BODY
	r.POST("/webhook", func(c *gin.Context) {
		// Read raw body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}

		contentType := c.GetHeader("Content-Type")

		c.JSON(http.StatusOK, gin.H{
			"content_type": contentType,
			"body_length":  len(body),
			"body":         string(body),
		})
	})

	r.POST("/xml", func(c *gin.Context) {
		var data map[string]interface{}
		
		// ShouldBindXML for XML data
		if err := c.ShouldBindXML(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "XML data received",
			"data":    data,
		})
	})

	// 8. COOKIES
	r.GET("/set-cookie", func(c *gin.Context) {
		// Set cookie
		c.SetCookie(
			"session_id",           // name
			"abc123",               // value
			3600,                   // maxAge (seconds)
			"/",                    // path
			"localhost",            // domain
			false,                  // secure
			true,                   // httpOnly
		)

		c.JSON(http.StatusOK, gin.H{"message": "Cookie set successfully"})
	})

	r.GET("/get-cookie", func(c *gin.Context) {
		// Get cookie
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Session cookie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"session_id": sessionID,
			"message":    "Cookie retrieved successfully",
		})
	})

	// 9. COMBINATION REQUESTS
	r.POST("/complex/:id", func(c *gin.Context) {
		// Path parameter
		userID := c.Param("id")
		
		// Query parameters
		include := c.QueryArray("include")
		format := c.DefaultQuery("format", "json")
		
		// Headers
		authorization := c.GetHeader("Authorization")
		userAgent := c.GetHeader("User-Agent")
		
		// JSON body
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Cookie
		sessionID, _ := c.Cookie("session_id")

		c.JSON(http.StatusOK, gin.H{
			"user_id":       userID,
			"include":       include,
			"format":        format,
			"authorization": authorization,
			"user_agent":    userAgent,
			"session_id":    sessionID,
			"request_data":  requestData,
		})
	})

	// 10. REQUEST VALIDATION
	r.POST("/validate", func(c *gin.Context) {
		var user User
		
		// Use ShouldBindJSON with validation tags
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		// Additional custom validation
		if user.Age < 18 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User must be at least 18 years old",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User data is valid",
			"user":    user,
		})
	})

	// 11. CONTENT TYPE SPECIFIC BINDINGS
	r.POST("/bind-query", func(c *gin.Context) {
		var query struct {
			Name  string `form:"name" binding:"required"`
			Email string `form:"email" binding:"required,email"`
		}
		
		// Bind only query parameters
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Query bound successfully",
			"data":    query,
		})
	})

	r.GET("/bind-header", func(c *gin.Context) {
		var headers struct {
			Authorization string `header:"Authorization" binding:"required"`
			UserAgent     string `header:"User-Agent"`
			RequestID     string `header:"X-Request-ID"`
		}
		
		// Bind headers to struct
		if err := c.ShouldBindHeader(&headers); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Headers bound successfully",
			"headers": headers,
		})
	})

	// 12. CLIENT IP AND REQUEST INFO
	r.GET("/request-info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"client_ip":     c.ClientIP(),
			"method":        c.Request.Method,
			"url":           c.Request.URL.String(),
			"user_agent":    c.Request.UserAgent(),
			"content_length": c.Request.ContentLength,
			"host":          c.Request.Host,
			"remote_addr":   c.Request.RemoteAddr,
		})
	})

	// 13. CUSTOM RESPONSE TYPES
	r.GET("/text", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.GET("/html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://google.com")
	})

	r.GET("/data", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/octet-stream", []byte("binary data"))
	})

```