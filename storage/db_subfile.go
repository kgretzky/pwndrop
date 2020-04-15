package storage

type DbSubFile struct {
	ID         int    `json:"id" storm:"id,increment"`
	Fid        int    `json:"fid"`
	Uid        int    `json:"uid"`
	Name       string `json:"name"`
	Filename   string `json:"fname"`
	FileSize   int64  `json:"fsize"`
	CreateTime int64  `json:"create_time"`
}

func SubFileCreate(o *DbSubFile) (*DbSubFile, error) {
	err := db.Save(o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func SubFileGet(id int) (*DbSubFile, error) {
	var o DbSubFile
	err := db.One("ID", id, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func SubFileDelete(id int) error {
	f := &DbSubFile{
		ID: id,
	}
	err := db.DeleteStruct(f)
	if err != nil {
		return err
	}
	return nil
}
