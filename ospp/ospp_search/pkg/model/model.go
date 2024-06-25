package model

type Project struct {
	Id string
}

var pro Project

func SetProject(project Project) {

	pro = project
}

func GetProject() Project {

	return pro
}
