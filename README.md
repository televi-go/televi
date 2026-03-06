
# 🤖 Televi - Build Telegram Bots Like React

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Write Telegram bots using a declarative, component-based architecture inspired by React.**

Televi is a Go framework that brings React-like patterns to Telegram bot development. Define your bot's UI as composable **Scenes** and **Views**, manage state reactively, and let the framework handle the complexity of Telegram's API.

---

## ✨ Features

- 🧩 **Component-Based Architecture**: Build bots using reusable `Scene` and `View` components
- ⚛️ **React-Inspired API**: Familiar patterns like `State[T]`, context propagation, and declarative rendering
- 🔄 **Reactive State Management**: Automatic re-renders when state changes via `State[T].Set()`
- 🎨 **Rich Media Support**: Send text, photos, animations, keyboards, and inline buttons
- 🧭 **Scene Navigation**: Manage multi-step flows with transition policies (`SeparativeTransition`, `ReplacingTransition`)
- 📊 **Built-in Metrics**: Integrated profiling and metrics server for monitoring bot performance
- 🧵 **Concurrency-Safe**: Designed for high-throughput bots with per-user controller isolation
- 📦 **Asset Loading**: Preload and embed images/animations with `ImageAssetGroupLoader`

---

## 🚀 Quick Start

### Installation

```bash
go get github.com/televi-go/televi@latest
```

### Basic Example

```go
package main

import (
    "context"
    "log"
    "github.com/televi-go/televi"
    "github.com/televi-go/televi/models/pages"
)

// Define a Scene (like a React component)
type HelloScene struct {
    Count televi.State[int] // Reactive state
}

// Implement the View method to render your UI
func (s *HelloScene) View(ctx televi.BuildContext) {
    // Access state
    count := s.Count.Get()
    
    // Render text message
    ctx.Text("Hello! You've clicked: %d times", count)
    
    // Add an inline button that updates state
    ctx.InlineKeyboard().
        Button("Click me!", func() {
            s.Count.Set(count + 1) // Triggers re-render
        })
}

func main() {
    // Initialize your bot
    app, err := televi.NewApp(
        "YOUR_BOT_TOKEN",           // Telegram Bot API token
        "https://api.telegram.org", // API endpoint
        func(platform televi.Platform) televi.ActionScene {
            return &HelloScene{
                Count: televi.StateOf(0), // Initialize state
            }
        },
        nil, // Optional metrics server
    )
    if err != nil {
        log.Fatal(err)
    }

    // Start the bot
    app.Run(context.Background())
}
```

---

## 📚 Core Concepts

### Scenes & Views

```go
type Scene interface {
    View(ctx BuildContext)
}

type View interface {
    Build(ctx PageBuildContext)
}
```

- **Scene**: A top-level component representing a bot "screen" or conversation state
- **View**: A composable UI element that builds part of a message

### Reactive State

```go
type State[T any] struct { /* ... */ }

// Usage:
type CounterScene struct {
    Value televi.State[int]
}

func (s *CounterScene) View(ctx televi.BuildContext) {
    // Read state
    val := s.Value.Get()
    
    // Update state (triggers re-render)
    s.Value.Set(val + 1)
    
    // Or use functional updates
    s.Value.SetFn(func(prev int) int { return prev + 1 })
}
```

### BuildContext API

The `BuildContext` provides methods to construct Telegram messages:

```go
ctx.Text(format string, args ...any)              // Send text
ctx.Photo(asset televi.ImageAsset)                // Send photo
ctx.Animation(asset televi.ImageAsset)            // Send animation
ctx.InlineKeyboard()                              // Create inline keyboard
ctx.ReplyKeyboard()                               // Create reply keyboard
ctx.Navigate(scene Scene, policy TransitPolicy)   // Transition to new scene
```

### Asset Management

```go
// Preload images for efficient sending
loader := televi.NewAssetLoader()
var logo televi.ImageAsset

err := loader.
    Add("assets/logo.png", &logo).
    Load()
if err != nil {
    log.Fatal(err)
}

// Use in your scene
func (s *MyScene) View(ctx televi.BuildContext) {
    logo.Embed(ctx.PhotoConsumer()) // Send preloaded image
}
```

---

## 🗂️ Project Structure

```
televi/
├── core/              # Core framework: App, Controller, Scene management
├── models/pages/      # Page building primitives: BuildContext, State, transitions
├── telegram/          # Telegram API client and DTOs
├── connector/         # UI element connectors: keyboards, text, photos
├── examples/          # Working bot examples
├── menu_bot/          # Reference implementation
└── exports.go         # Public API surface
```

---

## 🔧 Advanced Usage

### Scene Transitions

```go
type WizardScene struct {
    Step televi.State[int]
}

func (s *WizardScene) View(ctx televi.BuildContext) {
    switch s.Step.Get() {
    case 0:
        ctx.Text("Step 1: Enter your name")
        // ... handle input
    case 1:
        ctx.Text("Step 2: Confirm details")
        ctx.InlineKeyboard().
            Button("Back", func() {
                ctx.Navigate(s, pages.ReplacingTransition) // Go to previous
            }).
            Button("Finish", func() {
                ctx.Navigate(&FinalScene{}, pages.SeparativeTransition)
            })
    }
}
```

### Per-User Controllers

Televi automatically creates isolated `Controller` instances per chat, maintaining scene state and message history independently for each user.

### Metrics & Profiling

```go
// Enable built-in metrics server
server := metrics.NewServerImpl(dbConnection)

app, _ := televi.NewApp(token, endpoint, initScene, server)

// Access throughput stats
app.Profiler.WriteStats(os.Stdout)
```

---

## 🧪 Testing

```go
// Use the test runner for scene logic
import "github.com/televi-go/televi/runner/test"

func TestHelloScene(t *testing.T) {
    runner := test.NewRunner(&HelloScene{})
    
    // Simulate user interaction
    runner.Dispatch(&dto.Update{
        Message: &dto.Message{Text: "/start"},
    })
    
    // Assert expected output
    // ...
}
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feat/amazing-feature`
5. Open a Pull Request

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

## ⚠️ Project Status

**Note**: This project is discontinued in may 2023

---
