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

const (
	userCol        = "users"
	userSessionCol = "userSessions"
	userSessionTTL = 24 * time.Hour
)

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
		CreatedAt: u.ID.Time().UTC(),
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

// UserSession is stored and represents a logged in user.
type UserSession struct {
	ID        string    `bson:"_id"`
	UpdatedAt time.Time `bson:"updatedAt"`
	UserID    string    `bson:"userID"`
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

	err = col.EnsureIndex(index)
	if err != nil {
		return err
	}

	// Sessions
	col = s.session.db.C(userSessionCol)

	// Sessions expiration
	index = mgo.Index{
		Key:         []string{"updatedAt"},
		Sparse:      true,
		ExpireAfter: userSessionTTL,
	}

	return col.EnsureIndex(index)
}

// CreateUser creates a new User.
func (s *UserService) CreateUser(uc *pulpe.UserRegistration) (*pulpe.User, error) {
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
			return nil, pulpe.ErrUserAuthenticationFailed
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd))
	if err != nil {
		return nil, pulpe.ErrUserAuthenticationFailed
	}

	return FromMongoUser(&u), nil
}

// CreateSession store a user session in the database.
func (s *UserService) CreateSession(user *pulpe.User) (*pulpe.UserSession, error) {
	sid, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}

	session := UserSession{
		ID:        sid,
		UserID:    user.ID,
		UpdatedAt: s.session.now,
	}

	err = s.session.db.C(userSessionCol).Insert(&session)
	if err != nil {
		return nil, err
	}

	return &pulpe.UserSession{
		ID:        sid,
		UserID:    user.ID,
		UpdatedAt: s.session.now,
		ExpiresAt: s.session.now.Add(userSessionTTL),
	}, nil
}

// GetSession gets a session and resets the session expiration date.
func (s *UserService) GetSession(sid string) (*pulpe.UserSession, error) {
	var us UserSession

	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"updatedAt": s.session.now,
		}},
		ReturnNew: true,
	}

	_, err := s.session.db.C(userSessionCol).FindId(sid).Apply(change, &us)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrUserSessionUnknownSid
		}
		return nil, err
	}

	return &pulpe.UserSession{
		ID:        us.ID,
		UserID:    us.UserID,
		UpdatedAt: us.UpdatedAt.UTC(),
		ExpiresAt: us.UpdatedAt.Add(userSessionTTL).UTC(),
	}, nil
}
