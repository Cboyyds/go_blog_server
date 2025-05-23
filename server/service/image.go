package service

import (
	"mime/multipart"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"server/utils/upload"
	
	"gorm.io/gorm"
)

type ImageService struct {
}

// 本项目采用local存图片，这里即把图片存在本地也存在了数据库里面
func (imageService *ImageService) ImageUpload(file *multipart.FileHeader) (string, error) {
	oss := upload.NewOss()
	url, filename, err := oss.UploadImage(file) // local和host都实现该方法，都实现了oss接口
	if err != nil {
		return "", err
	}
	
	return url, global.DB.Create(&database.Image{
		Name:     filename,
		URL:      url,
		Category: appTypes.Null,
		Storage:  global.Config.System.Storage(),
	}).Error
}

func (imageService *ImageService) ImageDelete(req request.ImageDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	
	var images []database.Image
	if err := global.DB.Find(&images, req.IDs).Error; err != nil {
		return err
	}
	for _, image := range images {
		if err := global.DB.Transaction(func(tx *gorm.DB) error { // 使用   global.DB.Transaction   开始一个新的数据库事务。它接受一个函数作为参数，该函数在事务的上下文中执行一系列数据库操作。如果函数返回错误，事务将被回滚。
			oss := upload.NewOssWithStorage(image.Storage)
			if err := global.DB.Delete(&image).Error; err != nil {
				return err
			}
			return oss.DeleteImage(image.Name)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (imageService *ImageService) ImageList(info request.ImageList) (list interface{}, total int64, err error) {
	db := global.DB
	
	if info.Name != nil {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	
	if info.Category != nil {
		category := appTypes.ToCategory(*info.Category)
		db = db.Where("category = ?", category)
	}
	
	if info.Storage != nil {
		storage := appTypes.ToStorage(*info.Storage)
		db = db.Where("storage = ?", storage)
	}
	
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	return utils.MySQLPagination(&database.Image{}, option)
}
