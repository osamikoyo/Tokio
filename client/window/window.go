package window

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

func Window() {
	// Инициализация GTK
	gtk.Init(nil)

	// Создание нового окна
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatalf("Unable to create window: %v", err)
	}
	window.SetTitle("Hello GTK")
	window.SetDefaultSize(300, 200)

	// Обработчик события закрытия окна
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Создание метки
	label, err := gtk.LabelNew("Hello, GTK!")
	if err != nil {
		log.Fatalf("Unable to create label: %v", err)
	}

	// Создание кнопки
	button, err := gtk.ButtonNewWithLabel("Quit")
	if err != nil {
		log.Fatalf("Unable to create button: %v", err)
	}
	button.Connect("clicked", func() {
		os.Exit(0)
	})

	// Создание вертикального контейнера и добавление метки и кнопки
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatalf("Unable to create box: %v", err)
	}
	vbox.PackStart(label, true, true, 0)
	vbox.PackStart(button, true, true, 0)

	// Добавление контейнера в окно
	window.Add(vbox)

	window.ShowAll()
	gtk.Main()
}