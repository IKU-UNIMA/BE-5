package env

import "os"

func GetPublikasiFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_PUBLIKASI_FOLDER_ID")
}

func GetPatenFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_PATEN_FOLDER_ID")
}

func GetPengabdianFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_PENGABDIAN_FOLDER_ID")
}
