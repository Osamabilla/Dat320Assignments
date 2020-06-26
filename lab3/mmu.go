// +build !solution

package lab3

// MMU is the memory management unit which also contains the structure for the physical memory in this simulation
type MMU struct {
	frames    [][]byte           // contains memory content in form of frames[frameIndex][offset]
	freeList  []bool             // tracks free physical frames
	processes map[int]*PageTable // contains page table for each process (key=pid)
}

// OffsetLUT defines a Lookup Table for offset extracting values (used to translate virtual addresses)
// OffsetLUT[0] --> 0000 ... 0000
// OffsetLUT[1] --> 0000 ... 0001
// OffsetLUT[2] --> 0000 ... 0011
// OffsetLUT[3] --> 0000 ... 0111
// OffsetLUT[8] --> 0000 ... 1111 1111
// etc.
var OffsetLUT = []int{
	// 0000, 0001, 0011, 0111, 1111, etc.
	0x0, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f, 0xff, 0x1ff, 0x3ff, 0x7ff, 0xfff, 0x1fff, 0x3fff, 0x7fff, 0xffff, 0x1ffff, 0x3ffff, 0x7ffff,
	0xfffff, 0x1fffff, 0x3fffff, 0x7fffff, 0xffffff, 0x1ffffff, 0x3ffffff, 0x7ffffff, 0xfffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff, 0xffffffff,
}

// NewMMU creates a new MMU with a memory of `memSize` bytes
// `memSize` should be >= 1 and a multiple of frameSize
func NewMMU(memSize, frameSize int) *MMU {
	return &MMU{}
}

// NumFreeFrames returns the number of free frames in the free list
func (mmu *MMU) NumFreeFrames() int {
	n := 0
	for _, frame := range mmu.freeList {
		if frame {
			n++
		}
	}
	return n
}

// setMemoryContent sets the memory content (mmu.frames) to a certain state.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setMemoryContent(frames [][]byte) {
	mmu.frames = frames
}

// setFreeList sets the free list to a certain state.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setFreeList(freeList []bool) {
	mmu.freeList = freeList
}

// setProcesses sets the state of multiple processes.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setProcesses(processes map[int]*PageTable) {
	mmu.processes = processes
}

// setProcesses sets the state of a single process.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setProcess(pid int, process *PageTable) {
	mmu.processes[pid] = process
}

// Alloc is called by process `pid` which requests to allocate `n` bytes.
// Allocated memory is added to the process' page table.
// The process is given a page table if it doesn't already have one, but only
// if the result of Alloc is to allocate memory to the process (i.e. not
// when an out of memory error occurs).
func (mmu *MMU) Alloc(pid, n int) error {
	return nil
}

// Write is called by process `pid` which wants to write `content` to its address space starting from `virtualAddress`
func (mmu *MMU) Write(pid, virtualAddress int, content []byte) error {
	return nil
}

// Read is called by process `pid` which wants to read `n` bytes from its address space starting from `virtualAddress`
func (mmu *MMU) Read(pid, virtualAddress, n int) (content []byte, err error) {
	return
}

// Free is called by a process' Free() function to free some of its allocated memory
func (mmu *MMU) Free(pid, n int) error {
	return nil
}
