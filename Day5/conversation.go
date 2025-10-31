

package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// 🧱 Conversation represents a conversation topic
type Conversation struct {
    ID       uint      `gorm:"primaryKey" json:"id"`
    Topic    string    `gorm:"unique;not null" json:"topic"`
    Messages []Message `json:"messages"`
}

// 💬 Message represents a single chat message
type Message struct {
    ID             uint   `gorm:"primaryKey" json:"id"`
    ConversationID uint   `json:"conversation_id"`
    Content        string `gorm:"type:text;not null" json:"content"`
}

var db *gorm.DB

// 🧩 Initialize MySQL + AutoMigrate
func initDB() {
    dsn := "root:jvt123@tcp(127.0.0.1:3306)/familydb?charset=utf8mb4&parseTime=True&loc=Local"
    var err error

    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("❌ Failed to connect to database:", err)
    }

    err = db.AutoMigrate(&Conversation{}, &Message{})
    if err != nil {
        log.Fatal("❌ AutoMigrate failed:", err)
    }

    log.Println("✅ Connected to MySQL & Migrated schema")
}

// 🧠 Create a new conversation
func createConversation(c *gin.Context) {
    var conv Conversation
    if err := c.ShouldBindJSON(&conv); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if result := db.Create(&conv); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Conversation created successfully", "data": conv})
}

// 📚 Get a conversation (with messages)
func readConversation(c *gin.Context) {
    topic := c.Param("topic")

    var conv Conversation
    if err := db.Preload("Messages").Where("topic = ?", topic).First(&conv).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }

    c.JSON(http.StatusOK, conv)
}

// ✏️ Add new message to a conversation
func addMessage(c *gin.Context) {
    topic := c.Param("topic")

    var conv Conversation
    if err := db.Where("topic = ?", topic).First(&conv).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }

    var input struct {
        Message string `json:"message"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newMessage := Message{
        ConversationID: conv.ID,
        Content:        input.Message,
    }
    db.Create(&newMessage)

    c.JSON(http.StatusOK, gin.H{"message": "Message added", "data": newMessage})
}

// 🗑️ Delete a conversation
func deleteConversation(c *gin.Context) {
    topic := c.Param("topic")

    var conv Conversation
    if err := db.Where("topic = ?", topic).First(&conv).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }

    // Delete messages first (to maintain referential integrity)
    db.Where("conversation_id = ?", conv.ID).Delete(&Message{})
    db.Delete(&conv)

    c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted successfully"})
}

// 🧾 List all conversations
func listConversations(c *gin.Context) {
    var convs []Conversation
    db.Find(&convs)
    c.JSON(http.StatusOK, gin.H{"conversations": convs})
}

// 🚀 Main entry point
func main() {
    initDB()

    router := gin.Default()

    router.POST("/conversations", createConversation)
    router.GET("/conversations", listConversations)
    router.GET("/conversations/:topic", readConversation)
    router.PUT("/conversations/:topic", addMessage)
    router.DELETE("/conversations/:topic", deleteConversation)

    log.Println("🚀 Wife-Husband Conversation API running on http://localhost:8080")
    router.Run(":8080")
}

/*

Perfect — you’re very close! 🧠
Your API is already **running** at `http://localhost:8080`, so you can absolutely test it using **Google Chrome** (or any browser).

Let’s go step-by-step 👇

---

## 🟢 STEP 1: Start the server

Make sure your Go API is running in the terminal:

```
🚀 Wife-Husband Conversation API running on http://localhost:8080
```

If you see that, ✅ your backend is live.

---

## 🟡 STEP 2: Test simple GET routes in Chrome

### 🧩 1️⃣ View all conversations

Open this URL in Chrome:

```
http://localhost:8080/conversations
```

If your DB is empty, you’ll see:

```json
{"conversations":[]}
```

Otherwise, you’ll see something like:

```json
{"conversations":[{"id":1,"topic":"weekend"}]}
```

---

### 🧩 2️⃣ View a specific conversation

For example, open:

```
http://localhost:8080/conversations/weekend
```

If it exists, you’ll see:

```json
{
  "id": 1,
  "topic": "weekend",
  "messages": [
    { "id": 1, "conversation_id": 1, "content": "Let’s go out!" }
  ]
}
```

If it doesn’t exist yet, you’ll see:

```json
{"error": "Topic not found"}
```

---

## 🔵 STEP 3: Test POST / PUT / DELETE routes

Those routes **require sending JSON data**, so Chrome alone can’t do that easily.
You have 3 easy options:

---

### ✅ Option 1: Use **Postman** (recommended)

Download [Postman](https://www.postman.com/downloads/) — a free API testing tool.

Then:

#### Create conversation:

* **Method:** `POST`
* **URL:** `http://localhost:8080/conversations`
* **Body → raw → JSON:**

  ```json
  {
    "topic": "weekend"
  }
  ```

Click **Send** → you’ll get a success JSON.

---

#### Add message:

* **Method:** `PUT`
* **URL:** `http://localhost:8080/conversations/weekend`
* **Body → raw → JSON:**

  ```json
  {
    "message": "Let’s go for dinner!"
  }
  ```

Click **Send** → message gets added.

---

#### Delete conversation:

* **Method:** `DELETE`
* **URL:** `http://localhost:8080/conversations/weekend`

---

### ✅ Option 2: Use a **Chrome Extension**

Install a simple “REST API Client” extension, like:

* 🧩 **Talend API Tester**
* 🧩 **RESTED**
* 🧩 **Hoppscotch (Web)** → [https://hoppscotch.io](https://hoppscotch.io)

Then use the same requests as above.

---

### ✅ Option 3: Use **curl** in your terminal (simple)

Example:

```bash
curl -X POST http://localhost:8080/conversations -H "Content-Type: application/json" -d "{\"topic\": \"weekend\"}"
```

---

### ✅ Quick Summary

| Task                 | Method                         | Chrome works directly? | Use Postman instead? |
| -------------------- | ------------------------------ | ---------------------- | -------------------- |
| List conversations   | GET `/conversations`           | ✅ Yes                  | Optional             |
| Get one conversation | GET `/conversations/:topic`    | ✅ Yes                  | Optional             |
| Create conversation  | POST `/conversations`          | ❌ No                   | ✅ Yes                |
| Add message          | PUT `/conversations/:topic`    | ❌ No                   | ✅ Yes                |
| Delete conversation  | DELETE `/conversations/:topic` | ❌ No                   | ✅ Yes                |

---

Would you like me to give you a **ready-made Postman collection** (you can import it and test all endpoints instantly)?


*/