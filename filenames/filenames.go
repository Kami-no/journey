package filenames

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kabukky/journey/flags"
)

var (
	// Determine the work path for built-in resources - needed to load relative assets
	WorkPath = determineWorkPath()

	// Determine the path to the variable assets directory (default: Journey work directory)
	AssetPath = determineAssetPath()

	// For assets that are created, changed, our user-provided while running journey
	ConfigFilename   = filepath.Join(AssetPath, "config.yaml")
	ContentFilepath  = filepath.Join(AssetPath, "content")
	DatabaseFilepath = filepath.Join(ContentFilepath, "data")
	DatabaseFilename = filepath.Join(ContentFilepath, "data", "journey.db")
	ThemesFilepath   = filepath.Join(ContentFilepath, "themes")
	ImagesFilepath   = filepath.Join(ContentFilepath, "images")
	PluginsFilepath  = filepath.Join(ContentFilepath, "plugins")
	PagesFilepath    = filepath.Join(ContentFilepath, "pages")
	StaticFilepath   = filepath.Join(ContentFilepath, "static")

	// For https
	HttpsFilepath     = filepath.Join(ContentFilepath, "https")
	HttpsCertFilename = filepath.Join(ContentFilepath, "https", "cert.pem")
	HttpsKeyFilename  = filepath.Join(ContentFilepath, "https", "key.pem")

	// For built-in files (e.g. the admin interface)
	AdminFilepath  = filepath.Join(WorkPath, "built-in", "admin")
	PublicFilepath = filepath.Join(WorkPath, "built-in", "public")
	HbsFilepath    = filepath.Join(WorkPath, "built-in", "hbs")

	// For blog  (this is a url string)
	// TODO: This is not used at the moment because it is still hard-coded into the create database string
	DefaultBlogLogoFilename  = "/public/images/blog-logo.jpg"
	DefaultBlogCoverFilename = "/public/images/blog-cover.jpg"

	// For users (this is a url string)
	DefaultUserImageFilename = "/public/images/user-image.jpg"
	DefaultUserCoverFilename = "/public/images/user-cover.jpg"
)

func init() {
	// Create content directories if they are not created already
	err := createDirectories()
	if err != nil {
		log.Fatal("Error: Couldn't create directories:", err)
	}

}

func createDirectories() error {
	paths := []string{DatabaseFilepath, ThemesFilepath, ImagesFilepath, HttpsFilepath, PluginsFilepath, PagesFilepath}
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Println("Creating " + path)
			err := os.MkdirAll(path, 0776)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func bashPath(contentPath string) string {
	if strings.HasPrefix(contentPath, "~") {
		return strings.Replace(contentPath, "~", os.Getenv("HOME"), 1)
	}
	return contentPath
}

func determineAssetPath() string {
	if flags.CustomPath != "" {
		contentPath, err := filepath.Abs(bashPath(flags.CustomPath))
		if err != nil {
			log.Fatal("Error: Couldn't read from custom path: ", err)
		}
		return contentPath
	}
	return determineWorkPath()
}

func determineWorkPath() string {
	// Get the path of work directory
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("Error: Couldn't determine what work directory: ", err)
	}
	path = filepath.Clean(path)
	return path
}
