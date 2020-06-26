# Lab 3: Segmented Memory Allocation with Paging

| Lab 3:                | Memory Management                                |
| --------------------- | ------------------------------------------------ |
| Subject:              | DAT320 Operating Systems and Systems Programming |
| Deadline:             | **TODO**                                         |
| Expected effort:      | 10-15 hours                                      |
| Grading:              | Pass/fail                                        |
| Submission:           | Individually                                     |

## Introduction

In this task you will implement various low-level functions used by the Memory Management Unit (*MMU*).
By the end of this task you will have a working MMU which manages a memory divided into fixed-size physical page frames.
Frames can be allocated to or deallocated from processes, and processes will also be able to write to and read from their allocated memory.
The memory segmentation, address translation, etc. in this lab is based on [Chapter 18 of the textbook](http://pages.cs.wisc.edu/~remzi/OSTEP/vm-paging.pdf).

We have included a few files which you must use as a starting point for your assignment.

* `mmu.go`. This is where most of your code will be.
  It includes a skeleton of the `MMU` struct with some of its necessary data structures, as well as the function header for some necessary functions.
* `process.go`. The code in this file is already complete.
  It contains the `Process` struct which calls various functions of the `MMU` struct to allocate and use memory.
* `pagetable.go`. This file contains the `PageTable` struct which is used by the MMU to track translations from virtual page numbers to physical frame numbers.

After most tasks you should run a test to verify that your code is working as expected.
This step is important, because these tests and more are used in Autograder, and later parts of the lab depends on the code working as expected.

## Memory Management Questions

Answer the following questions in the file [`mmu-answers.md`](mmu-answers.md).

1. How is memory divided using memory management based on simple paging?
   * a. fixed-sized chunks
   * b. variable-sized chunks
   * c. hybrid chunks
   * d. virtual chunks

2. Which option is *not* a problem when using paging?
   * a. internal fragmentation
   * b. storing page tables wastes memory
   * c. external fragmentation
   * d. address translation requires extra memory accesses (reading the page table)

3. What is the purpose of the free list in paged memory?
   * a. keeping track of unused virtual page frames
   * b. keeping track of unused virtual memory segments
   * c. keeping track of unused virtual address spaces
   * d. keeping track of unused physical page frames

## Implementing the MMU for Paged Memory

### Architecture

In this section you will implement the core functionality of the MMU.
In this lab, the memory is divided using a method called paging, where the memory is called *paged memory*.
Paged memory consists of several fixed-size *physical frames*. Physical frames have two possible states: *free* and *allocated*.
The state of all physical frames of the memory is tracked in the *free list*.
The architecture is illustrated in the figure below (note that the what is denoted as *core map* in the figure is what we call the *free list*).

![An image illustrating the paged memory architecture described in this section](./fig/paging-architecture.png "Paged memory architecture")

Each physical frame can be be allocated to a process.
Allocated frames are added to the process' *virtual address space*.
Here, the process looks at its allocated memory as a continuous memory segment starting from the address 0, which the process has full control over.
Later on, a process will be able to do many operations on its address space, such as read from and write to it, expand it, and free physical frames from it.

To keep track of the virtual address space of each process, we use a structure called a *page table*.
The page table is a per-process data structure which contains translations from virtual page numbers (VPNs) to physical frame numbers.
Whenever a physical frame is allocated to a process, a translation from a virtual page number to a physical frame number is appended to its page table.

As mentioned, the state of physical frames will be tracked in the *free list*, `MMU.freeList`.
The length of the free list equals the number of physical frames.
The free list contains an entry for each physical frame which contains a boolean to indicate whether a physical frame is free, i.e.
available for allocation, or already allocated.
For simplicity, the index of an entry in the free list corresponds to the physical frame number it represents, e.g. `freeList[0]` shows the state of physical frame 0 and `freeList[7]` shows the state of physical frame 7.

`MMU.processes` is a map used to store references to each process' page table, where the key is the PID and the value is a pointer to the corresponding process' page table.

### Task 1: Memory Allocation

In this task you will structure the memory and implement allocation.
You must implement the `NewMMU()` and `MMU.Alloc()` functions.

#### Initialization

First, you must implement the `NewMMU` function.
Note that the memory is divided with `MMU.frames[i][j]`, where `i` signifies the physical frame number and `j` signifies the offset within that frame.
The frame size is a parameter to `NewMMU()`.
Use these factors to divide `MMU.frames` into the correct amount of frames.
Also, each frame must have an entry in `mmu.freeList` to represent its state.

#### Allocating Memory

Furthermore, you must implement allocation.
Allocation occurs when a process sends a request to the MMU with `Process.Malloc()`, which requests a certain amount of bytes of memory.
The MMU must check whether enough free memory is available, and if it is, must find a number of pages by checking the free list, and add these pages to the process' page table.

Add the following functionality to `MMU.Alloc()`:

1. Calculate the number of frames the process needs.
   Memory is only allocated as physical frames.
   As such, the amount of memory allocated to a process is always a multiple of the frame size.
   E.g. if the frame size is 4, and a process requests 6 bytes, the MMU will allocate 2 frames to the process which have room for 8 bytes in total.
   If there are not enough frames available, the function should return an error.
2. Find which frames to be allocated by scanning the free list.
   You should always pick free frames with the lowest physical frame number (i.e. lowest free list index).
   Remember to update the free list to reflect which frames have been allocated.
3. Add the selected frames to the process' page table.
   If the process doesn't have a page table yet you need to create one and add it to `MMU.processes`.
   Use `PageTable.Append()` to add pages to the process' page table.

#### Tests

Run the following commands to test your solution:

```sh
go test -run TestNewMMU
go test -run TestAlloc
go test -run TestPTAppend
```

### Task 2: Writing and Reading

In this task, you will implement write and read functions for the MMU.

#### Information: Virtual Addresses

As explained earlier, processes have a virtual address space which the MMU keeps track off by storing a page table for each process.
Since each virtual address space starts from the address 0 in the perspective of the process, processes use something called a *virtual address* to determine where they want to write to or read from within their address space.
A valid virtual address can be translated to a single physical memory address, which in the case of this lab refers to a memory cell containing a single byte.

A virtual address has a special format.
In the binary format, the *len* - *n* highest order bits determine the *virtual page number*, while the *n* lowest order bits determine the *offset* within the *physical frame* pointed to by the aforementioned virtual page.
*len* is the total length of the virtual address, while *n* is given by `log_2(frame_size)`.

The virtual address translation process is illustrated in the figure below:

![An illustration of the virtual address translation process as described in this section. A process p sends a request Read(19), which is translated into VPN=2 and offset=3 (split after 3rd bit from the right since frame size = 8 = 2^3).](./fig/paging-virtual-address.png "Illustration of virtual address translation in paged memory")

In the figure, *p* sends the request `Read(19)` to the MMU.
Since the frame size is 8, the address is split after the third bit.
`VPN = 2` means we use the physical frame number translation at index 2 of the page table of *p*, which gives us the physical frame 5.
Finally we access the byte at `offset=3` within physical frame 5.

#### Tasks

1. Implement `PageTable.Lookup()` which must be used for virtual page number to physical frame number translation for the process that wants to read/write.
2. Implement `MMU.Write()`. A process wants to write *n* bytes to its address space, starting from a certain virtual address *a*.
   For each of the *n* bytes in the sequence you have to write to the next virtual address in order.
   1. If the starting address is illegal, i.e. outside of the process' address space, return an error.
   2. If the starting address is legal, but starting address + *n* exceeds the process' address space, the MMU should try to allocate as much memory as necessary to the process in order for it to write.
   3. If there is not enough free memory to allocate to the process in the case above, return an error without allocating anything.
3. Implement `MMU.Read()`. Similarly, a process wants to read *n* bytes from its address space, starting from a certain virtual address *a*.
   The *n* bytes have to be returned in the order of their virtual address.
   1. If the starting address or the final address is illegal, i.e. outside of the process' address space, return an error.

#### Tests

Run the following commands to test your solution:

```sh
go test -run TestPTLookup
go test -run TestWrite
go test -run TestRead
```

### Task 3: Freeing Memory

In this task, you must implement the `MMU.Free()` function to allow processes to free pages.
`Process.Free()` tells the MMU that the process wants to free *n* pages from its address space.
The process to free memory involves several steps.

1. Implement `PageTable.Free()` which must free the *n* last entries from the process' page table.
2. Zero out the content of the memory that was freed, so that other processes cannot read the previous memory contents if the memory is reallocated.
3. Add the frames that were freed to the free list.

#### Tests

Run the following commands to test your solution:

```sh
go test -run TestFree
go test -run TestPTFree
```

You should also run the `TestSequences` test which will test a combination of allocation, reading, writing and freeing:

```sh
go test -run TestSequences
```
