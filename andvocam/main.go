package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//Параметры соединения
const (
	ConnHost = ""
	ConnPort = "3333"
	ConnType = "tcp"
)
const maxUploadSize = 300 * 1024 * 1024 // 40 mb
var uploadPath = "C:\\tmp\\"

const tmp = "\\tmp\\"

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	uploadPath = dir

	http.HandleFunc("/upload", uploadFileHandler())
	if _, err := os.Stat(uploadPath + tmp); os.IsNotExist(err) {
		os.Mkdir(uploadPath+tmp, 0777)
	}
	// copy(uploadPath+"\\Type1\\App_E_Dog.exe", uploadPath+tmp+"App_E_Dog.exe")
	// copy(uploadPath+"\\Type1\\cygwin1.dll", uploadPath+tmp+"cygwin1.dll")

	fs := http.FileServer(http.Dir(uploadPath + tmp))
	fmt.Println(fs)
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
		// req.ParseForm()                     // Parses the request body
		modelType := req.Form.Get("model_type") // x will be "" if parameter is not set
		fmt.Println("Model type: " + modelType)
		workDir := uploadPath + tmp + fmt.Sprintf("%v", makeTimestamp()) + "\\"
		fmt.Println("Workdir" + workDir)
		if _, err := os.Stat(workDir); os.IsNotExist(err) {
			os.Mkdir(workDir, 0777)
		}

		sendEmail(req.Form.Get("emailaddress"), req.Form.Get("comment"))
		uploadFileName := ""
		responseFileName := ""
		if modelType == "Type1" {
			copy(uploadPath+"\\Type1\\App_E_Dog.exe", workDir+"App_E_Dog.exe")
			copy(uploadPath+"\\Type1\\cygwin1.dll", workDir+"cygwin1.dll")
			uploadFileName = "Advocam_speedcam_V1.txt"
			responseFileName = "e_dog_data.txt"
		}
		if modelType == "Type2" {
			copy(uploadPath+"\\Type2\\App_E_Dog.exe", workDir+"App_E_Dog.exe")
			copy(uploadPath+"\\Type2\\cygwin1.dll", workDir+"cygwin1.dll")
			uploadFileName = "speedcam22.txt"
			responseFileName = "DATA_T.BIN"
		}

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
				// if outfile, err = os.Create(workDir + hdr.Filename); nil != err {
				if outfile, err = os.Create(workDir + uploadFileName); nil != err {
					status = http.StatusInternalServerError
					return
				}
				// 32K buffer copy
				var written int64
				if written, err = io.Copy(outfile, infile); nil != err {
					status = http.StatusInternalServerError
					return
				}

				fmt.Println([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
				convertFile(workDir)

				Filename := workDir + responseFileName
				Openfile, err := os.Open(Filename)
				if err != nil {
					//File not found, send 404
					http.Error(w, "File not found.", 404)
					return
				}
				defer Openfile.Close()        //Close after function return
				defer removeContents(workDir) //Clean folder after return

				// Get the content
				contentType, err := getFileContentType(Openfile)
				if err != nil {
					panic(err)
				}

				fmt.Println("Content Type: " + contentType)
				//File is found, create and send the correct headers
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

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
func removeContents(dir string) error {
	fmt.Println("Cleaning " + dir + " dir")
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
