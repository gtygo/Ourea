package server

//DispatchCommand gfhjkljhgfhjkl;
func (s *Server) DispatchCommand(info []string) ([]string, error) {

	cmd := info[0]
	data := info[1:]
	var error error

	switch cmd {
	case "HSET":
		if !checkDataLen(len(data), 3) {
			return nil, ErrHsetArgs
		}
		s.c.Hset(data[0], data[1], data[2])
	case "HGET":
		if !checkDataLen(len(data), 2) {
			return nil, ErrHgetArgs
		}
		v, err := s.c.Hget(data[0], data[1])
		if err != nil {
			return nil, err
		}
		return []string{v}, nil

	case "INCR":

	case "DECR":

	case "KEYS":
		if len(data) < 1 {
			return nil, ErrKeysArgs
		}
		if data[0] == "*" {
			return s.c.GetAllKey(), nil
		}
	case "SET":
		if len(data) < 2 {
			return nil, ErrSetArgs
		}
		s.c.Set(data[0], data[1])
	case "GET":
		if len(data) < 1 {
			return nil, ErrGetArgs
		}
		v, err := s.c.Get(data[0])
		if err != nil {
			return nil, err
		}
		return []string{v}, nil
	case "DEL":
		if len(data) < 1 {
			return nil, ErrDelArgs
		}
		s.c.Del(data[0])

	case "SAVE":
		if !checkDataLen(len(data),0){
			return nil,ErrSaveArgs
		}
		if err:=s.c.Dump();err!=nil{
			return nil,err
		}

	case "BGSAVE":


	}
	if error != nil {
		return nil, error
	}
	return []string{"OK"}, nil
}

func checkDataLen(gotLen int,wantLen int)bool{
	if gotLen<wantLen{
		return false
	}
	return true
}
