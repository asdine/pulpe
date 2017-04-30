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

// Ensure UserService implements pulpe.UserService.
var _ pulpe.UserService = new(UserService)

// user representation stored in MongoDB.
type user struct {
	ID        bson.ObjectId `bson:"_id"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	FullName  string        `bson:"fullName"`
	Login     string        `bson:"login"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
}

func (u *user) toPulpeUser() *pulpe.User {
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

// Register creates a new User.
func (s *UserService) Register(uc *pulpe.UserRegistration) (*pulpe.User, error) {
	var err error
	col := s.session.db.C(userCol)

	u := user{
		ID:       bson.NewObjectId(),
		FullName: uc.FullName,
		Login:    strings.Replace(slugify.Slugify(uc.FullName), "-", "", -1),
		Email:    uc.Email,
	}

	passwd, err := bcrypt.GenerateFromPassword([]byte(uc.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.Password = string(passwd)

	u.Login, err = resolveSlugAndDo(col, "", "login", u.Login, "", func(login string) error {
		u.Login = login
		return col.Insert(u)
	})

	if err != nil && mgo.IsDup(err) && strings.Contains(err.Error(), "email") {
		return nil, pulpe.ErrUserEmailConflict
	}

	return u.toPulpeUser(), nil
}

// User returns a User by ID.
func (s *UserService) User(id string) (*pulpe.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrUserNotFound
	}

	return s.userBy(bson.M{
		"_id": bson.ObjectIdHex(id),
	})
}

func (s *UserService) userBy(query bson.M) (*pulpe.User, error) {
	var u user

	err := s.session.db.C(userCol).Find(query).One(&u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrUserNotFound
		}

		return nil, err
	}

	return u.toPulpeUser(), nil
}

// MatchPassword checks is the login or email and password are correct.
func (s *UserService) MatchPassword(loginOrEmail, passwd string) (string, error) {
	var u user
	var field string

	if govalidator.IsEmail(loginOrEmail) {
		field = "email"
	} else {
		field = "login"
	}

	query := s.session.db.C(userCol).Find(bson.M{field: loginOrEmail})
	query = query.Select(bson.M{"_id": 1, "password": 1})
	err := query.One(&u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return "", pulpe.ErrUserAuthenticationFailed
		}

		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd))
	if err != nil {
		return "", pulpe.ErrUserAuthenticationFailed
	}

	return u.ID.Hex(), nil
}

// userSession is stored and represents a logged in user.
type userSession struct {
	ID        string    `bson:"_id"`
	UpdatedAt time.Time `bson:"updatedAt"`
	UserID    string    `bson:"userID"`
}

// UserSessionService represents a service for managing user sessions.
type UserSessionService struct {
	session *Session
}

// Ensure UserSessionService implements pulpe.UserSessionService.
var _ pulpe.UserSessionService = new(UserSessionService)

func (s *UserSessionService) ensureIndexes() error {
	col := s.session.db.C(userSessionCol)

	// Sessions expiration
	index := mgo.Index{
		Key:         []string{"updatedAt"},
		Sparse:      true,
		ExpireAfter: userSessionTTL,
	}

	return col.EnsureIndex(index)
}

// CreateSession store a user session in the database.
func (s *UserSessionService) CreateSession(user *pulpe.User) (*pulpe.UserSession, error) {
	sid, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}

	session := userSession{
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
func (s *UserSessionService) GetSession(id string) (*pulpe.UserSession, error) {
	var us userSession

	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"updatedAt": s.session.now,
		}},
		ReturnNew: true,
	}

	_, err := s.session.db.C(userSessionCol).FindId(id).Apply(change, &us)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrUserSessionUnknownID
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

// Login user with login or email and password.
func (s *UserSessionService) Login(loginOrEmail, password string) (*pulpe.UserSession, error) {
	id, err := s.session.UserService().MatchPassword(loginOrEmail, password)
	if err != nil {
		return nil, err
	}

	return s.CreateSession(&pulpe.User{ID: id})
}

// DeleteSession removes the given session.
func (s *UserSessionService) DeleteSession(id string) error {
	err := s.session.db.C(userSessionCol).RemoveId(id)
	if err == mgo.ErrNotFound {
		return pulpe.ErrUserSessionUnknownID
	}
	return err
}

// Ensure Authenticator implements pulpe.Authenticator.
var _ pulpe.Authenticator = new(Authenticator)

// Authenticator is a service for user authentication.
type Authenticator struct{}

// Authenticate returns the current authenticate user.
func (a Authenticator) Authenticate(session pulpe.Session, token string) (*pulpe.User, error) {
	us, err := session.UserSessionService().GetSession(token)
	if err != nil {
		return nil, err
	}

	return session.UserService().User(us.UserID)
}
