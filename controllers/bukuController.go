package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Buku struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Image_url    string    `json:"image_url"`
	Release_year int       `json:"release_year"`
	Price        int       `json:"price"`
	Total_page   int       `json:"total_page"`
	Thickness    string    `json:"thickness"`
	Category_id  int       `json:"category_id"`
	Created_at   time.Time `json:"created_at"`
	Created_by   string    `json:"created_by"`
	Modified_at  time.Time `json:"modified_at"`
	Modified_by  string    `json:"modified_by"`
}

type Kategori struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Created_at  time.Time `json:"created_at"`
	Created_by  string    `json:"created_by"`
	Modified_at time.Time `json:"modified_at"`
	Modified_by string    `json:"modified_by"`
}

type User struct {
	ID          int
	Username    string
	Password    string
	Created_at  time.Time
	Created_by  string
	Modified_at time.Time
	Modified_by string
}

const psqlInfo = "host=postgres port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

var db *sql.DB

func ConnectDB() {
	// Load Env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Peringatan: File .env tidak ditemukan, menggunakan nilai default.")
	}

	psqlInfo := os.Getenv("DATABASE_URL")
	if psqlInfo == "" {
		// Ambil data konfigurasi dari file .env
		host := getEnv("PGHOST", "localhost")
		user := getEnv("PGUSER", "postgres")
		password := getEnv("PGPASSWORD", "123")
		dbname := getEnv("PGDATABASE", "db-go-buku")

		// Konversi teks port di .env menjadi angka integer untuk fmt.Sprintf
		portStr := getEnv("PGPORT", "5432")
		port, _ := strconv.Atoi(portStr)

		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	// Set Database
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// Ping
	if err := db.Ping(); err != nil {
		db.Close()
		panic(err)
	}

	// Tabel kategori
	createTableSQL1 := `
	CREATE TABLE IF NOT EXISTS kategori (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(255) NOT NULL,
		modified_at TIMESTAMP,
		modified_by VARCHAR(255)
	);`

	_, err = db.Exec(createTableSQL1)
	if err != nil {
		panic(fmt.Sprintf("Gagal membuat tabel otomatis: %v", err))
	}
	fmt.Println("Tabel 'kategori' siap digunakan atau sudah ada.")

	// Tabel users
	createTableSQL2 := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(255) NOT NULL,
		modified_at TIMESTAMP,
		modified_by VARCHAR(255)
	);`

	_, err = db.Exec(createTableSQL2)
	if err != nil {
		panic(fmt.Sprintf("Gagal membuat tabel otomatis: %v", err))
	}
	fmt.Println("Tabel 'users' siap digunakan atau sudah ada.")

	// Membuat tabel otomatis jika belum ada di database Railway
	createTableSQL3 := `
	CREATE TABLE IF NOT EXISTS buku (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		image_url VARCHAR(255) NOT NULL,
		release_year INTEGER NOT NULL,
		price INTEGER NOT NULL,
		total_page INTEGER NOT NULL,
		thickness VARCHAR(255) NOT NULL,
		category_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(255) NOT NULL,
		modified_at TIMESTAMP,
		modified_by VARCHAR(255),
		CONSTRAINT fk_buku_kategori FOREIGN KEY (category_id) REFERENCES kategori(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(createTableSQL3)
	if err != nil {
		panic(fmt.Sprintf("Gagal membuat tabel otomatis: %v", err))
	}
	fmt.Println("Tabel 'buku' siap digunakan atau sudah ada.")

}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Detail Fitur Authentification
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uname, pwd, ok := ctx.Request.BasicAuth()

		if !ok {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Username dan password tidak boleh kosong"})
			return
		}

		// query ke db (login gagal jika data username dan pass tidak sesuai)
		var storedUsername string
		var storedPassword string
		query := `SELECT username, password FROM users WHERE username = $1 LIMIT 1`
		err := db.QueryRow(query, uname).Scan(&storedUsername, &storedPassword)

		if err != nil || pwd != storedPassword {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Username atau password and salah"})
			return
		}

		// Success
		ctx.Next()
	}
}

// Detail Fitur Kategori
func GetCategories(ctx *gin.Context) {
	var listKategori []Kategori

	// Tidak perlu binding karena client tidak menambahkan data, hanya meminta data dari database

	// SQL Statement
	sqlStatement := `SELECT id, name, created_at, created_by, modified_at, modified_by FROM kategori ORDER BY id ASC`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Gagal Menampilkan Data Kategori"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var k Kategori
		err := rows.Scan(&k.ID, &k.Name, &k.Created_at, &k.Created_by, &k.Modified_at, &k.Modified_by)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca baris data"})
			return
		}
		listKategori = append(listKategori, k)
	}

	if err = rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memproses data"})
		return
	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{"data": listKategori})
}

func CreateCategories(ctx *gin.Context) {
	var newKategori Kategori

	// Binding Request
	if err := ctx.ShouldBindJSON(&newKategori); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newKategori.Created_by == "" {
		newKategori.Created_by = "admin"
	}
	if newKategori.Modified_by == "" {
		newKategori.Modified_by = "admin"
	}

	currentTime := time.Now()
	newKategori.Created_at = currentTime
	newKategori.Modified_at = currentTime

	// Database Set SQL Statement
	sqlStatement := `INSERT INTO kategori (name, created_at, created_by, modified_at, modified_by)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	// QueryRow
	if err := db.QueryRow(sqlStatement, newKategori.Name, newKategori.Created_at, newKategori.Created_by, newKategori.Modified_at, newKategori.Modified_by).Scan(&newKategori.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori ke database"})
		return
	}

	// Success
	ctx.JSON(http.StatusCreated, gin.H{"kategori": newKategori})
}

func GetEachCategories(ctx *gin.Context) {
	// Mengambil dari URL /:id
	idParam := ctx.Param("id")

	var eachKategori Kategori

	// SQL Statement
	sqlStatement := `SELECT id, name FROM kategori WHERE id = $1`

	// Result
	err := db.QueryRow(sqlStatement, idParam).Scan(&eachKategori.ID, &eachKategori.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari database"})
	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{"data": eachKategori})
}

func DeleteCategories(ctx *gin.Context) {
	// Mengambil dari URL /:id
	idParam := ctx.Param("id")

	// Tidak perlu binding

	sqlStatement := `DELETE FROM kategori WHERE id = $1`

	// Result
	res, err := db.Exec(sqlStatement, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Gagal menghapus data"})
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"err": "Kategori tidak ditemukan"})
		return
	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Kategori dengan id:%s berhasil dihapus", idParam)})
}

func GetEachBookCategories(ctx *gin.Context) {
	// mengambil :id URL
	categoryID := ctx.Param("id")

	var eachKategori Kategori

	// SQL Statement
	sqlStatement1 := `SELECT id, name FROM kategori WHERE id = $1`

	// Result
	err1 := db.QueryRow(sqlStatement1, categoryID).Scan(&eachKategori.ID, &eachKategori.Name)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data kategori"})
		return
	}

	var bukuInCategory []Buku

	sqlStatement2 := `SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM buku WHERE category_id = $1`

	// Result
	res, err2 := db.Query(sqlStatement2, categoryID)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data buku"})
		return
	}

	for res.Next() {
		var b Buku

		err := res.Scan(&b.ID, &b.Title, &b.Description, &b.Image_url, &b.Release_year, &b.Price, &b.Total_page, &b.Thickness,
			&b.Category_id, &b.Created_at, &b.Created_by, &b.Modified_at, &b.Modified_by)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data buku"})
			return
		}
		bukuInCategory = append(bukuInCategory, b)

	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{
		"kategori": eachKategori,
		"buku":     bukuInCategory})
}

// Detail Fitur Buku
func GetBuku(ctx *gin.Context) {
	var listBuku []Buku

	// Tidak perlu binding karena client tidak menambahkan data, hanya meminta data dari database

	// SQL Statement
	sqlStatement := `SELECT id, title, description, image_url, release_year, price, total_page, 
	thickness, category_id, created_at, created_by, modified_at, modified_by FROM buku ORDER BY id ASC`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Menampilkan Data Buku"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var b Buku
		err := rows.Scan(
			&b.ID, &b.Title, &b.Description, &b.Image_url, &b.Release_year,
			&b.Price, &b.Total_page, &b.Thickness, &b.Category_id, &b.Created_at,
			&b.Created_by, &b.Modified_at, &b.Modified_by)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca baris data"})
			return
		}
		listBuku = append(listBuku, b)
	}

	if err = rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memproses data"})
	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{"data": listBuku})

}

func CreateBuku(ctx *gin.Context) {
	var newBuku Buku

	// Binding Request
	if err := ctx.ShouldBindJSON(&newBuku); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newBuku.Created_by == "" {
		newBuku.Created_by = "admin"
	}
	if newBuku.Modified_by == "" {
		newBuku.Modified_by = "admin"
	}

	// Menentukan ketebalan buku
	if newBuku.Total_page <= 100 {
		newBuku.Thickness = "tipis"
	} else {
		newBuku.Thickness = "tebal"
	}

	currentTime := time.Now()
	newBuku.Created_at = currentTime
	newBuku.Modified_at = currentTime

	// Database Set SQL Statement
	sqlStatement := `INSERT INTO buku (title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id`

	// QueryRow
	if err := db.QueryRow(
		sqlStatement, newBuku.Title, newBuku.Description, newBuku.Image_url,
		newBuku.Release_year, newBuku.Price, newBuku.Total_page,
		newBuku.Thickness, newBuku.Category_id,
		newBuku.Created_at, newBuku.Created_by,
		newBuku.Modified_at, newBuku.Modified_by).Scan(&newBuku.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data buku ke database"})
		return
	}

	// Success
	ctx.JSON(http.StatusCreated, gin.H{"buku": newBuku})

}

func GetEachBuku(ctx *gin.Context) {
	// mengambil URL /:id
	idParam := ctx.Param("id")

	var eachBuku Buku

	// SQL Statement
	sqlStatement := `SELECT id, title FROM buku WHERE id = $1`

	// Result
	err := db.QueryRow(sqlStatement, idParam).Scan(
		&eachBuku.ID, &eachBuku.Title, &eachBuku.Description, &eachBuku.Image_url,
		&eachBuku.Release_year, &eachBuku.Price, &eachBuku.Total_page,
		&eachBuku.Thickness, &eachBuku.Category_id,
		&eachBuku.Created_at, &eachBuku.Created_by,
		&eachBuku.Modified_at, &eachBuku.Modified_by)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Query Gagal"})
		return
	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Buku yang memiliki id: %s berhasil ditampilkan", idParam),
		"data": eachBuku})
}

func DeleteBuku(ctx *gin.Context) {
	// mengambil URL /:id
	idParam := ctx.Param("id")

	// Tidak perlu binding

	sqlStatement := `DELETE FROM buku WHERE id = $1`

	// Result
	res, err := db.Exec(sqlStatement, idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Gagal menghapus data"})
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"err": "Data buku tidak ditemukan"})
		return
	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Buku dengan id:%v berhasil dihapus",
		idParam)})

}
