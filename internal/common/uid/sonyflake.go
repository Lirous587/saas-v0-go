package uid

import (
	"saas/internal/common/utils"
	"sync"
	"time"

	"github.com/sony/sonyflake/v2"
)

var (
	sonyInstance *sonyflake.Sonyflake
	once         sync.Once
)

func Init() {
	once.Do(func() {
		startTime, err := time.Parse(time.RFC3339, utils.GetEnv("SONYFLAKE_START_TIME"))
		if err != nil {
			panic("Invalid SONYFLAKE_START_TIME format")
		}

		machineID := utils.GetEnvAsInt("SONYFLAKE_MACHINE_ID")

		sony, err := sonyflake.New(sonyflake.Settings{
			StartTime: startTime,
			MachineID: func() (int, error) {
				return machineID, nil
			},
		})
		if err != nil {
			panic(err)
		}
		sonyInstance = sony
	})
}

func Gen() (int64, error) {
	return sonyInstance.NextID()
}
