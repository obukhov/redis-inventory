package scanner
import (
	"github.com/mediocregopher/radix/v4"
)
type RedisScanner struct {
	client radix.Client
}

func NewScanner(client radix.Client) *RedisScanner {
	return &RedisScanner{
		client: client,
	}
}

func (*RedisScanner) Scan() {

}
