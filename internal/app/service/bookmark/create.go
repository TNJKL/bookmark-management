package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	"gorm.io/gorm"
)

// CreateBookmark generates a unique short code and creates a new bookmark for a user
func (b *bookmarkService) CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	//ban đầu e tính đặt Base62 vào struct model Bookmark luôn ,
	//sau đó e  tính gọi hook BeforeCreate để generate sẵn "code_int" rồi encode Base62 nó luôn rồi mới tạo record bookmark
	//--> Nhưng theo nguyên tắc thì model nên là Data Object thuần túy chỉ chứa dữ liệu, ko dc chứa dependency bên ngoài thì phải
	//--> Với lại hình như GORM hook ko có cơ chế nào để inject dependency luôn đúng ko a
	//Nên e tìm hiểu rồi chốt lại phương án cuối cùng đó là INSERT và UPDATE bookmark trong 1 Transaction
	//Cách này thì nếu UPDATE lỗi --> INSERTcũng bị xóa ---> DB không bao giờ có record nào bị thiếu "code" cả :v

	//Cái này được gọi là Trade Off đúng ko anh , mình sẽ tốn thêm 1 lần query và 1 tí hiệu năng ở DB
	//---> Nhưng về mặt lâu dài thì code lưu trong DB sẽ ko bao giờ bị trùng lặp khi hệ thống scale lên hàng triệu bookmarks

	var result *model.Bookmark
	//commit nếu hàm trả về nil
	//rollback nếu hàm trả về lỗi
	err := b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//insert bookmark ( lúc này code = "" ,code_int --> auto_incr)
		bm := &model.Bookmark{
			Description: description,
			URL:         url,
			UserID:      userID,
			Code:        "pending", //Tạm thời mình để "pending" vì field "code" NOT NULL
		}
		created, err := b.repo.CreateBookmarkTx(ctx, tx, bm)
		if err != nil {
			return err //có lỗi --> transaction tự động rollback
		}

		//encode code_int --> base62
		base62code := b.base62.Encode(created.CodeInt)

		//ghep prefix + encode
		code := utils.GetSQLPrefix() + base62code

		//update bookmark : update column "code" trong cung`transaction
		err = b.repo.UpdateCodeTx(ctx, tx, created.ID, code)
		if err != nil {
			return err ///có lỗi --> transaction tự động rollback ( Xóa luôn Insert ở bước 1 )
		}

		created.Code = code
		result = created

		return nil //commit transaction cả 2 bước'
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
