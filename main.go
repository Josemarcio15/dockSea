package main

import (
	"embed"
	"log"
	"os"
	"runtime"

	"go-walis/internal/builder"
	"go-walis/internal/config"
	"go-walis/internal/containers"
	"go-walis/internal/core/db"
	"go-walis/internal/dashboard"
	"go-walis/internal/extras"
	"go-walis/internal/images"
	"go-walis/internal/locale"
	"go-walis/internal/networks"
	"go-walis/internal/profiles"
	"go-walis/internal/servers"
	"go-walis/internal/stacks"
	"go-walis/internal/theme"
	"go-walis/internal/updater"
	"go-walis/internal/volumes"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed build/appicon.png
var appIcon []byte

//go:embed VERSION
var appVersion string

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Workaround: o compositor Wayland do LordOS não repassa o tamanho da janela
	// corretamente para o WebKitGTK (WebView fica cortado). Força XWayland nele.
	if runtime.GOOS == "linux" && os.Getenv("XDG_CURRENT_DESKTOP") == "LordOS" {
		os.Setenv("GDK_BACKEND", "x11")
		os.Unsetenv("WAYLAND_DISPLAY")
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Erro ao inicializar banco SQLite: %v", err)
	}
	defer database.Close()

	configService := config.NewConfigService(database)
	builderService := builder.NewService(database)
	dashboardService := dashboard.NewService(database)
	containerService := containers.NewContainerService(database)
	extraService := extras.NewExtraService(database)
	imageService := images.NewImageService(database)
	networkService := networks.NewNetworkService(database)
	profileService := profiles.NewService(database)
	serverService := servers.NewService(database)
	stackService := stacks.NewStackService(database)
	themeService := theme.NewThemeService(database.GetThemesDir())
	localeService := locale.NewLocaleService(database.GetLocalesDir())
	volumeService := volumes.NewVolumeService(database)
	updaterService := updater.NewUpdaterService(appVersion)

	app := application.New(application.Options{
		Name:        "DockSea",
		Description: "DockSea Container Manager",
		Icon:        appIcon,
		Linux: application.LinuxOptions{
			ProgramName: "DockSea",
		},
		Services: []application.Service{
			application.NewService(builderService),
			application.NewService(dashboardService),
			application.NewService(configService),
			application.NewService(containerService),
			application.NewService(extraService),
			application.NewService(imageService),
			application.NewService(localeService),
			application.NewService(networkService),
			application.NewService(profileService),
			application.NewService(serverService),
			application.NewService(stackService),
			application.NewService(themeService),
			application.NewService(updaterService),
			application.NewService(volumeService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "DockSea",
		Width:  1200,
		Height: 750,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(15, 23, 42),
		URL:              "/",
	})

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
