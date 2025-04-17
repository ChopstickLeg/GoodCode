package handler

import (
	"fmt"
	"io"
	"net/http"

	middleware "github.com/chopstickleg/good-code/api/_middleware"
)

func AddPRHandler(w *http.ResponseWriter, r http.Request) {
	middleware.AllowMethods(http.MethodPost)(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}
		

		file, _, err := r.FormFile("diff")
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusInternalServerError)
			return
		}

		fmt.Println("File content:", string(fileBytes))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
	})
}
