### 第二周作业答案 
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，应该 Wrap 这个 error，抛给上层。因为Error需要立即处理，如果不wrap这个error那么上层不知道这个Error是不是没有找到内容还是数据库查询的其他错误。
具体代码如下:
``` Go
// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("name = ? and status = 1", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}
```