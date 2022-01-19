package services

type Cloner interface {
	CloneAll()
	CloneProject(ID int)
	CloneGroupProjects(ID int)
	CloneGroupProjectsRecursive(ID int)
	CloneGroups(IDs []int)
	CloneGroupsRecursive(IDs int)
	GetGroupList()
	GetProjectsOfGroup(ID int) []interface{}
}
