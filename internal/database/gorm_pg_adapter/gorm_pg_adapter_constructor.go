package gorm_pg_adapter

type GormPgAdapterConstructorOption func(adapter *GormPgAdapter)

func CreateGormPgAdapter(
	host string,
	user string,
	password string,
	port int,
	dbname string,
	opts ...GormPgAdapterConstructorOption,
) (*GormPgAdapter, error) {

	newAdapter := &GormPgAdapter{
		host:     host,
		user:     user,
		password: password,
		port:     port,
		dbname:   dbname,
		timezone: "Asia/Singapore",
	}

	for _, opt := range opts {
		opt(newAdapter)
	}

	return newAdapter, nil
}
