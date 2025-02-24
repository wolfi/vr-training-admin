// internal/handlers/avatars.go
package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/saladinomario/vr-training-admin/internal/models"
	"github.com/saladinomario/vr-training-admin/templates/components/avatars"
	"github.com/saladinomario/vr-training-admin/templates/pages"
)

var avatarStore *models.AvatarStore

func init() {
	avatarStore = models.NewAvatarStore()
}

// AvatarsHandler handles the avatars index page
func AvatarsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/avatars" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allAvatars := avatarStore.GetAll()
	component := pages.AvatarsIndex(allAvatars)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering avatars page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AvatarNewHandler handles the new avatar form
func AvatarNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	component := pages.AvatarNew()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering new avatar form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AvatarEditHandler handles the edit avatar form
func AvatarEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/avatars/edit/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// Get avatar by ID
	avatar, err := avatarStore.GetByID(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	component := pages.AvatarEdit(avatar)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering edit avatar form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AvatarCreateHandler handles avatar creation
func AvatarCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	avatar := parseAvatarForm(r)

	// Create avatar
	_, err := avatarStore.Create(avatar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Redirect to avatars page
	http.Redirect(w, r, "/avatars", http.StatusSeeOther)
}

// AvatarUpdateHandler handles avatar updates
func AvatarUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/avatars/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	avatar := parseAvatarForm(r)

	// Update avatar
	err := avatarStore.Update(idStr, avatar)
	if err != nil {
		if err == models.ErrAvatarNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Redirect to avatars page
	http.Redirect(w, r, "/avatars", http.StatusSeeOther)
}

// AvatarDeleteHandler handles avatar deletion
func AvatarDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/avatars/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// Delete avatar
	err := avatarStore.Delete(idStr)
	if err != nil {
		if err == models.ErrAvatarNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If this is an HTMX request, return the updated avatar list
	if r.Header.Get("HX-Request") == "true" {
		allAvatars := avatarStore.GetAll()
		component := avatars.AvatarList(allAvatars)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering avatar list: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Standard redirect for non-HTMX requests
	http.Redirect(w, r, "/avatars", http.StatusSeeOther)
}

// AvatarSearchHandler handles avatar search
func AvatarSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	foundAvatars := avatarStore.Search(query)

	component := avatars.AvatarList(foundAvatars)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering avatar search results: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Helper function to parse avatar form data
func parseAvatarForm(r *http.Request) avatars.Avatar {
	knowledgeLevel, _ := strconv.Atoi(r.FormValue("knowledge_level"))
	aggressivenessLevel, _ := strconv.Atoi(r.FormValue("aggressiveness_level"))
	patienceLevel, _ := strconv.Atoi(r.FormValue("patience_level"))
	emotionalReactivity, _ := strconv.Atoi(r.FormValue("emotional_reactivity"))
	speakingSpeed, _ := strconv.Atoi(r.FormValue("speaking_speed"))

	return avatars.Avatar{
		Name:                r.FormValue("name"),
		Description:         r.FormValue("description"),
		PersonalityType:     r.FormValue("personality_type"),
		CommunicationStyle:  r.FormValue("communication_style"),
		KnowledgeLevel:      knowledgeLevel,
		AggressivenessLevel: aggressivenessLevel,
		PatienceLevel:       patienceLevel,
		EmotionalReactivity: emotionalReactivity,
		VoiceType:           r.FormValue("voice_type"),
		SpeakingSpeed:       speakingSpeed,
		ImageURL:            r.FormValue("image_url"),
		Keywords:            r.FormValue("keywords"),
	}
}

// SetupAvatarRoutes registers all avatar-related routes
func SetupAvatarRoutes(mux *http.ServeMux) {
	log.Println("Setting up avatar routes...")

	// List and Create
	mux.HandleFunc("/avatars", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling avatar request: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			AvatarsHandler(w, r)
		case http.MethodPost:
			AvatarCreateHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// New form
	mux.HandleFunc("/avatars/new", AvatarNewHandler)

	// Search
	mux.HandleFunc("/avatars/search", AvatarSearchHandler)

	// Edit form
	mux.HandleFunc("/avatars/edit/", AvatarEditHandler)

	// Update and Delete
	mux.HandleFunc("/avatars/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodPost:
			AvatarUpdateHandler(w, r)
		case http.MethodDelete:
			AvatarDeleteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Avatar routes registered successfully")
}
