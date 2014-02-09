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
	Ib
	Ibl
	Ilabel
	Istmfd
	Ildmfd //10
	Icmp
	Ible
	Istrb
	Ildrb
	Ibls
	Ibne
	Imovw
	Imovt
	Ismull
	Irsb  //20
)

//Sections
const (
	Mtext = iota
	Mdata
)

var registers = []string{
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
	if s[len(s)-1] == '!' {
		if p.verbose {
			fmt.Println("Mod In Place Reg.")
		}
		s = s[:len(s)-1]
	}
	reg,ok := p.regtab[s]
	if !ok {
		fmt.Printf("Invalid register: '%s'\n", s)
		reg = -1
	}
	return reg
}

func (p *Parser) ParseImmediate(s string) int {
	s = strings.Trim(s," ")
	if s[0] == '#' {
		s = s[1:]
	}
	if s[0] == ':' {
		vs := strings.Split(s,":")
		if len(vs) < 3 {
			fmt.Println(s)
			panic("INVALID/UNRECOGNIZED IMMEDIATE")
		}
		switch vs[1] {
		case "lower16":
			di := p.dattab[vs[2]]
			v := p.target.datalabels[di]
			v = v & 0xffff
			return v
		case "upper16":
			v := p.target.datalabels[p.dattab[vs[2]]]
			v = v >> 16
			return v
		default:
			fmt.Println(s)
			panic("INVALID!")
		}
		return -1
	} else {
		n,err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Immediate: %d\n", n)
		return n
	}
}

func (p *Parser) parseValue(s string) []*Val {
	s = strings.Trim(s, " ")
	if p.verbose {
		fmt.Printf("Parsing: '%s'\n",s)
	}
	v := new(Val)
	if s[0] == '#' || isNum(s[0]) {
		v.reg = -1
		v.offs = p.ParseImmediate(s)
	} else if s[0] == '[' {
		ops := strings.Split(s[1:len(s)-1], ",")
		v.reg = p.regValue(ops[0])
		if len(ops) > 1 {
			v.offs = p.ParseImmediate(ops[1])
		}
	} else if s[0] == '{' {
		vals := []*Val{}
		regs := strings.Split(s[1:len(s)-1], ",")
		for _,r := range regs {
			vals = append(vals, p.parseValue(r)...)
		}
		return vals
	} else {
		v.reg = p.regValue(s)
	}
	return []*Val{v}
}

func readToken(b *bytes.Buffer) string {
	out := new(bytes.Buffer)
	in := false
	for b.Len() > 0 {
		b,_ := b.ReadByte()
		if in {
			if b == ']' || b == '}' {
				out.WriteByte(b)
				return out.String()
			}
			out.WriteByte(b)
			continue
		}
		switch b {
		case ' ',',':
			if out.Len() > 0 {
				return out.String()
			}
		case '[','{':
			out.WriteByte(b)
			in = true
		default:
			out.WriteByte(b)
		}
	}
	return out.String()
}

func (p *Parser) parseArgsAlt(s string) ([]*Val, error) {
	var vals []*Val
	var toks []string
	buf := bytes.NewBufferString(s)
	for buf.Len() > 0 {
		toks = append(toks, readToken(buf))
	}
	for i := 0; i < len(toks); i++ {
		if len(toks) > i + 1 {
			switch toks[i+1] {
				case "asr":
					c := p.ParseImmediate(toks[i+2])
					reg := p.parseValue(toks[i])[0]
					reg.shift = -1 * c
					vals = append(vals, reg)
					i += 2
					continue
				case "asl":
					c := p.ParseImmediate(toks[i+2])
					reg := p.parseValue(toks[i])[0]
					reg.shift = c
					vals = append(vals, reg)
					i += 2
					continue
			}
		}
		v := p.parseValue(toks[i])
		vals = append(vals, v...)
	}
	return vals,nil
}

/*
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
					v.offs = p.ParseImmediate(s)
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
				if s[len(s)-1] == '!' {
					fmt.Println("Modify register in place!")
					s = s[:len(s)-1]
				}
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
		case '{','}':
			continue
		default:
			out.WriteByte(b)
		}
	}
	return nil,nil
}
*/

type Parser struct {
	instab map[string]int
	regtab map[string]int
	jmptab map[string]int
	dattab map[string]int
	mode int
	target *Machine
	memloc int
	verbose bool
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
	p.instab["b"]   = Ib
	p.instab["cmp"] = Icmp
	p.instab["ble"] = Ible
	p.instab["stmfd"] = Istmfd
	p.instab["ldmfd"] = Ildmfd
	p.instab["ldrb"]  = Ildrb
	p.instab["strb"]  = Istrb
	p.instab["bls"]   = Ibls
	p.instab["bne"]   = Ibne
	p.instab["movw"] = Imovw
	p.instab["movt"] = Imovt
	p.instab["smull"] = Ismull
	p.instab["rsb"] = Irsb

	//Set up register mappings
	p.regtab = make(map[string]int)
	for i,v := range registers {
		p.regtab[v] = i
	}
	p.regtab["fp"] = 12

	p.dattab = make(map[string]int)

	p.memloc = 4096

	//Map for labels to jumppoints
	p.jmptab = make(map[string]int)
	p.jmptab["putc"] = -1
	return p
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9'
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
	return n == Ibl || n == Ib || n == Ible || n == Ibls ||
	n == Ibne
}

func (p *Parser) ParseInstruction(ss string) *Instruction {
	var err error
	ss = strings.Replace(ss, "\t", " ", -1)
	ss = strings.Trim(ss, " ")
	instr := first(ss)
	switch instr {
	case ".text":
		p.mode = Mtext
		if p.verbose {
			fmt.Println("Reading text data.")
		}
		return nil
	case ".section":
		arg := strings.Split(ss," ")
		if len(arg) > 1 {
			param := strings.Trim(arg[1]," ")
			switch param {
			case ".rodata":
				p.mode = Mdata
				if p.verbose {
					fmt.Println("Reading data.")
				}
			default:
				fmt.Printf("Unknown mode: %s\n", param)
			}
		}
		return nil
	case ".arch",".global":
		if p.verbose {
			fmt.Println(ss)
		}
		return nil
	case ".fpu",".file",".eabi_attribute",".ident",".align":
		return nil
	case ".ascii":
		str := strings.TrimLeft(ss[len(instr)+1:], " ")
		str = strings.Replace(str[1:len(str)-1], "\\000", "", -1)
		str = strings.Replace(str,"\\012", "\n", -1)
		fmt.Printf("String: '%s'\n", str)

		for _,v := range str {
			p.target.setMem(p.memloc, 1, int(v))
			p.memloc++
		}
		return nil
	}

	if instr[0] == '@' {
		return nil
	}
	ins := new(Instruction)
	if instr[len(instr)-1] == ':' {
		if p.mode == Mdata {
			l := instr[:len(instr)-1]
			if p.verbose {
				fmt.Printf("Setting label '%s'\n", l)
			}
			p.dattab[l] = p.target.setDataLabel(p.memloc)
			return nil
		} else if p.mode == Mtext {
			jind := p.jumpMap(instr[:len(instr)-1])
			ins.params = []*Val{&Val{jind,0,0}}
			ins.op = Ilabel
			return ins
		}
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
		ins.params = []*Val{&Val{jv,0,0}}
		return ins
	}

	ins.params,err = p.parseArgsAlt(strings.TrimLeft(ss[len(instr):], " \t"))
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
	p.target = m
	m.srcp = p
	for buf.Scan() {
		i := p.ParseInstruction(buf.Text())
		if i != nil {
			m.addInstruction(i)
		}
	}
	m.Run(p.jumpMap("main"))
	fmt.Println("## Program Execution Finished ##")
	fmt.Println("Final register values:")
	for s,v := range m.regs {
		fmt.Printf("%s = %d\n", registers[s],v)
	}
}

func (p *Parser) getInsName(op int) string {
	for name,i := range p.instab {
		if i == op {
			return name
		}
	}
	return "Unknown"
}
