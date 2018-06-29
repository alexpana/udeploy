package models

type DeployStatus int

const (
	RUNNING   DeployStatus = 0
	STOPPED   DeployStatus = 1
	DEPLOYING DeployStatus = 2
)

type Project struct {
	Id         string
	Name       string
	ConfigPath string
	Status     DeployStatus
}
