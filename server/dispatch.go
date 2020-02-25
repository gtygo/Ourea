package server

//DispatchCommand gfhjkljhgfhjkl;
func (s *Server) DispatchCommand(info []string) (string, error) {

	cmd := info[0]
	data := info[1:]
	var error error

	switch cmd {
	case "SET":
		error = s.Rpc.Set([]byte(data[0]), []byte(data[1]))
	case "GET":
		v, err := s.Rpc.Get([]byte(data[0]))
		if err != nil {
			return err.Error(), err
		}
		return string(v), nil
	case "DEL":
		error = s.Rpc.Del([]byte(data[0]))
	}
	if error != nil {
		return "", error
	}
	return "OK", nil
}
