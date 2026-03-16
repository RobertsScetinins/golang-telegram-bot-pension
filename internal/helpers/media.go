package helpers

import (
	"mime"
	"path"

	"github.com/go-telegram/bot/models"
)

type MediaData struct {
	Type   string
	FileId string
}

type MediaHandler func(mediaType string, mediaObj interface{}) (*MediaData, error)

var mediaHandlerMap = map[string]MediaHandler{
	"photo": HandlePhoto,
}

func GetMediaByType(message *models.Message) (string, interface{}) {
	switch {
	case message.Sticker != nil:
		return "sticker", message.Sticker
	case message.Photo != nil:
		return "photo", message.Photo
	case message.Document != nil:
		return "document", message.Document
	case message.Animation != nil:
		return "animation", message.Animation
	case message.Voice != nil:
		return "voice", message.Voice
	case message.Video != nil:
		return "video", message.Video
	case message.VideoNote != nil:
		return "video_note", message.VideoNote
	default:
		return "", nil
	}
}

func HandlePhoto(mediaType string, mediaObj interface{}) (*MediaData, error) {
	photos := mediaObj.([]models.PhotoSize)
	largest := photos[len(photos)-1]

	return &MediaData{
		Type:   mediaType,
		FileId: largest.FileID,
	}, nil
}

func ProcessMedia(update *models.Update) (*MediaData, error) {
	mediaType, mediaObj := GetMediaByType(update.Message)

	if handler, exists := mediaHandlerMap[mediaType]; exists {
		return handler(mediaType, mediaObj)
	}

	return &MediaData{}, nil
}

func GetMimeTypeFromUrl(url string) string {
	ext := path.Ext(url)
	mimeType := mime.TypeByExtension(ext)

	if mimeType == "" {
		return ""
	}

	return mimeType
}
