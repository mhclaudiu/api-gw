package handler

func (c *ClientxOBJ) Update(data ClientxOBJ) {

	Mu.Lock()

	*c = data

	defer Mu.Unlock()
}

func (c ClientMap) Add(data ClientxOBJ) {

	Mu.Lock()

	c[data.ClientAddr] = &data

	defer Mu.Unlock()
}

func (c *ClientxOBJ) GetRequests() int {

	return c.ClientRequests
}
