package storages

import (
	"gitlab.com/andrsolo21/courseproject/internal/lots"
	"gitlab.com/andrsolo21/courseproject/internal/sessions"
	"gitlab.com/andrsolo21/courseproject/internal/users"
	"gitlab.com/andrsolo21/courseproject/pkg/log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"time"
	// Register some standard stuff
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DataB struct {
	DB *gorm.DB
}

type INTT interface {
	CreateTables() DataB
	AddUser(use users.User) error
	AddSession(use sessions.Session) error
	CheckEmail(email string) bool
	GetUserByID(id int) (users.User, error)
	GetSesByToken(token string) (sessions.Session, error)
	GetUserByEmPass(em string, pass string) (el users.User, err error)
	ChangeUser(us users.User, id int) users.User
	CountUsers() int
	AddLot(l lots.Lot) (lots.Lot, error)
	GetLots(typ string) (lots []lots.Lot, d DataB)
	GetUsersLots(id int, role string) (lots []lots.Lot)
	Db() DataB
	GetLotByID(id int) (el lots.Lot, err error)
	SMBBuyIT(userID int, l lots.Lot, price float64) (el lots.Lot, err error)
	UpdateLot(l lots.Lot, id int) lots.Lot
	KillBadLots(logger log.Logger)
	DeleteCrLor(id int)
	CheckDelLot(id int) (bool, lots.Lot)
}

func NewDataB() (d DataB, err error) {
	user := "gorm_user"
	gormDataBase := "gorm_db"
	pas := "gormDB"

	dsn := "postgres://" + user + ":" + pas + "@localhost:5432/" + gormDataBase +
		"?sslmode=disable&fallback_application_name=fintech-app"
	database, err := gorm.Open("postgres", dsn)

	//database.

	/*if err != nil {
		fmt.Printf("can't connect to db: %s", err)
	}*/

	return DataB{DB: database}, err
}

func (db DataB) CreateTables() DataB {

	db.DB = db.DB.AutoMigrate(&users.User{}, &sessions.Session{}, &lots.Lot{})

	//db.DB = db.DB.CreateTable(&userS.ShortUser{})

	return db
}

func (db DataB) Db() DataB {
	return db
}

func (db DataB) AddUser(use users.User) error {
	db.DB = db.DB.Create(&use)

	if err := db.DB.Error; err != nil {
		return err
	}
	return nil
}

func (db DataB) AddSession(use sessions.Session) error {
	use.SessionID = "Bearer " + use.SessionID
	db.DB = db.DB.Create(&use)
	if err := db.DB.Error; err != nil {
		return err
	}
	return nil
}

func (db DataB) CheckEmail(email string) bool {
	var el users.User

	db.DB.Where("email = ?", email).First(&el)

	return el.ID != 0

}

func (db DataB) GetUserByID(id int) (users.User, error) {
	var el users.User

	db.DB.Where("ID = ?", id).First(&el)

	if el.ID == 0 {
		return el, errors.New("this user dosen't exist")
	}

	return el, nil
}

func (db DataB) GetSesByToken(token string) (sessions.Session, error) {
	var el sessions.Session

	db.DB.Where("session_id = ?", token).First(&el)

	if el.UserID == 0 {
		return el, errors.New("this sesion dosen't exist")
	}

	if el.ValidUntil.Before(time.Now()) {
		return el, errors.New("expiration date of token is out")
	}

	return el, nil
}

func (db DataB) GetUserByEmPass(em string, pass string) (el users.User, err error) {

	db.DB.Where("email = ? AND password = ?", em, pass).First(&el)

	if el.ID != 0 {
		return el, nil
	}
	return el, errors.New("this user is not exist")
}

func (db DataB) ChangeUser(us users.User, id int) users.User {

	var el users.User

	db.DB.Where(users.User{ID: id}).Assign(users.User{
		FirstName: us.FirstName,
		LastName:  us.LastName,
		Birthday:  us.Birthday,
		UpdatedAt: time.Now(),
	}).FirstOrCreate(&el)

	return el
}

func (db DataB) CountUsers() int {

	var count int
	var usrs []users.User
	db.DB.Find(&usrs).Count(&count)

	return count
}

func (db DataB) AddLot(l lots.Lot) (lots.Lot, error) {
	db.DB = db.DB.Create(&l)

	if err := db.DB.Error; err != nil {
		return l, err
	}
	return l, nil
}

func (db DataB) GetLots(typ string) (lots []lots.Lot, d DataB) {

	if typ == "" {
		db.DB.Find(&lots)

		return lots, db
	}
	db.DB.Where("status = ?", typ).Find(&lots)

	return lots, db
}

func (db DataB) GetUsersLots(id int, role string) (lots []lots.Lot) {

	switch role {
	case "":
		db.DB.Where("creator_id = ?", id).Or("buyer_id = ?", id).Find(&lots)

	case "own":
		db.DB.Where("creator_id = ?", id).Find(&lots)

	case "buyed":
		db.DB.Where("buyer_id = ?", id).Find(&lots)

	}
	//db.DB.Where("status = ?", typ).Find(&lots)

	return lots
}

func (db DataB) GetLotByID(id int) (el lots.Lot, err error) {
	db.DB.Where("ID = ?", id).First(&el)

	if el.ID == 0 {
		return el, errors.New("this lot dosen't exist")
	}

	return el, err
}

func (db DataB) SMBBuyIT(userID int, l lots.Lot, price float64) (el lots.Lot, err error) {

	l.BuyPrice = price
	l.BuyerID = userID

	db.DB.Where(lots.Lot{ID: l.ID}).Assign(l).FirstOrCreate(&el)

	return el, err
}

func (db DataB) UpdateLot(l lots.Lot, id int) lots.Lot {

	var el lots.Lot

	db.DB.Where(lots.Lot{ID: id}).Assign(lots.Lot{
		Title:       l.Title,
		Description: l.Description,
		MinPrice:    l.MinPrice,
		PriceStep:   l.PriceStep,
		Status:      l.Status,
		EndAt:       l.EndAt,
		UpdatedAt:   l.UpdatedAt,
		DeletedAt:   l.DeletedAt,
	}).FirstOrCreate(&el)

	return el
}

func (db DataB) KillBadLots(logger log.Logger) {
	var lts []lots.Lot

	for {
		db.DB.Where("end_at < now() AND status != ?", "finished").Assign(lots.Lot{
			Status: "finished",
		}).FirstOrCreate(&lts)

		for _, el := range lts {
			logger.Infof("lot %d was finished", el.ID)
		}

		time.Sleep(time.Second)
	}
}

func (db DataB) DeleteCrLor(id int) {
	db.DB.Where("id = ?", id).Delete(lots.Lot{})
}

func (db DataB) CheckDelLot(id int) (bool, lots.Lot) {
	var c int
	var tr lots.Lot
	db.DB.Where("id = ?", id).Count(&c).First(&tr)

	return c == 0, tr
}
