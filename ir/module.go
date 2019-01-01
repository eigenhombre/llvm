package ir

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/internal/enc"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/rickypai/natsort"
)

// === [ Modules ] =============================================================

// Module is an LLVM IR module, which consists of top-level declarations and
// definitions.
type Module struct {
	// Type definitions.
	TypeDefs []types.Type
	// Global variable declarations and definitions.
	Globals []*Global
	// Function declarations and definitions.
	Funcs []*Func

	// extra.

	// (optional) Source filename; or empty if not present.
	SourceFilename string
	// (optional) Data layout; or empty if not present.
	DataLayout string
	// (optional) Target triple; or empty if not present.
	TargetTriple string
	// (optional) Module-level inline assembly.
	ModuleAsms []string
	// (optional) Comdat definitions.
	ComdatDefs []*ComdatDef
	// (optional) Aliases.
	Aliases []*Alias
	// (optional) IFuncs.
	IFuncs []*IFunc
	// (optional) Attribute group definitions.
	AttrGroupDefs []*AttrGroupDef
	// (optional) Named metadata definitions.
	NamedMetadataDefs map[string]*metadata.NamedDef
	// (optional) Metadata definitions.
	MetadataDefs []metadata.Definition
	// (optional) Use-list order directives.
	UseListOrders []*UseListOrder
	// (optional) Basic block specific use-list order directives.
	UseListOrderBBs []*UseListOrderBB
}

// NewModule returns a new LLVM IR module.
func NewModule() *Module {
	return &Module{
		NamedMetadataDefs: make(map[string]*metadata.NamedDef),
	}
}

// String returns the string representation of the module in LLVM IR assembly
// syntax.
func (m *Module) String() string {
	buf := &strings.Builder{}
	// Source filename.
	if len(m.SourceFilename) > 0 {
		// 'source_filename' '=' Name=StringLit
		fmt.Fprintf(buf, "source_filename = %s\n", quote(m.SourceFilename))
	}
	// Data layout.
	if len(m.DataLayout) > 0 {
		// 'target' 'datalayout' '=' DataLayout=StringLit
		fmt.Fprintf(buf, "target datalayout = %s\n", quote(m.DataLayout))
	}
	// Target triple.
	if len(m.TargetTriple) > 0 {
		// 'target' 'triple' '=' TargetTriple=StringLit
		fmt.Fprintf(buf, "target triple = %s\n", quote(m.TargetTriple))
	}
	// Module-level inline assembly.
	if len(m.ModuleAsms) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, asm := range m.ModuleAsms {
		// 'module' 'asm' Asm=StringLit
		fmt.Fprintf(buf, "module asm %s\n", quote(asm))
	}
	// Type definitions.
	if len(m.TypeDefs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, t := range m.TypeDefs {
		// Alias=LocalIdent '=' 'type' Typ=OpaqueType
		//
		// Alias=LocalIdent '=' 'type' Typ=Type
		fmt.Fprintf(buf, "%s = type %s\n", t, t.LLString())
	}
	// Comdat definitions.
	if len(m.ComdatDefs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, def := range m.ComdatDefs {
		fmt.Fprintln(buf, def.LLString())
	}
	// Global declarations and definitions.
	if len(m.Globals) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, g := range m.Globals {
		fmt.Fprintln(buf, g.LLString())
	}
	// Aliases.
	if len(m.Aliases) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, alias := range m.Aliases {
		fmt.Fprintln(buf, alias.LLString())
	}
	// IFuncs.
	if len(m.IFuncs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, ifunc := range m.IFuncs {
		fmt.Fprintln(buf, ifunc.LLString())
	}
	// Function declarations and definitions.
	if len(m.Funcs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for i, f := range m.Funcs {
		if i != 0 {
			buf.WriteString("\n")
		}
		fmt.Fprintln(buf, f.LLString())
	}
	// Attribute group definitions.
	if len(m.AttrGroupDefs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, a := range m.AttrGroupDefs {
		fmt.Fprintln(buf, a.LLString())
	}
	// Named metadata definitions; output in natural sorting order.
	var mdNames []string
	for mdName := range m.NamedMetadataDefs {
		mdNames = append(mdNames, mdName)
	}
	natsort.Strings(mdNames)
	if len(m.NamedMetadataDefs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, mdName := range mdNames {
		// Name=MetadataName '=' '!' '{' MDNodes=(MetadataNode separator ',')* '}'
		md := m.NamedMetadataDefs[mdName]
		fmt.Fprintf(buf, "%s = %s\n", md.Ident(), md.LLString())
	}
	// Metadata definitions.
	if len(m.MetadataDefs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, md := range m.MetadataDefs {
		// ID=MetadataID '=' Distinctopt MDNode=MDTuple
		//
		// ID=MetadataID '=' Distinctopt MDNode=SpecializedMDNode
		fmt.Fprintf(buf, "%s = %s\n", md.Ident(), md.LLString())
	}
	// Use-list orders.
	if len(m.UseListOrders) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, u := range m.UseListOrders {
		fmt.Fprintln(buf, u)
	}
	// Basic block specific use-list orders.
	if len(m.UseListOrderBBs) > 0 && buf.Len() > 0 {
		buf.WriteString("\n")
	}
	for _, u := range m.UseListOrderBBs {
		fmt.Fprintln(buf, u)
	}
	return buf.String()
}

// ~~~ [ Comdat Definition ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ComdatDef is a comdat definition top-level entity.
type ComdatDef struct {
	// Comdat name (without '$' prefix).
	Name string
	// Comdat kind.
	Kind enum.SelectionKind
}

// String returns the string representation of the Comdat definition.
func (c *ComdatDef) String() string {
	return fmt.Sprintf("comdat(%s)", enc.Comdat(c.Name))
}

// LLString returns the LLVM syntax representation of the Comdat definition.
func (c *ComdatDef) LLString() string {
	// Name=ComdatName '=' 'comdat' Kind=SelectionKind
	return fmt.Sprintf("%s = comdat %s", enc.Comdat(c.Name), c.Kind)
}

// ~~~ [ Attribute Group Definition ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AttrGroupDef is an attribute group definition.
type AttrGroupDef struct {
	// Attribute group ID (without '#' prefix).
	ID int64
	// Function attributes.
	FuncAttrs []FuncAttribute
}

// String returns the string representation of the attribute group definition.
func (a *AttrGroupDef) String() string {
	return enc.AttrGroupID(a.ID)
}

// LLString returns the LLVM syntax representation of the attribute group
// definition.
func (a *AttrGroupDef) LLString() string {
	// 'attributes' ID=AttrGroupID '=' '{' Attrs=FuncAttribute* '}'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "attributes %s = { ", enc.AttrGroupID(a.ID))
	for i, attr := range a.FuncAttrs {
		if i != 0 {
			buf.WriteString(" ")
		}
		switch attr := attr.(type) {
		case Align:
			// Note, alignment is printed as `align = 8` in attribute groups.
			fmt.Fprintf(buf, "align = %d", uint64(attr))
		case AlignStack:
			// Note, stack alignment is printed as `alignstack = 8` in attribute
			// groups.
			fmt.Fprintf(buf, "alignstack = %d", uint64(attr))
		default:
			buf.WriteString(attr.String())
		}
	}
	buf.WriteString(" }")
	return buf.String()
}

// ~~~ [ Use-list Order Directives ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UseListOrder is a use-list order directive.
type UseListOrder struct {
	// Value.
	Value value.Value
	// Use-list order.
	Indices []uint64
}

// String returns the string representation of the use-list order directive
// definition.
func (u *UseListOrder) String() string {
	//  'uselistorder' TypeValue ',' '{' Indices=(UintLit separator ',')+ '}'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "uselistorder %s, { ", u.Value)
	for i, index := range u.Indices {
		if i != 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(buf, "%d", index)
	}
	buf.WriteString(" }")
	return buf.String()
}

// UseListOrderBB is a basic block specific use-list order directive.
type UseListOrderBB struct {
	// Function.
	Func *Func
	// Basic block.
	Block *Block
	// Use-list order.
	Indices []uint64
}

// String returns the string representation of the basic block specific use-
// list order directive definition.
func (u *UseListOrderBB) String() string {
	//  'uselistorder_bb' Func=GlobalIdent ',' Block=LocalIdent ',' '{'
	//  Indices=(UintLit separator ',')+ '}'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "uselistorder_bb %s, %s, { ", u.Func.Ident(), u.Block.Ident())
	for i, index := range u.Indices {
		if i != 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(buf, "%d", index)
	}
	buf.WriteString(" }")
	return buf.String()
}
