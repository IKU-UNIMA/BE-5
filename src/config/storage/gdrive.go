package storage

import (
	"be-5/src/util"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var gDriveSrv *drive.Service

func InitGDrive() {
	ctx := context.Background()

	// get Google Drive service account
	client := serviceAccount("credentials.json")

	var err error
	gDriveSrv, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(err.Error())
	}
}

func serviceAccount(secretFile string) *http.Client {
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal("error while reading the credential file", err)
	}
	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &s)
	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(context.Background())
	return client
}

func CreateFile(fileHeader *multipart.FileHeader, parentId string) (*drive.File, error) {
	if err := util.CheckFileType(fileHeader); err != nil {
		return nil, err
	}

	f := &drive.File{
		MimeType: fileHeader.Header["Content-Type"][0],
		Name:     fileHeader.Filename,
		Parents:  []string{parentId},
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	file, err := gDriveSrv.Files.Create(f).Media(src).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func DeleteFile(fileId string) error {
	err := gDriveSrv.Files.Delete(fileId).Do()
	return err
}
