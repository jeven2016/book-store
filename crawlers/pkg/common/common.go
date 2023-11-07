package common

import (
	"context"
	"image"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	_ "image/jpeg" // 导入 image/jpeg 包，以便解码 JPG 图片
)

func CreateMongoClient(log *zap.Logger) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://db_user:db_pwd@localhost:27017/books?retryWrites=true&w=majority&authSource=admin&maxPoolSize=20"))

	if err != nil {
		log.Error("Cannot connect to mongodb", zap.Error(err))
		return nil, err
	}

	// 检测MongoDB是否连接成功
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error("Failed to Ping for mongodb", zap.Error(err))
		return nil, err
	}

	return client, nil
}

func ValidJpgImage(jpgPath string) (valid bool, err error) {
	file, err := os.Open(jpgPath)
	if err != nil {
		return
	}
	defer file.Close()
	//image.Decode 函数用于解码图像文件，并返回一个 image.Image 接口类型的对象，代表解码后的图像。这个函数会自动识别图像的格式，并根据格式进行解码。使用 image.Decode 函数可以获取完整的图像数据，可以对图像进行处理、修改和保存等操作。
	//image.DecodeConfig 函数用于获取图像文件的基本信息，而不需要完全解码图像。它返回一个 image.Config 类型的对象，包含图像的宽度、高度、颜色模式等信息，但不包含图像的像素数据。使用 image.DecodeConfig 函数可以快速获取图像的基本信息，而无需完全解码图像。
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return
	}

	if format != "jpeg" {
		return false, err
	}

	return true, err
}

func GetFileExtFromUrl(urlPath string) (string, error) {
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(parsedURL.Path)
	ext := filepath.Ext(filename)
	return ext, nil
}
