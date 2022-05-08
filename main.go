package main

// import fyne
import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

//agrigate function for parse()
func findValue(str string, like string) string {
	result := ""
	if strings.Contains(str, like) {
		result = strings.Trim(str, " ")
		result = strings.Replace(result, like, "", -1)
		return result
	} else {
		return "NotFound"
	}
}

//Function for parsing csv file
func parse(path_for_read string, path_for_save string) {
	csvFile, _ := os.Open(path_for_read)
	reader := csv.NewReader(csvFile)
	outputFile, err := os.Create(path_for_save)
	if err != nil {
		log.Fatal(err)
	}
	csvwriter := csv.NewWriter(outputFile)
	//Prepare slice and put title there
	columnTitle := []string{"EMAIL", "PHONE", "NAME", "LASTNAME", "CHAT", "COMPANY", "COUNTRY", "COMMNENT"}
	_ = csvwriter.Write(columnTitle)
	//In cycle split string data from csv row and search for needed fields whoth help of findValue()
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		array := strings.Split(line[0], "\n")
		row := []string{"", "", "", "", "", "", "", ""}
		for _, item := range array {
			name := findValue(item, "* Name:")
			email := findValue(item, "* Email:")
			phone := findValue(item, "* Phone number:")
			sname := findValue(item, "* Last Name:")
			chat := findValue(item, "* Your Skype/WeChat/WhatsApp ID:")
			company := findValue(item, "* Company Name:")
			country := findValue(item, "* Country:")
			comment := findValue(item, "* Comment:")
			if name != "NotFound" {
				row[2] = name
			}
			if email != "NotFound" {
				row[0] = email
			}
			if phone != "NotFound" {
				row[1] = phone
			}
			if sname != "NotFound" {
				row[3] = sname
			}
			if chat != "NotFound" {
				row[4] = chat
			}
			if company != "NotFound" {
				row[5] = company
			}
			if country != "NotFound" {
				row[6] = country
			}
			if comment != "NotFound" {
				row[7] = comment
			}
		}

		if row[0] != "" {
			_ = csvwriter.Write(row)
		}
	}
	csvwriter.Flush()
	csvFile.Close()
	outputFile.Close()
}

func main() {
	// New app
	a := app.New()
	//New title and window
	w := a.NewWindow("Parce csv")
	// resize window
	w.Resize(fyne.NewSize(600, 400))
	//set icon
	icon, _ := fyne.LoadResourceFromPath("ico.png")
	w.SetIcon(icon)
	//info label about output file
	output_info := widget.NewLabel("You dosn't determined file path to save it")
	output_info.Hide() //in deafault this label is hiden
	// New Buttton
	btn := widget.NewButton("Open .csv files", func() {
		// Using dialogs to open files
		// first argument func(fyne.URIReadCloser, error)
		// 2nd is parent window in our case "w"
		// r for reader
		// _ is ignore error
		file_Dialog := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, _ error) {
				// read files
				path_for_read := r.URI().Path()

				data, _ := ioutil.ReadAll(r)
				// reader will read file and store data
				// now result
				//fmt.Println(string(data))
				result := fyne.NewStaticResource("name", data)
				// lets display our data in label or entry
				//entry := widget.NewMultiLineEntry()
				entry := widget.NewEntry()
				entry.Resize(fyne.NewSize(600, 350))
				// string() function convert byte to string
				entry.SetText(string(result.StaticContent))

				entry_save_path := widget.NewEntry()
				entry_save_path.Resize(fyne.NewSize(500, 50))
				entry_save_path.Move(fyne.NewPos(0, 350))
				// Lets show and setup content
				// tile of our new window
				win := fyne.CurrentApp().NewWindow(
					string(result.StaticName)) // title/name
				btn_parse := widget.NewButton("Parse", func() {
					parse(path_for_read, entry_save_path.Text) //run main parsing function
					output_info.SetText("You file is located in: " + entry_save_path.Text)
					output_info.Show() //unhide output label
					win.Hide()         // Hide current windows
				})
				btn_parse.Resize(fyne.NewSize(100, 50))
				btn_parse.Move(fyne.NewPos(500, 350))
				//w.SetContent(container.NewScroll(entry))
				win.SetContent(container.NewWithoutLayout(entry, entry_save_path, btn_parse))
				win.Resize(fyne.NewSize(600, 400))
				// show/display content
				win.Show()
				// we are almost done
			}, w)
		// fiter to open .csv files only
		// array/slice of strings/extensions
		file_Dialog.SetFilter(
			storage.NewExtensionFileFilter([]string{".csv"}))
		file_Dialog.Show()
		// Show file selection dialog.
	})
	// lets show button in parent fmt.Println(line)window
	w.SetContent(container.NewVBox(
		btn,
		output_info,
	))
	w.ShowAndRun()
}
