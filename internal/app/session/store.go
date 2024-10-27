package session

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/JDinABox/yapa/internal/ilog"
	"github.com/JDinABox/yapa/internal/sqlc/db"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/blake2b"
)

type Store struct {
	db      *db.Queries
	key     []byte
	Options *sessions.Options
}

func NewStore(queries *db.Queries, key []byte, options *sessions.Options) *Store {
	return &Store{
		db:      queries,
		key:     key,
		Options: options,
	}
}

// Get should return a cached session.
func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = &sessions.Options{
		Path:     s.Options.Path,
		Domain:   s.Options.Domain,
		MaxAge:   s.Options.MaxAge,
		Secure:   s.Options.Secure,
		HttpOnly: s.Options.HttpOnly,
	}
	session.IsNew = true

	c, err := r.Cookie(name)
	// no cookie found
	if err != nil {
		return session, nil
	}

	encodedValue := strings.TrimSpace(c.Value)
	// empty cookie value
	if encodedValue == "" {
		return session, nil
	}
	// unencode cookie value
	encryptedValue, err := base64.URLEncoding.DecodeString(encodedValue)
	// invalid base64
	if err != nil {
		return session, nil
	}

	// decrypt cookie value
	idByte, err := s.decrypt(encryptedValue)
	// unable to decrypt cookie value
	if err != nil {
		return session, nil
	}

	id, err := uuid.FromBytes(idByte)
	// invalid uuid
	if err != nil {
		// should not happen
		ilog.Error(err)
		return session, nil
	}
	session.ID = id.String()

	err = s.load(r.Context(), session)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) &&
			!uuid.IsInvalidLengthError(err) &&
			!(err.Error() == "invalid UUID format") &&
			!strings.HasPrefix(err.Error(), "invalid urn prefix") {
			ilog.Error(err)
		}

		return session, nil
	}

	session.IsNew = false

	return session, nil
}

// Save should persist session to the underlying store implementation.
func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.ID == "" {
		if err := s.insert(r.Context(), session); err != nil {
			return err
		}
	} else if err := s.update(r.Context(), session); err != nil {
		return err
	}

	id, err := uuid.Parse(session.ID)
	if err != nil {
		// should not happen
		ilog.Error(err)
		return err
	}

	cookieValue, err := s.encrypt(id[:])
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), base64.URLEncoding.EncodeToString(cookieValue), session.Options))
	return nil
}

// Delete should delete session from the underlying store implementation.
func (s *Store) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	options := *session.Options
	options.MaxAge = -1
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", &options))
	clear(session.Values)
	return s.delete(r.Context(), session)
}

// insert inserts session into database
func (s *Store) insert(ctx context.Context, session *sessions.Session) error {
	sessionID := uuid.New().String()
	sM := db.Session{}
	var err error

	sM.ID, err = sIDToDBKey(sessionID)
	if err != nil {
		return err
	}

	now := time.Now()
	sM.Created = now.Unix()
	sM.Updated = now.Unix()
	sM.Expires = now.Unix() + int64(session.Options.MaxAge)

	sM.Data, err = s.encode(session.Values)
	if err != nil {
		return err
	}

	sM, err = s.db.CreateSession(ctx, db.CreateSessionParams(sM))
	if err != nil {
		return err
	}

	session.ID = sessionID
	saveTimes(&sM, session)
	return nil
}

func (s *Store) update(ctx context.Context, session *sessions.Session) error {
	if session.IsNew {
		return s.insert(ctx, session)
	}
	key, err := sIDToDBKey(session.ID)
	if err != nil {
		return err
	}
	now := time.Now()

	expires, ok := session.Values["expires_on"].(time.Time)
	if !ok || expires.Sub(now.Add(time.Second*time.Duration(session.Options.MaxAge))) < 0 {
		expires = now.Add(time.Second * time.Duration(session.Options.MaxAge))
	}

	data, err := s.encode(session.Values)
	if err != nil {
		return err
	}

	unix := now.Unix()
	err = s.db.UpdateSessionData(ctx, db.UpdateSessionDataParams{
		ID:      key,
		Updated: unix,
		Expires: unix,
		Data:    data,
	})
	if err != nil {
		return err
	}

	// update session times
	session.Values["expires_on"] = expires
	session.Values["updated_on"] = now

	return nil
}

// load loads session from database
func (s *Store) load(ctx context.Context, session *sessions.Session) error {
	key, err := sIDToDBKey(session.ID)
	if err != nil {
		return err
	}

	dbSession, err := s.db.GetSession(ctx, key)
	if err != nil {
		return err
	}

	if int64(dbSession.Expires) < time.Now().Unix() {
		return sql.ErrNoRows
	}
	return s.toSession(&dbSession, session)
}

func (s *Store) delete(ctx context.Context, session *sessions.Session) error {
	key, err := sIDToDBKey(session.ID)
	if err != nil {
		return err
	}

	return s.db.DeleteSession(ctx, key)
}

func (s *Store) toSession(dbS *db.Session, session *sessions.Session) error {
	err := s.decode(dbS.Data, &session.Values)
	if err != nil {
		return err
	}

	saveTimes(dbS, session)
	return nil
}

// saveTimes saves session times to session.Values
func saveTimes(dbS *db.Session, session *sessions.Session) {
	session.Values["created_on"] = time.Unix(dbS.Created, 0)
	session.Values["updated_on"] = time.Unix(dbS.Created, 0)
	session.Values["expires_on"] = time.Unix(dbS.Created, 0)
}

// sIDToDBKey converts uuid.UUID to a blake2b hash
func sIDToDBKey(idStr string) ([]byte, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	key := blake2b.Sum256(id[:])
	return key[:], nil
}

// encode encodes values using msgpack and encrypts them using chacha20poly1305
// does not encode created_on, updated_on, expires_on
func (s *Store) encode(values map[any]any) ([]byte, error) {
	valuesCopy := make(map[any]any, len(values))

	for id, v := range values {
		if id == "created_on" || id == "updated_on" || id == "expires_on" {
			continue
		}
		valuesCopy[id] = v
	}

	b, err := msgpack.Marshal(valuesCopy)
	if err != nil {
		return nil, err
	}

	return s.encrypt(b)
}

var errBadCipherData = errors.New("bad cipher data")

// decode decrypts values using chacha20poly1305 and decodes them using msgpack
func (s *Store) decode(ciphertext []byte, values *map[any]any) error {
	plainData, err := s.decrypt(ciphertext)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(plainData, values)
}

func (s *Store) encrypt(plaintext []byte) ([]byte, error) {
	c, err := chacha20poly1305.NewX(s.key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, c.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return c.Seal(nonce, nonce, plaintext, nil), nil
}

func (s *Store) decrypt(ciphertext []byte) ([]byte, error) {
	c, err := chacha20poly1305.NewX(s.key)
	if err != nil {
		return nil, err
	}
	nonceSize := c.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errBadCipherData
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return c.Open(ciphertext[:0], nonce, ciphertext, nil)
}
