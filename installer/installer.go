package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"tent/discord"
	"tent/logger"
	"time"

	"github.com/manifoldco/promptui"
)

func Selector() string {
	installs, err := discord.Discords()
	if err != nil {
		logger.Error(err)
		return ""
	}

	options := []string{}
	if installs.Discord {
		options = append(options, "Discord")
	}

	if installs.DiscordCanary {
		options = append(options, "Discord Canary")
	}

	if installs.DiscordPTB {
		options = append(options, "Discord PTB")
	}

	if len(options) == 0 {
		logger.Error("No Discord installations found")
		return ""
	}

	prompt := promptui.Select{
		Label: "Select Discord Version",
		Items: options,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | faint }}?",
			Active:   "{{ `>` | faint }} {{ . | white }}",
			Inactive: "  {{ . | faint }}",
			Selected: "{{ `✔` | green }} {{ . | white }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("Prompt failed: %v\n", err))
		return ""
	}

	return result
}

func Install(build string) {
	var err error
	switch build {
	case "Discord":
		err = install("Discord")
	case "Discord Canary":
		err = install("Discord Canary")
	case "Discord PTB":
		err = install("Discord PTB")
	default:
		logger.Error("Invalid selection")
		return
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Installation failed: %v\n", err))
	} else {
		logger.Success("Installation completed successfully!")
	}
}

func Discord()       { install("Discord") }
func DiscordCanary() { install("Discord Canary") }
func DiscordPTB()    { install("Discord PTB") }

func install(build string) (err error) {
	logger.Info("Installing Shelter for", build+"...")
	if build != "Discord" && build != "Discord Canary" && build != "Discord PTB" {
		logger.Error("Invalid build")
		return
	}
	no_spaces := strings.Join(strings.Split(build, " "), "")
	cache, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	install_folder := path.Join(cache, no_spaces)

	app_folder := ""
	files, err := os.ReadDir(install_folder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "app-") {
			app_folder = path.Join(install_folder, file.Name())
			break
		}
	}

	if app_folder == "" {
		return fmt.Errorf("could not find app folder")
	}

	resources_folder := path.Join(app_folder, "resources")
	resources_app_folder := path.Join(resources_folder, "app")

	exe_location := path.Join(app_folder, no_spaces+".exe")
	if _, err := os.Stat(exe_location); err != nil {
		return err
	}

	cmd := exec.Command("taskkill", "/IM", no_spaces+".exe", "/F")
	_, err = cmd.CombinedOutput()
	if err != nil {
		logger.Errorf("Failed to kill process: The process \"%s.exe\" is not running.\n", no_spaces)
	}

	time.Sleep(2 * time.Second)

	if _, err := os.Stat(resources_app_folder); err == nil {
		files, err := os.ReadDir(resources_app_folder)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return nil
		}

		logger.Infof("Resources app folder already exists. Press enter to clear it or any other key to exit: ")
		var input string
		fmt.Scanln(&input)
		if input == "" {
			err = os.RemoveAll(resources_app_folder)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	err = os.MkdirAll(resources_app_folder, 0755)
	if err != nil {
		return err
	}

	logger.Infof("Renaming %s's asar file...\n", build)
	asar_location := path.Join(resources_folder, "app.asar")
	original_asar_location := path.Join(resources_folder, "original.asar")
	err = os.Rename(asar_location, original_asar_location)
	if err != nil {
		return err
	}

	logger.Info("Downloading injectors...")
	indexjs_location := path.Join(resources_app_folder, "index.js")
	preloadjs_location := path.Join(resources_app_folder, "preload.js")
	packagejson_location := path.Join(resources_app_folder, "package.json")

	indexjs_link := "https://raw.githubusercontent.com/uwu/shelter/main/injectors/desktop/app/index.js"
	preloadjs_link := "https://raw.githubusercontent.com/uwu/shelter/main/injectors/desktop/app/preload.js"
	packagejson_link := "https://raw.githubusercontent.com/uwu/shelter/main/injectors/desktop/app/package.json"

	err = download(indexjs_link, indexjs_location)
	if err != nil {
		return err
	}

	err = download(preloadjs_link, preloadjs_location)
	if err != nil {
		return err
	}

	err = download(packagejson_link, packagejson_location)
	if err != nil {
		return err
	}

	logger.Info("Starting", build+"...")
	cmd = exec.Command(exe_location)
	err = cmd.Start()

	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	logger.Success("Started", build+"!")

	return nil
}

func download(link string, location string) (err error) {
	logger.Infof("Downloading %s...\n", strings.Split(location, "/")[len(strings.Split(location, "/"))-1])
	file, err := os.Create(location)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := http.Get(link)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	n, err := io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	logger.Infof("Downloaded %s (%d bytes.)\n", strings.Split(location, "/")[len(strings.Split(location, "/"))-1], n)
	return nil
}
