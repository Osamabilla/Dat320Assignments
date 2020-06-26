package lab3

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewMMU(t *testing.T) {
	for i, test := range NewMMUTests {
		mmu := NewMMU(test.memSize, test.frameSize)

		var explain bool
		if diff := cmp.Diff(test.wantFrames, mmu.frames); diff != "" {
			explain = true
			t.Errorf("TestNewMMU %d: Unexpected state of memory (mmu.frames) after NewMMU(memSize = %d, frameSize = %d); (-want +got):\n%s", i, test.memSize, test.frameSize, diff)
		}

		if diff := cmp.Diff(test.wantFreeList, mmu.freeList); diff != "" {
			explain = true
			t.Errorf("TestNewMMU %d: Unexpected free list state after NewMMU(memSize = %d, frameSize = %d); (-want +got):\n%s", i, test.memSize, test.frameSize, diff)
		}

		if explain {
			t.Logf("Description of this test: \n\t%s", test.desc)
		}
	}
}

func TestAlloc(t *testing.T) {
	for i, test := range AllocTests {
		mmu := NewMMU(test.memSize, test.frameSize)
		mmu.setFreeList(test.freeList)

		err := mmu.Alloc(0, test.in)

		if test.err == nil && err != nil || test.err != nil && err == nil {
			t.Logf("Explanation of this test case: %s", test.desc)
			t.Fatalf("AllocTests %d: Unexpected error result after Alloc operation; want '%v', got '%v'", i, test.err, err)
		}

		if diff := cmp.Diff(test.pageTable, mmu.processes[0], cmpOptPageTable); diff != "" {
			t.Errorf("AllocTests %d: Invalid page table for process %d; (-want +got):\n%s", i, 0, diff)
			t.Logf("Explanation of this test case: %s", test.desc)
		}
	}
}

func TestAllocMultiple(t *testing.T) {
	for i, test := range AllocMultipleTests {
		mmu := NewMMU(test.memSize, test.frameSize)
		mmu.setFreeList(test.freeListState)
		for j, operation := range test.operations {
			var explain bool

			err := mmu.Alloc(operation.pid, operation.in)
			if operation.wantError == nil && err != nil || operation.wantError != nil && err == nil {
				t.Logf("Explanation of this step: %s", operation.desc)
				t.Fatalf("AllocMultipleTests %d, operation %d: Unexpected error result after Alloc operation; want '%v', got '%v'", i, j, operation.wantError, err)
			}

			if diff := cmp.Diff(operation.wantPageTable, mmu.processes[operation.pid], cmpOptPageTable); diff != "" {
				explain = true

				if operation.wantError != nil {
					t.Errorf("AllocMultipleTests %d, operation %d: Invalid page table for process %d. \nNOTE: expected error for this operation and thus no changes should have occured since the last successful operation; (-want +got):\n%s", i, j, operation.pid, diff)
				} else {
					t.Errorf("AllocMultipleTests %d, operation %d: Invalid page table for process %d; (-want +got):\n%s", i, j, operation.pid, diff)
				}
			}

			if diff := cmp.Diff(operation.wantFreeList, mmu.freeList); diff != "" {
				explain = true

				if operation.wantError != nil {
					t.Errorf("AllocMultipleTests %d, operation %d: Invalid free list state. \nNOTE: expected error for this operation and thus no changes should have occured since the last successful operation; (-want +got):\n%s", i, j, diff)
				} else {
					t.Errorf("AllocMultipleTests %d, operation %d: Invalid free list state; (-want +got):\n%s", i, j, diff)
				}
			}

			if explain {
				t.Logf("Explanation of this step: %s", operation.desc)
			}
		}
	}
}

func TestRead(t *testing.T) {
	for i, test := range ReadTests {
		mmu := NewMMU(test.memSize, test.frameSize)
		mmu.setMemoryContent(test.frames)
		if test.pageTable != nil {
			mmu.setProcess(0, test.pageTable)
		}

		content, err := mmu.Read(0, test.addr, test.n)
		if test.err == nil && err != nil || test.err != nil && err == nil {
			t.Logf("Description of this test: \n\t%s", test.desc)
			t.Fatalf("ReadTests %d: Unexpected error result after Read operation; want '%v', got '%v'", i, test.err, err)
		}

		if diff := cmp.Diff(test.content, content); diff != "" {
			t.Errorf("ReadTests %d: Unexpected content from read; (-want +got):\n%s", i, diff)
			t.Logf("Description of this test: \n\t%s", test.desc)
		}
	}
}

func TestWrite(t *testing.T) {
	for i, test := range WriteTests {
		mmu := NewMMU(test.memSize, test.frameSize)
		mmu.setMemoryContent(test.frames)
		mmu.setFreeList(test.freeList)
		if test.pageTable != nil {
			mmu.setProcess(0, test.pageTable)
		}
		var explain bool

		err := mmu.Write(0, test.addr, test.content)
		if test.err == nil && err != nil || test.err != nil && err == nil {
			t.Logf("Description of this test: \n\t%s", test.desc)
			t.Fatalf("WriteTests %d: Unexpected error result after Write operation; want '%v', got '%v'", i, test.err, err)
		}

		if diff := cmp.Diff(test.wantPageTable, mmu.processes[0], cmpOptPageTable); diff != "" {
			explain = true
			t.Errorf("WriteTests %d: Unexpected page table content after write; (-want +got):\n%s", i, diff)
		}

		if diff := cmp.Diff(test.wantFrames, mmu.frames); diff != "" {
			explain = true
			t.Errorf("WriteTests %d: Unexpected memory content content after write; (-want +got):\n%s", i, diff)
		}

		if diff := cmp.Diff(test.wantFreeList, mmu.freeList); diff != "" {
			explain = true
			t.Errorf("WriteTests %d: Unexpected free list state after write; (-want +got):\n%s", i, diff)
		}

		if explain {
			t.Logf("Description of this test: \n\t%s", test.desc)
		}
	}
}

func TestFree(t *testing.T) {
	for i, test := range FreeTests {
		mmu := NewMMU(test.memSize, test.frameSize)
		mmu.setProcesses(test.processes)
		mmu.setFreeList(test.freeList)
		mmu.setMemoryContent(test.frames)

		err := mmu.Free(test.pid, test.n)

		if test.err == nil && err != nil || test.err != nil && err == nil {
			t.Logf("Description of this test: \n\t%s", test.desc)
			t.Fatalf("FreeTests %d: Unexpected error result after Free(pid = %d, n = %d) operation; want '%v', got '%v'", i, test.pid, test.n, test.err, err)
		}

		var explain bool

		if diff := cmp.Diff(test.wantProcesses, mmu.processes, cmpOptPageTable); diff != "" {
			explain = true
			t.Errorf("FreeTests %d: Unexpected state of processes' page tables after Free(pid = %d, n = %d) operation; (-want +got):\n%s", i, test.pid, test.n, diff)
		}

		if diff := cmp.Diff(test.wantFrames, mmu.frames); diff != "" {
			explain = true
			t.Errorf("FreeTests %d: Unexpected state of memory (mmu.frames) after Free(pid = %d, n = %d) operation; (-want +got):\n%s", i, test.pid, test.n, diff)
		}

		if diff := cmp.Diff(test.wantFreeList, mmu.freeList); diff != "" {
			explain = true
			t.Errorf("FreeTests %d: Unexpected state of free list after Free(pid = %d, n = %d) operation; (-want +got):\n%s", i, test.pid, test.n, diff)
		}

		if explain {
			t.Logf("Description of this test: \n\t%s", test.desc)
		}
	}
}

func TestSequences(t *testing.T) {
	for i, test := range SequenceTests {
		mmu := NewMMU(test.memSize, test.frameSize)

		// perform each command in sequence, checking the results of each command
		for j, cmd := range test.cmds {
			var failed bool

			switch cmd.cmd {
			case cmdAlloc:
				err := mmu.Alloc(cmd.pid, cmd.n)
				failed = testSequencesInternalCheckErr(t, i, j, cmd, err, nil)
			case cmdWrite:
				err := mmu.Write(cmd.pid, cmd.vaddr, cmd.content)
				failed = testSequencesInternalCheckErr(t, i, j, cmd, err, nil)
			case cmdRead:
				content, err := mmu.Read(cmd.pid, cmd.vaddr, cmd.n)
				failed = testSequencesInternalCheckErr(t, i, j, cmd, err, content)
			case cmdFree:
				err := mmu.Free(cmd.pid, cmd.n)
				failed = testSequencesInternalCheckErr(t, i, j, cmd, err, nil)
			}

			if failed {
				t.Fatalf("Stopping test %d after failure on step %d. Sequence of commands along with their status: \n%s", i, j, testSequencesCmdSequenceDescriber(test.cmds, j, true))
			}
		}

		var explain bool

		if diff := cmp.Diff(test.wantFreeList, mmu.freeList); diff != "" {
			explain = true
			t.Errorf("SequenceTests %d: Unexpected free list state after command sequence; (-want +got):\n%s", i, diff)
		}

		if diff := cmp.Diff(test.wantFrames, mmu.frames); diff != "" {
			explain = true
			t.Errorf("SequenceTests %d: Unexpected memory content (mmu.frames) after command sequence; (-want +got):\n%s", i, diff)
		}

		if diff := cmp.Diff(test.wantProcesses, mmu.processes, cmpOptPageTable); diff != "" {
			explain = true
			t.Errorf("SequenceTests %d: Unexpected state of processes' page tables after command sequence; (-want +got):\n%s", i, diff)
		}

		if explain {
			t.Logf("Sequence of commands in test %d along with their status: \n%s", i, testSequencesCmdSequenceDescriber(test.cmds, len(test.cmds)-1, false))
		}
	}
}

// checks if outputs from a command in the sequence matches what is expected
func testSequencesInternalCheckErr(t *testing.T, i, j int, cmd mmuCmd, gotErr error, gotContent []byte) (failed bool) {
	if cmd.err == nil && gotErr != nil || cmd.err != nil && gotErr == nil {
		t.Errorf("SequenceTests %d, step %d: Unexpected error result after command '%s'; want '%v', got '%v'", i, j, cmd, cmd.err, gotErr)
		failed = true
	}

	// in case of Read command, check if result matches
	if cmd.cmd == cmdRead {
		if diff := cmp.Diff(cmd.content, gotContent); diff != "" {
			t.Errorf("SequenceTests %d, step %d: Unexpected content returned from command '%s'; (-want +got):\n%s", i, j, cmd, diff)
			failed = true
		}
	}

	return
}

// produces a string describing all commands in the sequence, along with their
// status and description. The final command is indicated as a failure.
// Used to give students an idea of where the test went wrong.
func testSequencesCmdSequenceDescriber(cmds []mmuCmd, final int, finalFailed bool) string {
	var str, finalStatus string
	for i := 0; i < final; i++ {
		str += fmt.Sprintf("Status: SUCCESS; Command: %s; \n\tDescription: %s\n", cmds[i], cmds[i].desc)
	}

	if finalFailed {
		finalStatus = "FAILED"
	} else {
		finalStatus = "SUCCESS"
	}
	str += fmt.Sprintf("Status: %s; Command: %s; \n\tDescription: %s\n", finalStatus, cmds[final], cmds[final].desc)
	return str
}
