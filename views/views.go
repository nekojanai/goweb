package views

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"sleepy.systems/goweb/config"
	"sleepy.systems/goweb/utils"
)

const VIEWS_PATH = "views/html/*"

type Page struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Body      []byte
	FilePath  string
}

func (page *Page) SetID(id string) *Page {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		log.Fatal(err)
	}
	page.ID = parsedUUID
	return page
}

func (page *Page) init() {
	if page.ID == uuid.Nil {
		page.ID = uuid.New()
	}

	if page.CreatedAt.IsZero() {
		page.CreatedAt = time.Now()
	}

	page.UpdatedAt = time.Now()
}

func createFilePath(uuid string, config *config.Config) string {
	return fmt.Sprintf("%v%v%v", config.DataPath, uuid, ".md")
}

func New(title string, body []byte) *Page {
	page := Page{Title: title, Body: body}
	page.init()
	return &page
}

func (page *Page) Save(config *config.Config) {
	buffer := new(bytes.Buffer)

	page.init()

	err := toml.NewEncoder(buffer).Encode(map[string]interface{}{
		"ID":        page.ID,
		"CreatedAt": page.CreatedAt,
		"UpdatedAt": page.UpdatedAt,
		"Title":     page.Title,
		"Body":      page.Body,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(createFilePath(page.ID.String(), config), buffer.Bytes(), 0650)
	if err != nil {
		log.Fatal(err)
	}
}

func (page *Page) Exists(config *config.Config) bool {
	_, err := os.Stat(createFilePath(page.ID.String(), config))
	return err != nil
}

func (page *Page) Hash(config *config.Config) (string, error) {
	file, err := os.Open(createFilePath(page.ID.String(), config))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sha256 := crypto.SHA256.New()

	if _, err := io.Copy(sha256, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(sha256.Sum(nil)), nil
}

func (page *Page) Load(config *config.Config) (*Page, error) {
	_, err := toml.DecodeFile(createFilePath(page.ID.String(), config), page)
	return page, err
}

// TODO: continue here

func LoadAll(config config.Config) []Page {
	//r.URL.Path[len("/view/"):]
	return []Page{}
}

func parseAndExecuteTemplate(w io.Writer, tmpl string, data any) {
	t, err := template.ParseGlob(VIEWS_PATH)
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, tmpl, data)
}

func HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	utils.LogReq(r)
	parseAndExecuteTemplate(w, "index", nil)
}

func HandleSubPage(w http.ResponseWriter, r *http.Request, config *config.Config) {
	page, err := (&Page{}).SetID(r.PathValue("id")).Load(config)
	utils.ErrorHandler(w, r, http.StatusNotFound, err)

	utils.LogReq(r)
	parseAndExecuteTemplate(w, "page", nil)
}
