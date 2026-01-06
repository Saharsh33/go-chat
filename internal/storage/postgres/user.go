package postgres

import "log"

const CreateUserIfNotExistsQuery = `INSERT INTO users (username) VALUES ($1) ON CONFLICT (username) DO NOTHING;
`

func (s *Store) CreateUserIfNotExists(user string) {

	_, err := s.db.Exec(CreateUserIfNotExistsQuery, user)

	if err != nil {
		log.Println(err)
		return
	}

}
