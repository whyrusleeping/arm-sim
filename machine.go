package main

import (
	"fmt"
)

type Machine struct {
	regs []int
	mem []int
	jumps []int
	instr []*Instruction
}

func NewMachine() *Machine {
	m := new(Machine)
	m.regs = make([]int, len(registers))
	m.jumps = make([]int,1024)
	m.mem = make([]int, 8192)
	m.regs[Rsp] = 1024
	return m
}

func (m *Machine) askVal(v *Val) int {
	if v.reg == -1 {
		return v.offs
	}
	return m.regs[v.reg] + v.offs
}

func (m *Machine) Run() error {
	m.regs[Rpc] = 0
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
}

func (m *Machine) Exec(i *Instruction) {
	//fmt.Printf("instr: %d nargs: %d\n", i.op, len(i.params))
	switch i.op {
	case Iadd:
		a := m.askVal(i.params[1])
		b := m.askVal(i.params[2])
		//fmt.Printf("Adding %d and %d\n", a,b)
		m.store(a + b, i.params[0])
	case Isub:
		m.store(m.askVal(i.params[1]) - m.askVal(i.params[2]), i.params[0])
	case Istr:
		loc := m.askVal(i.params[1])
		v := m.askVal(i.params[0])
		//fmt.Printf("storing %d in memory location %d\n", v,loc)
		m.mem[loc] = v
	case Imov:
		v := m.askVal(i.params[1])
		//fmt.Printf("putting %d in %d\n", v, i.params[0].reg)
		m.regs[i.params[0].reg] = v
	case Ildr:
		m.regs[i.params[0].reg] = m.mem[m.askVal(i.params[1])]
	case Ibl:
		//fmt.Printf("Jumping to point: %d\n", i.params[0].reg)
		if i.params[0].reg < 0 {
			m.syscall(i.params[0].reg, i)
		} else {
			m.regs[Rlr] = m.regs[Rpc]
			m.regs[Rpc] = m.jumps[i.params[0].reg]
		}
	case Ibx:
		fmt.Println(i.params[0].reg)
		//fmt.Printf("Returning to: %d\n", m.regs[i.params[0].reg])
		m.regs[Rpc] = m.regs[i.params[0].reg]
	default:
		fmt.Printf("Unhandled instruction: %d\n", i.op)
	}
}

func (m *Machine) Mputchar(i *Instruction) {
	fmt.Printf("%c", m.regs[0])
	//fmt.Print(string([]byte{byte(m.regs[0])}))
}
