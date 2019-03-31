package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

//Параметры соединения
const (
	ConnHost = ""
	ConnPort = "3333"
	ConnType = "tcp"
)
const maxUploadSize = 300 * 1024 * 1024 // 40 mb
const uploadPath = "D:\\tmp\\"

func main() {
	http.HandleFunc("/upload", uploadFileHandler())
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, 0777)
	}
	copy("D:\\Type1\\App_E_Dog.exe", uploadPath+"App_E_Dog.exe")
	fmt.Println("filepath: " + uploadPath + "App_E_Dog.exe")
	copy("D:\\Type1\\cygwin1.dll", uploadPath+"cygwin1.dll")

	fs := http.FileServer(http.Dir(uploadPath))
	fmt.Println(fs)
	http.Handle("/files/", http.StripPrefix("/files", fs))

	log.Print("Server started on localhost:3333, use /upload for uploading files and /files/{fileName} for downloading")
	log.Fatal(http.ListenAndServe(":3333", nil))
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			status int
			err    error
		)
		defer func() {
			if nil != err {
				fmt.Println(err)

				http.Error(w, err.Error(), status)
			}
		}()
		// parse request
		// const _24K = (1 << 20) * 24
		if err = req.ParseMultipartForm(32 << 20); nil != err {
			status = http.StatusInternalServerError
			return
		}
		fmt.Println("No memory problem")
		for _, fheaders := range req.MultipartForm.File {
			for _, hdr := range fheaders {
				// open uploaded
				var infile multipart.File
				if infile, err = hdr.Open(); nil != err {
					status = http.StatusInternalServerError
					return
				}
				// open destination
				var outfile *os.File
				if outfile, err = os.Create(uploadPath + hdr.Filename); nil != err {
					status = http.StatusInternalServerError
					return
				}
				// 32K buffer copy
				var written int64
				if written, err = io.Copy(outfile, infile); nil != err {
					status = http.StatusInternalServerError
					return
				}

				// w.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
				fmt.Println([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
				convertFile(uploadPath)

				Filename := uploadPath + "e_dog_data.txt"
				Openfile, err := os.Open(Filename)
				if err != nil {
					//File not found, send 404
					http.Error(w, "File not found.", 404)
					return
				}
				defer Openfile.Close() //Close after function return

				// Get the content
				contentType, err := getFileContentType(Openfile)
				if err != nil {
					panic(err)
				}

				fmt.Println("Content Type: " + contentType)

				//File is found, create and send the correct headers

				//Get the Content-Type of the file
				//Create a buffer to store the header of the file in
				// FileHeader := make([]byte, 512)
				// //Copy the headers into the FileHeader buffer
				// Openfile.Read(FileHeader)
				// //Get content type of file
				// FileContentType := http.DetectContentType(FileHeader)

				//Get the file size
				FileStat, _ := Openfile.Stat()                     //Get info from file
				FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string
				fmt.Print(FileSize)

				//Send the headers
				fmt.Println(filepath.Base(Filename))
				w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(Filename))
				w.Header().Set("Content-Type", contentType)
				w.Header().Set("Content-Length", FileSize)

				//Send the file
				//We read 512 bytes from the file already, so we reset the offset back to 0
				Openfile.Seek(0, 0)
				io.Copy(w, Openfile) //'Copy' the file to the client
				// http.ServeContent(w, req, Openfile.Name(), FileStat.ModTime(), Openfile)
			}
		}
	})
}

//Ищем  content-type
func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// func uploadFileHandler() http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		fmt.Println(filepath.Join(uploadPath + "Advocam_speedcam_V1.txt"))

// 		in, _, err := req.FormFile("uploadfile")
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		defer in.Close()
// 		convertFile(uploadPath)
// 		Filename := uploadPath + "e_dog_data.txt"
// 		Openfile, err := os.Open(Filename)
// 		defer Openfile.Close() //Close after function return
// 		if err != nil {
// 			//File not found, send 404
// 			http.Error(w, "File not found.", 404)
// 			return
// 		}

// 		//File is found, create and send the correct headers

// 		//Get the Content-Type of the file
// 		//Create a buffer to store the header of the file in
// 		FileHeader := make([]byte, 512)
// 		//Copy the headers into the FileHeader buffer
// 		Openfile.Read(FileHeader)
// 		//Get content type of file
// 		FileContentType := http.DetectContentType(FileHeader)

// 		//Get the file size
// 		FileStat, _ := Openfile.Stat()                     //Get info from file
// 		FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

// 		//Send the headers
// 		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(Filename))
// 		w.Header().Set("Content-Type", FileContentType)
// 		w.Header().Set("Content-Length", FileSize)

// 		//Send the file
// 		//We read 512 bytes from the file already, so we reset the offset back to 0
// 		Openfile.Seek(0, 0)
// 		io.Copy(w, Openfile) //'Copy' the file to the client
// 	})
// }

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Конвертирует файл с использованием проприетарного по
func convertFile(dirname string) {

	fmt.Println(dirname + "App_E_Dog.exe")
	cmd := exec.Command("powershell", dirname+"App_E_Dog.exe")
	cmd.Dir = dirname
	log.Printf("Running command and waiting for it to finish...")
	error := cmd.Run()
	if error != nil {
		fmt.Println("Error launching:", error.Error())
		log.Printf("Command finished with error: %v", error)
	}
}
