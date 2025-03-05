package snowflakex

import (
	"github.com/sony/sonyflake"
	"time"
)

type Node interface {
	NextID() int64
}

type SonyFlake struct {
	flake *sonyflake.Sonyflake
}

func NewNode(startTime time.Time, machineId int) Node {
	return &SonyFlake{
		flake: sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: startTime,
			MachineID: func() (uint16, error) {
				return uint16(machineId), nil
			},
			CheckMachineID: func(u uint16) bool {
				return u == uint16(machineId)
			},
		}),
	}
}

func (node *SonyFlake) NextID() int64 {
	id, _ := node.flake.NextID()
	return int64(id)
}
