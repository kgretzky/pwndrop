package storage

type DbSession struct {
	ID         int    `json:"id" storm:"id,increment"`
	Uid        int    `json:"uid" storm:"index"`
	Token      string `json:"token" storm:"unique"`
	CreateTime int64  `json:"create_time"`
}

func SessionCreate(o *DbSession) (*DbSession, error) {
	err := db.Save(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func SessionGet(id int) (*DbSession, error) {
	var o DbSession
	err := db.One("ID", id, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func SessionGetByToken(token string) (*DbSession, error) {
	var o DbSession
	err := db.One("Token", token, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func SessionDelete(id int) error {
	o := &DbSession{
		ID: id,
	}
	err := db.DeleteStruct(o)
	if err != nil {
		return err
	}
	return nil
}

func SessionDeleteAll() error {
	err := db.Drop("DbSession")
	if err != nil {
		return err
	}
	return nil
}
