package main

import (
	"fmt"
	"bytes"
	"bufio"
	"os"
	"strings"
	"strconv"
)

const (
	Istr = iota
	Iadd
	Isub
	Imov
	Ildr
	Ibx
	Ibl
	Ilabel
	Istmfd
	Ildmfd
)

var registers []string = []string{
	"r0","r1","r2","r3","r4","r5","r6","r7","r8","r9","r10","r11","r12",
	"sp","lr","pc","apsr"}

const (
	Rr0 = iota
	Rr1
	Rr2
	Rr3
	Rr4
	Rr5
	Rr6
	Rr7
	Rr8
	Rr9
	Rr10
	Rr11
	Rr12
	Rsp
	Rlr
	Rpc
	Rapsr
)

func first(s string) string {
	s = strings.TrimLeft(s," \t")
	for i,v := range s {
		if v == ' ' || v == '\t' {
			return s[:i]
		}
	}
	return s
}

func (p *Parser) strToVal(s string) *Val {
	v := new(Val)
	var err error
	if s[0] == '#' {
		v.offs,err = strconv.Atoi(s[1:])
		if err != nil {
			panic(err)
		}
		return v
	}
	if s[0] == '[' {
		vals := strings.Split(s[1:len(s)-1], ",")
		reg,ok := p.regtab[vals[0]]
		if !ok {
			reg = -1
		}
		v.reg = reg
		v.offs,err = strconv.Atoi(strings.TrimLeft(vals[1], " ")[1:])
		if err != nil {
			panic(err)
		}
		return v
	}
	reg,ok := p.regtab[s]
	if !ok {
		reg = -1
	}
	v.reg = reg
	return v
}

func (p *Parser) regValue(s string) int {
	reg,ok := p.regtab[s]
	if !ok {
		fmt.Printf("Invalid register: '%s'\n", s)
		reg = -1
	}
	return reg
}

func (p *Parser) parseArgs(s string) ([]*Val,error) {
	var vals []*Val
	buf := bytes.NewBufferString(s)
	out := new(bytes.Buffer)
	for {
		b,err := buf.ReadByte()
		if err != nil {
			if out.Len() > 0 {
				s := out.String()
				v := new(Val)
				if s[0] == '#' {
					v.offs,err = strconv.Atoi(s[1:])
					if err != nil {
						panic(err)
					}
					v.reg = -1
				} else {
					v.reg = p.regValue(s)
				}
				vals = append(vals, v)
			}
			return vals,nil
		}
		switch b {
			case ',':
				if out.Len() == 0 {
					continue
				}
				s := out.String()
				out.Reset()
				v := new(Val)
				if s[0] == '#' {
					v.offs,err = strconv.Atoi(s[1:])
					if err != nil {
						panic(err)
					}
					v.reg = -1
				} else {
					v.reg = p.regValue(s)
				}
				vals = append(vals, v)
			case ' ':
				continue
			case '[':
				v := new(Val)
				str,err := buf.ReadString(']')
				spl := strings.Split(str,",")
				if len(spl) == 1 {
					v.reg = p.regValue(spl[0][:len(spl[0])-1])
				} else if len(spl) == 2 {
					v.reg = p.regValue(spl[0])
					num := strings.Trim(spl[1]," #]")
					v.offs,err = strconv.Atoi(num)
					if err != nil {
						panic(err)
					}
				}
				vals = append(vals, v)
				fmt.Println("got compound val.")
				fmt.Println(v)
			default:
				out.WriteByte(b)
		}
	}
	return nil,nil
}

type Parser struct {
	instab map[string]int
	regtab map[string]int
	jmptab map[string]int
}

func NewParser() *Parser {
	p := new(Parser)
	p.instab = make(map[string]int)
	p.instab["add"] = Iadd
	p.instab["sub"] = Isub
	p.instab["str"] = Istr
	p.instab["mov"] = Imov
	p.instab["ldr"] = Ildr
	p.instab["bl"]	= Ibl
	p.instab["bx"]  = Ibx
	p.instab["stmfd"] = Istmfd
	p.instab["ldmfd"] = Ildmfd

	//Set up register mappings
	p.regtab = make(map[string]int)
	for i,v := range registers {
		p.regtab[v] = i
	}
	p.regtab["fp"] = 12

	//Map for labels to jumppoints
	p.jmptab = make(map[string]int)
	p.jmptab["putc"] = -1
	return p
}

func (p *Parser) jumpMap(s string) int {
	i,ok := p.jmptab[s]
	if !ok {
		i = len(p.jmptab)
		p.jmptab[s] = i
	}
	return i
}

func isJump(n int) bool {
	return n == Ibl
}

func (p *Parser) ParseInstruction(ss string) *Instruction {
	var err error
	ss = strings.TrimLeft(ss, " \t")
	instr := first(ss)
	if instr[0] == '@' {
		return nil
	}
	ins := new(Instruction)
	if instr[len(instr)-1] == ':' {
		//fmt.Println("Found label.")
		jind := p.jumpMap(instr[:len(instr)-1])
		ins.params = []*Val{&Val{jind,0}}
		ins.op = Ilabel
		return ins
	}
	op,ok := p.instab[instr]
	if !ok {
		fmt.Printf("Unhandled: '%s'\n", instr)
		return nil
	}
	ins.op = op
	if isJump(op) {
		label := strings.Trim(ss[len(instr):], " \t")
		jv := p.jumpMap(label)
		ins.params = []*Val{&Val{jv,0}}
		return ins
	}

	ins.params,err = p.parseArgs(strings.TrimLeft(ss[len(instr):], " \t"))
	if err != nil {
		fmt.Println(ss)
		panic(err)
	}
	return ins
}

func main() {
	fmt.Println("ARM ASM parser.")
	in, err := os.Open("test.s")
	if err != nil {
		panic(err)
	}
	buf := bufio.NewScanner(in)
	m := NewMachine()
	p := NewParser()
	for buf.Scan() {
		i := p.ParseInstruction(buf.Text())
		if i != nil {
			m.addInstruction(i)
		}
	}
	m.Run()
	fmt.Println("## Program Execution Finished ##")
	fmt.Println("Final register values:")
	for s,v := range m.regs {
		fmt.Printf("%s = %d\n", registers[s],v)
	}
}
