package client

import "sync"

// ClientObjectPool is used to make sure client instances are being used optimized.
type ClientObjectPool struct {
	pool *sync.Pool
}

// NewObjectPool returns a client object pool instance.
func NewObjectPool() *ClientObjectPool {
	return &ClientObjectPool{
		pool: &sync.Pool{
			New: func() any {
				return &Client{}
			},
		},
	}
}

// GetClientObj returns a fresh client object.
func (c *ClientObjectPool) GetClientObj() *Client {
	return c.pool.Get().(*Client)
}

// PutClientObj puts back an unused client object.
func (c *ClientObjectPool) PutClientObj(obj *Client) {
	obj.reset()
	c.pool.Put(obj)
}
