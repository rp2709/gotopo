package gotopo

import "fmt"

type DatabaseConfig struct {
	Dbuser string
	Dbname string
	Dbpwd  string
	Dbhost string
	Dbport int
}

const connectionStringFmt = "user=%s dbname=%s sslmode=prefer host=%s password=%s port=%d"

func (config DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		connectionStringFmt,
		config.Dbuser,
		config.Dbname,
		config.Dbhost,
		config.Dbpwd,
		config.Dbport,
	)
}
