package kernel

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/infra/config"
)

var conf = config.GetInstance()
var urlSyscall = fmt.Sprintf("http://%s:%d/kernel/syscall", conf.IPKernel, conf.PortKernel)
var urlInterrupt = fmt.Sprintf("http://%s:%d/kernel/interrupt", conf.IPKernel, conf.PortKernel)
