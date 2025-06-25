package discordgo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// SlashHandlerFunc is what your command handlers should look like.
// Think of it as: "When someone uses this slash command, run this function"
type SlashHandlerFunc func(ctx context.Context, s *Session, i *InteractionCreate)

// MiddlewareFunc wraps around your handlers to add extra functionality.
// Like adding logging, permission checks, or rate limiting before the actual command runs
type MiddlewareFunc func(SlashHandlerFunc) SlashHandlerFunc

// SlashRouter is your command center for managing Discord slash commands.
// It keeps track of all your commands, applies middleware, and routes interactions
// to the right handlers.
type SlashRouter struct {
	mu           sync.RWMutex              // Protects concurrent access to our maps
	commands     map[string]SlashHandlerFunc // Maps command names to their handlers
	middlewares  []MiddlewareFunc          // Stack of middleware that runs before commands
	errorHandler SlashHandlerFunc          // What to do when a command isn't found
	
	// New optimization: command metadata for better debugging
	commandMeta  map[string]*CommandMetadata
}

// CommandMetadata holds useful info about each registered command
type CommandMetadata struct {
	Name        string
	Description string
	RegisteredAt time.Time
	UsageCount  int64
}

// NewSlashRouter creates a fresh router ready to handle your slash commands.
// Think of this as setting up your command headquarters!
func NewSlashRouter() *SlashRouter {
	return &SlashRouter{
		commands:    make(map[string]SlashHandlerFunc),
		commandMeta: make(map[string]*CommandMetadata),
	}
}

// Handle registers a new slash command with the router.
// This is where you tell the router "when someone types /commandName, run this handler"
func (r *SlashRouter) Handle(name string, handler SlashHandlerFunc) {
	r.HandleWithDescription(name, handler, "")
}

// HandleWithDescription is like Handle but lets you add a description for debugging
func (r *SlashRouter) HandleWithDescription(name string, handler SlashHandlerFunc, description string) {
	if handler == nil {
		panic(fmt.Sprintf("handler for command '%s' cannot be nil", name))
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	normalizedName := strings.ToLower(strings.TrimSpace(name))
	if normalizedName == "" {
		panic("command name cannot be empty")
	}
	
	r.commands[normalizedName] = handler
	r.commandMeta[normalizedName] = &CommandMetadata{
		Name:         normalizedName,
		Description:  description,
		RegisteredAt: time.Now(),
		UsageCount:   0,
	}
}

// Use adds middleware to your router. Middleware runs in the order you add it.
// Think of it like a security checkpoint - each middleware can inspect or modify
// the request before it reaches your actual command handler.
func (r *SlashRouter) Use(middlewares ...MiddlewareFunc) {
	if len(middlewares) == 0 {
		return
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, middlewares...)
}

// SetErrorHandler tells the router what to do when someone uses a command that doesn't exist.
// By default, it just ignores unknown commands, but you might want to send a helpful message!
func (r *SlashRouter) SetErrorHandler(handler SlashHandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errorHandler = handler
}

// Dispatch is the heart of the router. When Discord sends an interaction,
// this figures out which command to run and applies all your middleware.
func (r *SlashRouter) Dispatch(s *Session, i *InteractionCreate) {
	// Only handle slash commands, ignore other interaction types
	if i.Type != InteractionApplicationCommand {
		return
	}

	cmdData := i.ApplicationCommandData()
	cmdName := strings.ToLower(strings.TrimSpace(cmdData.Name))

	// Find the handler for this command
	r.mu.RLock()
	handler, exists := r.commands[cmdName]
	
	// Update usage stats while we have the lock
	if exists && r.commandMeta[cmdName] != nil {
		r.commandMeta[cmdName].UsageCount++
	}
	r.mu.RUnlock()

	// If command doesn't exist, try the error handler
	if !exists {
		if r.errorHandler != nil {
			r.errorHandler(context.Background(), s, i)
		}
		return
	}

	// Build the middleware chain (like wrapping presents - outermost middleware runs first)
	finalHandler := handler
	for idx := len(r.middlewares) - 1; idx >= 0; idx-- {
		finalHandler = r.middlewares[idx](finalHandler)
	}

	// Finally, run the command with all middleware applied
	finalHandler(context.Background(), s, i)
}

// GetCommandStats returns usage statistics for debugging and monitoring
func (r *SlashRouter) GetCommandStats() map[string]*CommandMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Return a copy to prevent external modification
	stats := make(map[string]*CommandMetadata)
	for name, meta := range r.commandMeta {
		statsCopy := *meta // Copy the struct
		stats[name] = &statsCopy
	}
	return stats
}

// RegisteredCommands returns a list of all command names registered with this router
func (r *SlashRouter) RegisteredCommands() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	commands := make([]string, 0, len(r.commands))
	for name := range r.commands {
		commands = append(commands, name)
	}
	return commands
}

// RegisterCommands is a helper that registers multiple Discord application commands at once.
// This talks directly to Discord's API to make your commands available in servers.
func (r *SlashRouter) RegisterCommands(s *Session, guildID string, commands []*ApplicationCommand) ([]*ApplicationCommand, error) {
	if s == nil {
		return nil, fmt.Errorf("session cannot be nil")
	}
	if s.State.User == nil {
		return nil, fmt.Errorf("session user is nil - make sure the bot is logged in")
	}
	
	return s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, commands)
}

// --- Prebuilt Middleware ---
// These are common middleware patterns you can use right away!

// LoggingMiddleware logs every command usage. Super helpful for debugging!
// Just pass in a function that says how you want to log (could be to console, file, database, etc.)
func LoggingMiddleware(logFunc func(cmd string, userID string)) MiddlewareFunc {
	if logFunc == nil {
		// Provide a sensible default if no log function is given
		logFunc = func(cmd string, userID string) {
			log.Printf("Command '%s' used by user %s", cmd, userID)
		}
	}
	
	return func(next SlashHandlerFunc) SlashHandlerFunc {
		return func(ctx context.Context, s *Session, i *InteractionCreate) {
			cmdName := i.ApplicationCommandData().Name
			userID := "unknown"
			
			// Safely get user ID (could be from member or user depending on context)
			if i.Member != nil && i.Member.User != nil {
				userID = i.Member.User.ID
			} else if i.User != nil {
				userID = i.User.ID
			}
			
			logFunc(cmdName, userID)
			next(ctx, s, i)
		}
	}
}

// RequirePermission checks if the user has specific Discord permissions before running the command.
// Great for admin-only commands! Use Discord's permission constants like PermissionManageMessages.
func RequirePermission(requiredPerm int64) MiddlewareFunc {
	return func(next SlashHandlerFunc) SlashHandlerFunc {
		return func(ctx context.Context, s *Session, i *InteractionCreate) {
			// Get the user ID safely
			var userID string
			if i.Member != nil && i.Member.User != nil {
				userID = i.Member.User.ID
			} else if i.User != nil {
				userID = i.User.ID
			} else {
				respondWithError(s, i, "Could not identify user")
				return
			}
			
			// Check their permissions in this channel
			perms, err := s.UserChannelPermissions(userID, i.ChannelID)
			if err != nil {
				respondWithError(s, i, "Could not check permissions")
				return
			}
			
			// Do they have the required permission?
			if perms&requiredPerm != requiredPerm {
				respondWithError(s, i, "You don't have permission to use this command")
				return
			}
			
			// All good! Continue to the actual command
			next(ctx, s, i)
		}
	}
}

// RateLimitMiddleware prevents users from spamming commands.
// Simple implementation - you might want something more sophisticated for production.
func RateLimitMiddleware(maxUsesPerMinute int) MiddlewareFunc {
	userLastUse := make(map[string][]time.Time)
	var mu sync.Mutex
	
	return func(next SlashHandlerFunc) SlashHandlerFunc {
		return func(ctx context.Context, s *Session, i *InteractionCreate) {
			var userID string
			if i.Member != nil && i.Member.User != nil {
				userID = i.Member.User.ID
			} else if i.User != nil {
				userID = i.User.ID
			} else {
				next(ctx, s, i)
				return
			}
			
			mu.Lock()
			now := time.Now()
			
			// Clean up old entries (older than 1 minute)
			if times, exists := userLastUse[userID]; exists {
				var recentTimes []time.Time
				for _, t := range times {
					if now.Sub(t) < time.Minute {
						recentTimes = append(recentTimes, t)
					}
				}
				userLastUse[userID] = recentTimes
			}
			
			// Check if user has exceeded rate limit
			if len(userLastUse[userID]) >= maxUsesPerMinute {
				mu.Unlock()
				respondWithError(s, i, "Slow down! You're using commands too quickly.")
				return
			}
			
			// Record this usage
			userLastUse[userID] = append(userLastUse[userID], now)
			mu.Unlock()
			
			next(ctx, s, i)
		}
	}
}

// Helper function to send error responses consistently
func respondWithError(s *Session, i *InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &InteractionResponseData{
			Content: message,
			Flags:   MessageFlagsEphemeral, // Only the user who ran the command sees this
		},
	})
}
