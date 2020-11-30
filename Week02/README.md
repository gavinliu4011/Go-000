### 问题

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

- model层伪代码

  ```golang
  type User struct {
  	ID   int    `sql:"primary_key"`
  	Name string `sql:"type:varchar(36);not null;"`
  }
  ```

- Dao层伪代码

  ```golang
  type Dao struct {
  	db *sql.DB
  }
  
  func NewDao(db *sql.DB) *Dao {
  	return &Dao{
  		db: db,
  	}
  }
  
  func (d *Dao) GetUserByID(id int) (*model.User, error) {
  	user := &model.User{}
  	err := d.db.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(user)
  	if errors.Is(err, sql.ErrNoRows) {
  		// 这里看自己团队设计而定，在这里我先返回pkg/errors
  		err = ErrDaoRecordNotFound
  	}
  	return user, pkgerr.Wrapf(err, "query user %d failed", id)
  }
  ```

- Service层伪代码

  ```golang
  type Service struct {
  	dao *dao.Dao
  }
  
  func NewService(d *dao.Dao) *Service {
  	return &Service{
  		dao: d,
  	}
  }
  
  func (s *Service) GetUserByID(id int) (*model.User, error) {
  	user, err := s.dao.GetUserByID(id)
  	if errors.Is(pkgerr.Cause(err), dao.ErrDaoRecordNotFound) {
  		// 这里可以处理没有查询到记录时的业务逻辑
  		......
  		return ......
  	}
  	return user, err
  }
  ```

  

