package mongo

import (
	"strings"
	"time"

	"github.com/Machiel/slugify"
	"github.com/asaskevich/govalidator"
	"github.com/blankrobot/pulpe"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const userCol = "users"

// User representation stored in MongoDB.
type User struct {
	ID        bson.ObjectId `bson:"_id"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	FullName  string        `bson:"fullName"`
	Login     string        `bson:"login"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
}

// ToMongoUser creates a mongo user from a pulpe user.
func ToMongoUser(p *pulpe.User) *User {
	id := bson.NewObjectId()
	p.ID = id.Hex()
	p.CreatedAt = id.Time()

	return &User{
		ID:        id,
		UpdatedAt: p.UpdatedAt,
		FullName:  p.FullName,
		Login:     p.Login,
		Email:     p.Email,
	}
}

// FromMongoUser creates a pulpe user from a mongo user.
func FromMongoUser(u *User) *pulpe.User {
	p := pulpe.User{
		ID:        u.ID.Hex(),
		CreatedAt: u.ID.Time(),
		FullName:  u.FullName,
		Login:     u.Login,
		Email:     u.Email,
	}

	if u.UpdatedAt != nil {
		t := (*u.UpdatedAt).UTC()
		p.UpdatedAt = &t
	}

	return &p
}

// UserService represents a service for managing users.
type UserService struct {
	session *Session
}

func (s *UserService) ensureIndexes() error {
	col := s.session.db.C(userCol)

	// Unique login
	index := mgo.Index{
		Key:    []string{"login"},
		Unique: true,
		Sparse: true,
	}

	err := col.EnsureIndex(index)
	if err != nil {
		return err
	}

	// Unique email
	index = mgo.Index{
		Key:    []string{"email"},
		Unique: true,
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateUser creates a new User.
func (s *UserService) CreateUser(uc *pulpe.UserCreation) (*pulpe.User, error) {
	var err error
	col := s.session.db.C(userCol)

	user := pulpe.User{
		FullName: uc.FullName,
		Login:    strings.Replace(slugify.Slugify(uc.FullName), "-", "", -1),
		Email:    uc.Email,
	}

	u := ToMongoUser(&user)

	passwd, err := bcrypt.GenerateFromPassword([]byte(uc.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.Password = string(passwd)

	user.Login, err = resolveSlugAndDo(col, "login", u.Login, "", func(login string) error {
		u.Login = login
		return col.Insert(u)
	})

	if err != nil && mgo.IsDup(err) && strings.Contains(err.Error(), "email") {
		return nil, pulpe.ErrUserEmailConflict
	}

	return &user, err
}

// User returns a User by login or ID.
func (s *UserService) User(selector string) (*pulpe.User, error) {
	var u User

	err := s.session.db.C(userCol).Find(getSelector("login", selector)).One(&u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrUserNotFound
		}

		return nil, err
	}

	return FromMongoUser(&u), nil
}

// Authenticate user with login or email and password.
func (s *UserService) Authenticate(loginOrEmail, passwd string) (*pulpe.User, error) {
	var u User
	var field string

	if govalidator.IsEmail(loginOrEmail) {
		field = "email"
	} else {
		field = "login"
	}

	err := s.session.db.C(userCol).Find(bson.M{field: loginOrEmail}).One(&u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrAuthenticationFailed
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd))
	if err != nil {
		return nil, pulpe.ErrAuthenticationFailed
	}

	return FromMongoUser(&u), nil
}
