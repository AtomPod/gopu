package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ngs24313/gopu/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

//Hardware is hardware information api of machine
type Hardware struct {
}

//Register handlers
func (h *Hardware) Register(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		v1.GET("/cpu", h.CPUGet)
		v1.GET("/cpu/times", h.CPUTimesGet)
		v1.GET("/memory", h.MemoryGet)
		v1.GET("/disk", h.DiskGet)
		v1.GET("/process", h.ProcessGet)
		v1.GET("/process/:pid/memory", h.ProcessMemoryGet)
	}
}

//CPUGet handles GET /v1/cpu
func (h *Hardware) CPUGet(c *gin.Context) {
	infos, err := cpu.Info()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}

//CPUTimesGet handles GET  /v1/cpu/times
func (h *Hardware) CPUTimesGet(c *gin.Context) {
	infos, err := cpu.Times(true)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}

//MemoryGet handles GET /v1/memory
func (h *Hardware) MemoryGet(c *gin.Context) {
	vm, err := mem.VirtualMemory()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	memory := &models.Memory{
		Total:       vm.Total,
		Available:   vm.Available,
		Used:        vm.Used,
		Freed:       vm.Free,
		UsedPercent: vm.UsedPercent,
	}
	c.JSON(http.StatusOK, memory)
}

//DiskGet handles GET /disk
func (h *Hardware) DiskGet(c *gin.Context) {
	stats, err := disk.Partitions(true)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	diskInfo := make([]struct {
		Partition *disk.PartitionStat `json:"partition,omitempty"`
		Usage     *disk.UsageStat     `json:"usage,omitempty"`
	}, len(stats))

	for i, stat := range stats {
		usage, err := disk.Usage(stat.Device)
		if err == nil {
			diskInfo[i].Usage = usage
		}
		var copyStat = stat
		diskInfo[i].Partition = &copyStat
	}
	c.JSON(http.StatusOK, diskInfo)
}

//ProcessGet handles GET /v1/process
func (h *Hardware) ProcessGet(c *gin.Context) {
	var page = int(1)
	var pageSize = int(16)

	if pageString := c.Query("page"); pageString != "" {
		p, err := strconv.Atoi(pageString)
		if err != nil {
			c.String(http.StatusBadRequest, "param [page] is not a integer")
			return
		}
		if p >= 1 {
			page = p
		}
	}

	if pageSizeString := c.Query("page_size"); pageSizeString != "" {
		size, err := strconv.Atoi(pageSizeString)
		if err != nil {
			c.String(http.StatusBadRequest, "param [page_size] is not a integer")
			return
		}
		if size <= 64 && size >= 16 {
			pageSize = size
		}
	}

	processes, err := process.Processes()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	plen := len(processes)
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= plen {
		startIndex = plen
	}

	if endIndex >= plen {
		endIndex = plen
	}

	processList := make([]*models.Process, endIndex-startIndex)
	for i := 0; i < len(processList); i++ {
		processList[i] = toProcessesModel(processes[i+startIndex])
	}
	c.JSON(http.StatusOK, processList)
}

//ProcessMemoryGet handles GET /v1/process/:pid/memory
func (h *Hardware) ProcessMemoryGet(c *gin.Context) {

	pidStr := c.Param("pid")
	if pidStr == "" {
		c.String(http.StatusBadRequest, "pid cannot be empty")
		return
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.String(http.StatusBadRequest, "pid must be a integer")
		return
	}

	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		if err == process.ErrorProcessNotRunning {
			c.String(http.StatusNotFound, "pid [%d] is not exists", pid)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	memoryInfo, err := proc.MemoryInfo()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, memoryInfo)
}

func toProcessesModel(p *process.Process) *models.Process {
	process := &models.Process{
		Pid: p.Pid,
	}

	if bg, err := p.Background(); err == nil {
		process.Background = bg
	}

	if name, err := p.Name(); err == nil {
		process.Name = name
	}

	if cpuPercent, err := p.CPUPercent(); err == nil {
		process.CPUPercent = cpuPercent
	}

	if cmdLine, err := p.Cmdline(); err == nil {
		process.Cmdline = cmdLine
	}

	if createTime, err := p.CreateTime(); err == nil {
		process.CreateTime = createTime
	}

	if cwd, err := p.Cwd(); err == nil {
		process.Cwd = cwd
	}

	if running, err := p.IsRunning(); err == nil {
		process.IsRunning = running
	}

	if numOfThread, err := p.NumThreads(); err == nil {
		process.NumOfThreads = numOfThread
	}

	if memPercent, err := p.MemoryPercent(); err == nil {
		process.MemoryPercent = memPercent
	}

	if numOfFD, err := p.NumFDs(); err == nil {
		process.NumOfFDs = numOfFD
	}

	return process
}
