package db

func WithInsertIgnore(b bool) InsertOption {
	return InsertOption(b)
}

