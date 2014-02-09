package main

import (
	"fmt"
)

type Machine struct {
	srcp *Parser
	regs []int
	mem []uint32
	jumps []int
	datalabels []int
	instr []*Instruction
	c_Z bool
	c_N bool
	c_V bool
	verbose bool
	super bool
}

func NewMachine() *Machine {
	m := new(Machine)
	m.regs = make([]int, len(registers))
	m.jumps = make([]int,1024)
	m.mem = make([]uint32, 8192)
	m.regs[Rsp] = 1024
	m.verbose = false
	m.super = false
	return m
}

func (m *Machine) askVal(v *Val) int {
	if v.reg == -1 {
		return v.offs
	}
	return m.regs[v.reg] + v.offs
}

func (m *Machine) setDataLabel(loc int) int {
	m.datalabels = append(m.datalabels, loc)
	return len(m.datalabels) - 1
}

func (m *Machine) setMem(addr, bytes, val int) {
	word := addr / 4
	b    := uint(addr % 4)
	//fmt.Printf("Mem write at %d\n", addr)
	switch bytes {
		case 4:
			if b != 0 {
				fmt.Printf("error in memory access! pc = %d\n", m.regs[Rpc])
				fmt.Printf("%d %d %d\n", addr, bytes, val)
				fmt.Printf("Instr: %d\n", m.instr[m.regs[Rpc]].op)
				fmt.Println(m.instr[m.regs[Rpc]])
				panic("oh god.")
				return
			}
			m.mem[word] = uint32(val)
		case 1:
			blnk := m.mem[word] & ^(0xff << (8*b))
			blnk |= uint32(byte(val)) << (8*b)
			m.mem[word] = blnk
		default:
			panic("I DONT KNOW WHAT TO DO")
	}
}
func (m *Machine) memAccess(addr, bytes int) int {
	word := addr / 4
	b    := uint32(addr % 4)
	if m.verbose {
		fmt.Printf("Mem access at: %d [%d]\n", word, b)
	}

	switch bytes {
		case 4:
			if b != 0 {
				fmt.Printf("error in memory access! pc = %d\n", m.regs[Rpc])
				return -1
			}
			return int(m.mem[word])
		case 1:
			v := m.mem[word]
			return int((v & (0xff << (b*8))) >> (b*8))
		default:
			panic("I DONT KNOW WHAT TO DO")
	}
	return -1
}

func (m *Machine) Run(main int) error {
	m.regs[Rpc] = m.jumps[main]
	fmt.Printf("Starting at: %d\n", m.jumps[main])
	fmt.Println("## Begin Program Execution ##")
	m.regs[Rlr] = 8192
	for ;m.regs[Rpc] < len(m.instr); m.regs[Rpc]++ {
		m.Exec(m.instr[m.regs[Rpc]])
	}
	return nil
}

func (m *Machine) store(v int, in *Val) {
	//fmt.Printf("Storing in register: %d\n", in.reg)
	m.regs[in.reg] = v
}

func (m *Machine) addInstruction(i *Instruction) {
	if i.op == Ilabel {
		m.jumps[i.params[0].reg] = len(m.instr)
		return
	}
	m.instr = append(m.instr, i)
}

func (m *Machine) getSysFuncID(s string) int {
	switch s {
		case "putc":
			return -1
	}
	return 0
}

func (m *Machine) syscall(n int, i *Instruction) {
	switch n {
		case -1:
			m.Mputchar(i)
	}
}

type Instruction struct {
	op int
	params []*Val
}

type Val struct {
	reg int
	offs int
	shift int
}

func (m *Machine) Exec(i *Instruction) {
	if m.verbose {
		fmt.Printf("pc = %d instr: %s nargs: %d\n", m.regs[Rpc], m.srcp.getInsName(i.op), len(i.params))
	}
	if m.super {
		fmt.Println("Registers:")
		for i := 0; i < len(registers); i++ {
			fmt.Printf("%s = %d\n", registers[i], m.regs[i])
		}
		fmt.Println("Stack:")
		start := m.regs[Rsp]
		for i := 32; i >= -12; i -= 4 {
			fmt.Printf("%d: %d", start + i, m.mem[(start + i)/4])
			if i == 0 {
				fmt.Printf(" <-sp")
			}
			fmt.Println()
		}
		var s string
		fmt.Scanln(&s)
	}
	switch i.op {
	case Iadd:
		a := m.askVal(i.params[1])
		b := m.askVal(i.params[2])
		//fmt.Printf("Adding %d and %d\n", a,b)
		m.store(a + b, i.params[0])
	case Isub:
		a := m.askVal(i.params[1])
		b := m.askVal(i.params[2])
		//fmt.Printf("Subtracting %d from %d\n", b, a)
		m.store(a - b, i.params[0])
	case Istr:
		loc := m.askVal(i.params[1])
		v := m.askVal(i.params[0])
		//fmt.Printf("storing %d in memory location %d\n", v,loc)
		m.setMem(loc, 4, v)
	case Imov:
		v := m.askVal(i.params[1])
		//fmt.Printf("putting %d in %d\n", v, i.params[0].reg)
		m.regs[i.params[0].reg] = v
	case Ildr:
		v := m.memAccess(m.askVal(i.params[1]), 4)
		//fmt.Printf("Loading %d into reg %d\n", v, i.params[0].reg)
		m.regs[i.params[0].reg] = v
	case Ibl:
		//fmt.Printf("Jumping to point: %d\n", i.params[0].reg)
		if i.params[0].reg < 0 {
			m.syscall(i.params[0].reg, i)
		} else {
			//fmt.Printf("Jumping to label, linking back to: %d\n", m.regs[Rpc])
			m.regs[Rlr] = m.regs[Rpc]
			//fmt.Printf("Jumping to %d\n", m.jumps[i.params[0].reg])
			m.regs[Rpc] = m.jumps[i.params[0].reg] - 1
		}
	case Ib:
		if i.params[0].reg < 0 {
			m.syscall(i.params[0].reg, i)
		} else {
			m.regs[Rpc] = m.jumps[i.params[0].reg] - 1
		}
	case Ibx:
		//fmt.Println(i.params[0].reg)
		//fmt.Printf("Returning to: %d\n", m.regs[i.params[0].reg])
		m.regs[Rpc] = m.regs[i.params[0].reg]
	case Istmfd:
		stk := m.regs[i.params[0].reg]
		for _,v := range i.params[1:] {
			stk -= 4
			num := m.regs[v.reg]
			//fmt.Printf("storing %d from %d into %d\n", num, v.reg, stk)
			m.setMem(stk, 4, num)
		}
		m.regs[i.params[0].reg] = stk
	case Ildmfd:
		stk := m.regs[i.params[0].reg]
		for n := len(i.params) - 1; n > 0; n-- {
			v := i.params[n]
			num := m.memAccess(stk, 4)
			//fmt.Printf("loading %d from %d into %d\n", num, stk, v.reg)
			m.regs[v.reg] = num
			stk += 4
		}
		m.regs[i.params[0].reg] = stk
	case Icmp:
		m.c_Z = false
		m.c_N = false
		m.c_V = false

		n := m.askVal(i.params[0]) - m.askVal(i.params[1])
		if m.verbose {
			fmt.Printf("Comparing: %d and %d\n", m.askVal(i.params[0]), m.askVal(i.params[1]))
		}
		if n == 0 {
			m.c_Z = true
		} else if n < 0 {
			m.c_N = true
		}
	case Ible:
		if m.c_N || m.c_Z {
			//fmt.Printf("Jumping to label %d %d\n", i.params[0].reg, i.params[0].offs)
			m.regs[Rpc] = m.jumps[i.params[0].reg] - 1
		}
	case Ibls:
		//Unsigned version of ble
		if m.c_N || m.c_Z {
			m.regs[Rpc] = m.jumps[i.params[0].reg] - 1
		}
	case Ibne:
		if !m.c_Z {
			m.regs[Rpc] = m.jumps[i.params[0].reg] - 1
		}
	case Istrb:
		b := m.regs[i.params[0].reg]
		addr := m.askVal(i.params[1])
		m.setMem(addr, 1, b)
	case Ildrb:
		addr := m.askVal(i.params[1])
		m.regs[i.params[0].reg] = m.memAccess(addr, 1)
	case Imovw:
		cv := m.regs[i.params[0].reg]
		cv = cv & ^0xffff
		m.regs[i.params[0].reg] = cv | i.params[1].offs
	case Imovt:
		cv := m.regs[i.params[0].reg]
		cv = cv & 0xffff
		m.regs[i.params[0].reg] = cv | i.params[1].offs << 16
	case Ismull:
		m1 := int64(m.askVal(i.params[2]))
		m2 := int64(m.askVal(i.params[3]))
		res := m1 * m2
		outl := res & 0xffffffff
		outh := res >> 32
		m.regs[i.params[0].reg] = int(outh)
		m.regs[i.params[1].reg] = int(outl)
	case Irsb:
		m.regs[i.params[0].reg] = m.askVal(i.params[2]) - m.askVal(i.params[1])

	default:
		fmt.Printf("Unhandled instruction: %d\n", i.op)
	}
}

func (m *Machine) Mputchar(i *Instruction) {
	fmt.Printf("%c", m.regs[0])
}
