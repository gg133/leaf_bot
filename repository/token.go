package repository

type Bucket string

//Constants for buckets' names
const (
	AccessTokens  Bucket = "access_tokens"
	RequestTokens Bucket = "request_tokens"
)

//TokenRepository is interface which implements DB Save and Get functions.
type TokenRepository interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
