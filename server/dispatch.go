package server

//DispatchCommand gfhjkljhgfhjkl;
func (s *Server) DispatchCommand(info []string) (string, error) {

	cmd := info[0]
	data := info[1:]
	var error error

	switch cmd {
	case "SET":
		s.c.Set(data[0],data[1])
	case "GET":
		return s.c.Get(data[0])
	case "DEL":
		s.c.Del(data[0])
	}
	if error != nil {
		return "", error
	}
	return "OK", nil
}
