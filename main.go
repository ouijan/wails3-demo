package main

import (
	"embed"
	"log"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "wails3-demo",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
			// ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new mainWindow with the necessary options.
	// 'Title' is the title of the mainWindow.
	// 'Mac' options tailor the mainWindow when running on macOS.
	// 'BackgroundColour' is the background colour of the mainWindow.
	// 'URL' is the URL that will be loaded into the webview.
	mainWindow := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:         400,
		Height:        600,
		Frameless:     true,
		AlwaysOnTop:   true,
		DisableResize: true,
		Hidden:        true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},

		// Mac: application.MacWindow{

		// InvisibleTitleBarHeight: 50,
		// Backdrop:                application.MacBackdropTranslucent,
		// TitleBar: application.MacTitleBarHiddenInset,
		// },
		// BackgroundColour: application.NewRGB(27, 38, 54),
		// URL:              "/",
	})

	systemTray := app.NewSystemTray()

	// Support for template icons on macOS
	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	} else {
		// Support for light/dark mode icons
		systemTray.SetDarkModeIcon(icons.SystrayDark)
		systemTray.SetIcon(icons.SystrayLight)
	}

	// Support for menu
	myMenu := app.NewMenu()
	myMenu.Add("Show / Hide").OnClick(func(_ *application.Context) {
		if mainWindow.IsVisible() {
			mainWindow.Hide()
		} else {
			mainWindow.Show()
			mainWindow.Focus()
		}
	})

	myMenu.AddSeparator()
	myMenu.Add("Quit").OnClick(func(_ *application.Context) {
		app.Quit()
	})
	systemTray.SetMenu(myMenu)

	// This will center the window to the systray icon with a 5px offset
	// It will automatically be shown when the systray icon is clicked
	// and hidden when the window loses focus
	systemTray.AttachWindow(mainWindow).WindowOffset(5)

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
