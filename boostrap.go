package gofkass

var boostrap *Boostrap

type Boostrap struct {
	PreSetupDB func()
	SetupApp   func()
	Register   func(aliasName, driverName, dataSource string, params ...int) error
	Sync       func(name string, force bool, verbose bool) error
}

func GetBoostrap() *Boostrap{
	return boostrap
}

func SetBoostrap(boostrap_ *Boostrap){
	boostrap = boostrap_
}