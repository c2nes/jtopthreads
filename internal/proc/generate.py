import re
import sys

fields = []


def print_option(sect):
    opt = sect[0]
    doc = sect[1:]
    docp = []

    i = 0
    while i < len(doc):
        x = doc[i]
        if len(docp) == 1:
            docp.append("//")
        if len(x) == 1:
            docp.append(f"//   {x}  {doc[i+1]}")
            i += 2
        else:
            docp.append(f"// {x}")
            i += 1

    num, name, spec = opt.split()[:3]
    golang_name = "".join(p.title() for p in name.split("_"))
    golang_type = {
        "%d": "int",
        "%u": "uint",
        "%ld": "int64",
        "%lld": "int64",
        "%lu": "uint64",
        "%llu": "uint64",
        "%s": "string",
        "%c": "rune",
    }[spec]

    # Special formatting rules for process state enumeration
    if name == "state":
        for i in range(2, len(docp)):
            dl = docp[i]
            if re.match(r"// [a-zA-Z]\.  [^ ]", dl):
                _, ch, rest = dl.split(None, 2)
                docp[i] = f"//   {ch[0]}  {rest}"

    # Wrapping
    for i in range(len(docp)):
        dl = docp[i]
        dl = dl[len("// ") :]
        indent = len(dl) - len(dl.lstrip())
        dl = dl[indent:]
        cols = 72 - indent

        # Keep emtpy lines as-is
        if not dl.strip():
            continue

        splits = []
        rest = dl
        while len(rest) > cols:
            sp = cols - 1
            while not rest[sp].isspace():
                sp -= 1
            splits.append(rest[:sp].rstrip())
            rest = rest[sp + 1 :].lstrip()
        if rest:
            splits.append(rest)
        indent = " " * indent
        splits = [f"// {indent}{l}" for l in splits]
        docp[i] = "\n".join(splits)

    # Add indent
    docp = "\n".join(docp)
    docp = "\n".join(f"\t{l}" for l in docp.split("\n"))

    fields.append((name, golang_name, golang_type))

    print(docp)
    print(f"\t{golang_name} {golang_type}")
    print()


# Skip to start of stat documentation
for l in sys.stdin:
    if l.startswith("/proc/[pid]/stat"):
        break

print(
    """// Code generated from proc(5). DO NOT EDIT.

package proc

type ProcStat struct {
"""
)

sect = []
for l in sys.stdin.read().splitlines():
    # Skip blank lines
    if not l.strip():
        continue

    # End of /proc/[pid]/stat section
    if not l.startswith(" "):
        break

    # Remove indentation
    l = l.lstrip()

    if l.startswith("(") and "%" in l:
        # New option
        if sect:
            print_option(sect)
        sect = [l]
        continue
    elif sect:
        sect.append(l)

if sect:
    print_option(sect)

print("}\n")

print(
    """
func (s *ProcStat) parseRest(fields []string) error {
	var err error
	for i, field := range fields {
		switch i {
"""
)

for i, field in enumerate(fields[2:]):
    name, goname, gotype = field

    val = "v"
    if gotype == "rune":
        setter = "setRuneField"
        rawtype = "rune"
    elif gotype.startswith("int"):
        setter = "setIntField"
        rawtype = "int64"
    elif gotype.startswith("uint"):
        setter = "setUintField"
        rawtype = "uint64"
    else:
        raise ValueError(f"unsupported field: {field}")

    if gotype != rawtype:
        val = f"{gotype}({val})"

    print(f"case {i}:")
    print(
        f'err = {setter}("{name}", field, func(v {rawtype}) {{ s.{goname} = {val} }})'
    )

print(
    """
		default:
			break
		}

		if err != nil {
			return err
		}
	}
	return nil
}
"""
)
