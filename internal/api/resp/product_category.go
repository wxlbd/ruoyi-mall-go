package resp

import "time"

// ProductCategoryResp 商品分类响应
type ProductCategoryResp struct {
	ID          int64     `json:"id"`
	ParentID    int64     `json:"parentId"`
	Name        string    `json:"name"`
	PicURL      string    `json:"picUrl"`
	Sort        int32     `json:"sort"`
	Status      int32     `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createTime"`
}
