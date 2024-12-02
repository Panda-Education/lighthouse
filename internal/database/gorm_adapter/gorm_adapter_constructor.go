package gorm_adapter

type GormAdapterConstructorOption func(adapter *GormDbAdapter)

func CreateGormAdapter(
	host string,
	user string,
	password string,
	port int,
	dbname string,
	opts ...GormAdapterConstructorOption,
) (*GormDbAdapter, error) {

	newAdapter := &GormDbAdapter{
		host:     host,
		user:     user,
		password: password,
		port:     port,
		dbname:   dbname,
		sslMode:  "disabled",
		timezone: "Asia/Singapore",
	}

	for _, opt := range opts {
		opt(newAdapter)
	}

	return newAdapter, nil
}
