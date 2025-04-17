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

func ParseUser(user any) string {

	if user == nil || len(user.(string)) < 1 {
		return "N/A"
	}

	return user.(string)
}
