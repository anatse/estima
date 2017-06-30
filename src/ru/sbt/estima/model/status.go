package model

import (
	"strconv"
	"ru/sbt/estima/conf"
)

type Status int

const (
	STATUS_NEW Status = iota
	STATUS_INWORK
	STATUS_REJECTED
	STATUS_APPROVED
	STATUS_ESTIMATION
	STATUS_ESTIMATED
	STATUS_ESTIMATION_APPROVED
	STATUS_DISABLED
)

const (
	RTE = ROLE_RTE
)

type IStatusFSM interface {
	CurrentStatus() Status
	NextStatuses() []Status
	MoveTo(status Status, roles []string) IStatusFSM
	CanMoveTo(status Status, roles []string) bool
}

type StatusFSM struct {
	currentStatus Status
	nextStatuses []*StatusFSM
	allowedRoles []string
}

var status *StatusFSM
var statuses []*StatusFSM

// Function
func GetStatusFSM () IStatusFSM {
	if status == nil {
		// Create statuses
		status = &StatusFSM{currentStatus:STATUS_NEW, allowedRoles: []string {ROLE_BA, ROLE_SA, ROLE_RTE, ROLE_SM}}
		inWork := &StatusFSM{currentStatus:STATUS_INWORK, allowedRoles: []string {ROLE_BA, ROLE_SA, ROLE_RTE, ROLE_SM}}
		rejected := &StatusFSM{currentStatus:STATUS_REJECTED, allowedRoles: []string {ROLE_SA, ROLE_RTE, ROLE_SM, ROLE_ARCHITECT}}
		approved := &StatusFSM{currentStatus:STATUS_APPROVED, allowedRoles: []string {ROLE_SA, ROLE_RTE, ROLE_SM, ROLE_ARCHITECT}}
		estimation := &StatusFSM{currentStatus:STATUS_ESTIMATION, allowedRoles: []string {ROLE_BA, ROLE_PO}}
		estimated := &StatusFSM{currentStatus:STATUS_ESTIMATED, allowedRoles: []string {ROLE_RTE, ROLE_SM, ROLE_ARCHITECT}}
		estimaApproved := &StatusFSM{currentStatus:STATUS_ESTIMATION_APPROVED, allowedRoles: []string {ROLE_ARCHITECT}}

		// Connect status to each other
		status.nextStatuses = []*StatusFSM {inWork}
		inWork.nextStatuses = []*StatusFSM {rejected, approved, status}
		rejected.nextStatuses = []*StatusFSM {inWork}
		approved.nextStatuses = []*StatusFSM {estimation, rejected}
		estimation.nextStatuses = []*StatusFSM {estimated, rejected}
		estimated.nextStatuses = []*StatusFSM {estimaApproved, estimation}

		statuses = []*StatusFSM {
			status,inWork,rejected,approved,estimation,estimated,estimaApproved,
		}
	}

	return *status
}

func FromStatus (s Status) IStatusFSM {
	if len(statuses) == 0 {
		GetStatusFSM ()
	}

	if int(s) > len(statuses) {
		conf.GetLog().Panicf("Status %d unknown", s)
	}

	return statuses[s]
}

func (fsm StatusFSM) CurrentStatus () Status {
	return fsm.currentStatus
}

func (fsm StatusFSM) NextStatuses () []Status {
	if fsm.nextStatuses == nil || len(fsm.nextStatuses) == 0 {
		return nil
	}

	rets := make([]Status, len(fsm.nextStatuses))
	for idx, s := range fsm.nextStatuses {
		rets[idx] = s.CurrentStatus()
	}

	return rets
}

func isAllowed (allowedRoles []string, roles []string) bool {
	for _, role := range roles {
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return true
			}
		}
	}

	return false
}

func (fsm StatusFSM) MoveTo (status Status, roles []string) IStatusFSM {
	for _, s := range fsm.nextStatuses {
		if s.CurrentStatus() == status && isAllowed(fsm.allowedRoles, roles) {
			return *s
		}
	}

	conf.GetLog().Panicf("Not allowed to change status from %d to %d for current user", fsm.currentStatus, status)
	return nil
}

func (fsm StatusFSM) CanMoveTo (status Status, roles []string) bool {
	for _, s := range fsm.nextStatuses {
		if s.CurrentStatus() == status && isAllowed(fsm.allowedRoles, roles) {
			return true
		}
	}

	return false
}

func (status *Status) UnmarshalJSON(b []byte) error {
	stringVal := string(b)
	iVal, err := strconv.Atoi(stringVal)
	*status = Status(iVal)
	return err
}