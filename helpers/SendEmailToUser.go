package helpers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

func SendEmailToUser(userEmail, username, _ string) error {

	// Mengambil tanggal saat ini
	currentDate := time.Now()
	// Menambahkan 1 hari ke tanggal saat ini
	startDate := currentDate.AddDate(0, 0, 1).Format("02/01/2006")
	// Menambahkan 3 bulan ke "Mulai Tanggal"
	endDate := currentDate.AddDate(0, 3, 0).Format("02/01/2006")

	// Mengambil konfigurasi server SMTP dari variabel lingkungan
	smtpServer := os.Getenv("SMTPSERVER")
	smtpPortStr := os.Getenv("SMTPPORT")
	smtpUsername := os.Getenv("SMTPUSERNAME")
	smtpPassword := os.Getenv("SMTPPASSWORD")

	// Konversi smtpPortStr menjadi int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	// Pesan email dalam format HTML
	message := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Selamat Bergabung Sebagai Mahasiswa Magang</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f3f3f3;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					padding: 20px;
					background-color: #ffffff;
					box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #0073e6;
				}
				p {
					font-size: 16px;
				}
				ul {
					list-style-type: disc;
				}
				ol {
					list-style-type: decimal;
				}
				.footer {
					background-color: #0073e6;
					color: #fff;
					padding: 10px 0;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div "header">
					<h1>Selamat Bergabung Sebagai Mahasiswa di PT Krisnadwipayana</h1>
				</div>
				<p>Halo, ` + username + `</p>

				<p>Kami dengan senang hati ingin memberitahukan bahwa Anda telah diterima sebagai (mahasiswa) magang di PT. Krisnadwipayana. Kami sangat bersemangat untuk menyambut Anda menjadi bagian dari tim kami. Keputusan ini didasari oleh potensi dan kualifikasi yang luar biasa yang Anda tunjukkan selama proses seleksi.</p>

				<p><strong>Detail Kontrak Magang:</strong></p>
				<ul>
					<li>Mulai Tanggal: ` + startDate + `</li>
					<li>Berakhir Tanggal: ` + endDate + `</li>
					<li>Jam Kerja: 08.00 - 16.00 WIB</li>
					<li>Lokasi: JL. Gatot Utomo Kav 5</li>
				</ul>

				<p><strong>Langkah Selanjutnya:</strong></p>
				<ol>
					<li>Konfirmasi Kehadiran: Tolong konfirmasi penerimaan ini dengan membalas email ini atau menghubungi Mutiakhoirunniza di <a href="mailto:mutiakhoirunniza@ac.id">mutiakhoirunniza@ac.id</a> atau 08316281026.</li>
					<li>Dokumen-dokumen: Kami akan mengirimkan Anda berkas persyaratan dan formulir yang perlu Anda isi sebagai persyaratan magang. Mohon lengkapi dan kembalikan dokumen-dokumen ini sesegera mungkin.</li>
				</ol>

				<p><strong>Pengenalan Tim:</strong> Anda akan dikenalkan kepada tim yang akan menjadi mentormu selama magang. Mereka akan membantu Anda beradaptasi dan memberikan bimbingan selama Anda belajar di perusahaan.</p>

				<p>Kami sangat menghargai kesempatan ini untuk bekerja bersama Anda dan membantu Anda dalam pengembangan karier Anda. Jika Anda memiliki pertanyaan atau perlu bantuan lebih lanjut, jangan ragu untuk menghubungi kami kapan saja.</p>

				<p>Kami berharap Anda akan mendapatkan pengalaman yang berharga selama magang di PT. Krisnadwipayana dan berkembang bersama kami. Selamat datang di tim kami!</p>

				<p>Terima kasih dan selamat bersiap-siap untuk memulai petualangan baru Anda.</p>
				<div class="footer">
					<p>Salam, Mutia Khoirunniza</p>
					<p>JL. Gatot Utomo Kav 5</p>
					<p>08316281026</p>
				</div>
			</div>
		</body>
		</html>
	`

	// Konfigurasi pengiriman email menggunakan gomail
	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUsername)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Selamat Bergabung Sebagai Magang")
	m.SetBody("text/html", message)

	// Mengirim email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("Gagal mengirim email: %w", err)
	}

	return nil
}
