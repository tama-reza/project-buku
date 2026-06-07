# project-buku

Cara penggunaan: project-buku memiliki dua subfolder dan satu file main.go yang berfungsi sebagai main entry dari aplikasi. Dua subfolder tersebut berisikan fungsi router dan controller. Pada bukuRouter.go, berisikan router dimana basic auth dibutuhkan. Sementara bukuController memiliki fungsi-fungsi untuk koneksi ke datatabase maupun fungsi-fungsi untuk menjalankan method GET, POST, DELETE yang dideklarasikan di subfolder router.

project-buku/                   <-- Ini folder utama proyek
├── controllers/
│   └── bukuController.go       <-- Berisi fungsi ConnectDB, CreateBuku, dll.
├── routers/
│   └── bukuRouter.go           <-- Berisi fungsi StartServer()
├── .env                        <-- 📄 DI SINI (Berisi password database)
├── .gitignore                  <-- 📄 DI SINI (Berisi teks ".env" agar password aman)
├── go.mod                      <-- File modul Go 
├── go.sum                      <-- File dependensi Go
├── main.go                     <-- File utama untuk menjalankan aplikasi
└── schema.sql                  <-- Berisi fungsi Tabel SQL