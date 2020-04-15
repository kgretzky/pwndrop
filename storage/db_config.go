package storage

type DbConfig struct {
	ID          int    `json:"id" storm:"id"`
	Hostname    string `json:"hostname"`
	SecretPath  string `json:"secret_path"`
	RedirectUrl string `json:"redirect_url"`
	CookieName  string `json:"cookie_name"`
	CookieToken string `json:"cookie_token"`
}

func ConfigCreate(o *DbConfig) (*DbConfig, error) {
	err := db.Save(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func ConfigGet(id int) (*DbConfig, error) {
	var o DbConfig
	err := db.One("ID", id, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func ConfigUpdate(id int, o *DbConfig) (*DbConfig, error) {
	o.ID = id
	if err := db.Save(o); err != nil {
		return nil, err
	}
	return o, nil
}

func ConfigDelete(id int) error {
	o := &DbConfig{
		ID: id,
	}
	err := db.DeleteStruct(o)
	if err != nil {
		return err
	}
	return nil
}
