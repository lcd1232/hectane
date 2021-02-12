package queue

func noopFailProcess(m *Message, s *Storage, err error) error {
	return nil
}
