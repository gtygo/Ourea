package server

//DispatchCommand gfhjkljhgfhjkl;
func (s *Server) DispatchCommand(info []string) (string, error) {

	cmd := info[0]
	data := info[1:]
	var error error

	switch cmd {
	case "SET":
		error = s.Db.Set([]byte(data[0]), []byte(data[1]))
	case "GET":
		v, err := s.Db.Get([]byte(data[0]))
		if err != nil {
			return "", err
		}
		return string(v), nil
	case "DEL":
		error = s.Db.Delete([]byte(data[0]))

	}
	if error != nil {
		return "", error
	}
	return "OK", nil
}
