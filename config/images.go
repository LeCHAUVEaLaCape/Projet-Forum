package config

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
type Progress struct {
	TotalSize int64
	BytesRead int64
}

func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	return
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var id string
	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a reference to the fileHeaders
	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = os.MkdirAll("./assets/uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name, _ := uuid.NewV4()
		f, err := os.Create(fmt.Sprintf("./assets/uploads/%s%s", name, filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()
		pr := &Progress{
			TotalSize: fileHeader.Size,
		}
		_, err = io.Copy(f, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		database, err := sql.Open("sqlite3", "./db-sqlite.db")
		CheckError(err)
		defer database.Close()
		rows, err := database.Query("SELECT id FROM pendingPosts ORDER BY id DESC LIMIT 1")
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id)
			CheckError(err)
		}
		var image string
		rows, err = database.Query("SELECT image FROM pendingPosts WHERE id = ?", id)
		CheckError(err)
		for rows.Next() {
			err := rows.Scan(&image)
			CheckError(err)
		}
		image += "," + "../assets/uploads/" + name.String() + filepath.Ext(fileHeader.Filename)

		//range over database
		tx, err := database.Begin()
		stmt, err := tx.Prepare("UPDATE pendingPosts SET image = ? WHERE id = ?")
		CheckError(err)
		_, err = stmt.Exec(image, id)
		CheckError(err)
		tx.Commit()
		rows.Close()
	}
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func Delete_image(state string, id string) {
	var image string
	var each_image []string
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// delete the image
	var rows *sql.Rows
	if state == "pendingPosts" {
		rows, err = database.Query("SELECT image FROM pendingPosts WHERE id = ?", id)
	} else if state == "mainPost" {
		rows, err = database.Query("SELECT image FROM posts WHERE id = ?", id)
	}
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&image)
		CheckError(err)
	}
	image = strings.Replace(image, ",.", "", 1)
	image = strings.Replace(image, "..", ".", -1)
	each_image = strings.Split(image, ",")

	for i := 0; i < len(each_image); i++ {
		err = os.Remove(each_image[i])
		CheckError(err)
	}
}
