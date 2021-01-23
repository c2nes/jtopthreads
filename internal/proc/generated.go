// Code generated from proc(5). DO NOT EDIT.

package proc

type ProcStat struct {

	// The process ID.
	Pid int

	// The filename of the executable, in parentheses. Strings longer than
	// TASK_COMM_LEN (16) characters (including the terminating null byte) are
	// silently truncated. This is visible whether or not the executable is
	// swapped out.
	Comm string

	// One of the following characters, indicating process state:
	//
	//   R  Running
	//   S  Sleeping in an interruptible wait
	//   T  Waiting in uninterruptible disk sleep
	//   U  Zombie
	//   V  Stopped (on a signal) or (before Linux 2.6.33) trace stopped
	//   t  Tracing stop (Linux 2.6.33 onward)
	//   W  Paging (only before Linux 2.6.0)
	//   X  Dead (from Linux 2.6.0 onward)
	//   x  Dead (Linux 2.6.33 to 3.13 only)
	//   K  Wakekill (Linux 2.6.33 to 3.13 only)
	//   L  Waking (Linux 2.6.33 to 3.13 only)
	//   M  Parked (Linux 3.9 to 3.13 only)
	State rune

	// The PID of the parent of this process.
	Ppid int

	// The process group ID of the process.
	Pgrp int

	// The session ID of the process.
	Session int

	// The controlling terminal of the process. (The minor device number is
	// contained in the combination of bits 31 to 20 and 7 to 0; the major
	// device number is in bits 15 to 8.)
	TtyNr int

	// The ID of the foreground process group of the controlling terminal of
	// the process.
	Tpgid int

	// The kernel flags word of the process. For bit meanings, see the PF_*
	// defines in the Linux kernel source file include/linux/sched.h. Details
	// depend on the kernel version.
	//
	// The format for this field was %lu before Linux 2.6.
	Flags uint

	// The number of minor faults the process has made which have not required
	// loading a memory page from disk.
	Minflt uint64

	// The number of minor faults that the process's waited-for children have
	// made.
	Cminflt uint64

	// The number of major faults the process has made which have required
	// loading a memory page from disk.
	Majflt uint64

	// The number of major faults that the process's waited-for children have
	// made.
	Cmajflt uint64

	// Amount of time that this process has been scheduled in user mode,
	// measured in clock ticks (divide by sysconf(_SC_CLK_TCK)). This includes
	// guest time, guest_time (time spent running a virtual CPU, see below),
	// so that applications that are not aware of the guest time field do not
	// lose that time from their calculations.
	Utime uint64

	// Amount of time that this process has been scheduled in kernel mode,
	// measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	Stime uint64

	// Amount of time that this process's waited-for children have been
	// scheduled in user mode, measured in clock ticks (divide by
	// sysconf(_SC_CLK_TCK)). (See also times(2).) This includes guest time,
	// cguest_time (time spent running a virtual CPU, see below).
	Cutime int64

	// Amount of time that this process's waited-for children have been
	// scheduled in kernel mode, measured in clock ticks (divide by
	// sysconf(_SC_CLK_TCK)).
	Cstime int64

	// (Explanation for Linux 2.6) For processes running a real-time
	// scheduling policy (policy below; see sched_setscheduler(2)), this is
	// the negated scheduling priority, minus one; that is, a number in the
	// range -2 to -100, corresponding to real-time priorities 1 to 99. For
	// processes running under a non-real-time scheduling policy, this is the
	// raw nice value (setpriority(2)) as represented in the kernel. The
	// kernel stores nice values as numbers in the range 0 (high) to 39 (low),
	// corresponding to the user-visible nice range of -20 to 19.
	//
	// Before Linux 2.6, this was a scaled value based on the scheduler
	// weighting given to this process.
	Priority int64

	// The nice value (see setpriority(2)), a value in the range 19 (low
	// priority) to -20 (high priority).
	Nice int64

	// Number of threads in this process (since Linux 2.6). Before kernel 2.6,
	// this field was hard coded to 0 as a placeholder for an earlier removed
	// field.
	NumThreads int64

	// The time in jiffies before the next SIGALRM is sent to the process due
	// to an interval timer. Since kernel 2.6.17, this field is no longer
	// maintained, and is hard coded as 0.
	Itrealvalue int64

	// The time the process started after system boot. In kernels before Linux
	// 2.6, this value was expressed in jiffies. Since Linux 2.6, the value is
	// expressed in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	//
	// The format for this field was %lu before Linux 2.6.
	Starttime uint64

	// Virtual memory size in bytes.
	Vsize uint64

	// Resident Set Size: number of pages the process has in real memory. This
	// is just the pages which count toward text, data, or stack space. This
	// does not include pages which have not been demand-loaded in, or which
	// are swapped out. This value is inaccurate; see /proc/[pid]/statm below.
	Rss int64

	// Current soft limit in bytes on the rss of the process; see the
	// description of RLIMIT_RSS in getrlimit(2).
	Rsslim uint64

	// The address above which program text can run.
	Startcode uint64

	// The address below which program text can run.
	Endcode uint64

	// The address of the start (i.e., bottom) of the stack.
	Startstack uint64

	// The current value of ESP (stack pointer), as found in the kernel stack
	// page for the process.
	Kstkesp uint64

	// The current EIP (instruction pointer).
	Kstkeip uint64

	// The bitmap of pending signals, displayed as a decimal number. Obsolete,
	// because it does not provide information on real-time signals; use
	// /proc/[pid]/status instead.
	Signal uint64

	// The bitmap of blocked signals, displayed as a decimal number. Obsolete,
	// because it does not provide information on real-time signals; use
	// /proc/[pid]/status instead.
	Blocked uint64

	// The bitmap of ignored signals, displayed as a decimal number. Obsolete,
	// because it does not provide information on real-time signals; use
	// /proc/[pid]/status instead.
	Sigignore uint64

	// The bitmap of caught signals, displayed as a decimal number. Obsolete,
	// because it does not provide information on real-time signals; use
	// /proc/[pid]/status instead.
	Sigcatch uint64

	// This is the "channel" in which the process is waiting. It is the
	// address of a location in the kernel where the process is sleeping. The
	// corresponding symbolic name can be found in /proc/[pid]/wchan.
	Wchan uint64

	// Number of pages swapped (not maintained).
	Nswap uint64

	// Cumulative nswap for child processes (not maintained).
	Cnswap uint64

	// Signal to be sent to parent when we die.
	ExitSignal int

	// CPU number last executed on.
	Processor int

	// Real-time scheduling priority, a number in the range 1 to 99 for
	// processes scheduled under a real-time policy, or 0, for non-real-time
	// processes (see sched_setscheduler(2)).
	RtPriority uint

	// Scheduling policy (see sched_setscheduler(2)). Decode using the SCHED_*
	// constants in linux/sched.h.
	//
	// The format for this field was %lu before Linux 2.6.22.
	Policy uint

	// Aggregated block I/O delays, measured in clock ticks (centiseconds).
	DelayacctBlkioTicks uint64

	// Guest time of the process (time spent running a virtual CPU for a guest
	// operating system), measured in clock ticks (divide by
	// sysconf(_SC_CLK_TCK)).
	GuestTime uint64

	// Guest time of the process's children, measured in clock ticks (divide
	// by sysconf(_SC_CLK_TCK)).
	CguestTime int64

	// Address above which program initialized and uninitialized (BSS) data
	// are placed.
	StartData uint64

	// Address below which program initialized and uninitialized (BSS) data
	// are placed.
	EndData uint64

	// Address above which program heap can be expanded with brk(2).
	StartBrk uint64

	// Address above which program command-line arguments (argv) are placed.
	ArgStart uint64

	// Address below program command-line arguments (argv) are placed.
	ArgEnd uint64

	// Address above which program environment is placed.
	EnvStart uint64

	// Address below which program environment is placed.
	EnvEnd uint64

	// The thread's exit status in the form reported by waitpid(2).
	ExitCode int
}

func (s *ProcStat) parseRest(fields []string) error {
	var err error
	for i, field := range fields {
		switch i {

		case 0:
			err = setRuneField("state", field, func(v rune) { s.State = v })
		case 1:
			err = setIntField("ppid", field, func(v int64) { s.Ppid = int(v) })
		case 2:
			err = setIntField("pgrp", field, func(v int64) { s.Pgrp = int(v) })
		case 3:
			err = setIntField("session", field, func(v int64) { s.Session = int(v) })
		case 4:
			err = setIntField("tty_nr", field, func(v int64) { s.TtyNr = int(v) })
		case 5:
			err = setIntField("tpgid", field, func(v int64) { s.Tpgid = int(v) })
		case 6:
			err = setUintField("flags", field, func(v uint64) { s.Flags = uint(v) })
		case 7:
			err = setUintField("minflt", field, func(v uint64) { s.Minflt = v })
		case 8:
			err = setUintField("cminflt", field, func(v uint64) { s.Cminflt = v })
		case 9:
			err = setUintField("majflt", field, func(v uint64) { s.Majflt = v })
		case 10:
			err = setUintField("cmajflt", field, func(v uint64) { s.Cmajflt = v })
		case 11:
			err = setUintField("utime", field, func(v uint64) { s.Utime = v })
		case 12:
			err = setUintField("stime", field, func(v uint64) { s.Stime = v })
		case 13:
			err = setIntField("cutime", field, func(v int64) { s.Cutime = v })
		case 14:
			err = setIntField("cstime", field, func(v int64) { s.Cstime = v })
		case 15:
			err = setIntField("priority", field, func(v int64) { s.Priority = v })
		case 16:
			err = setIntField("nice", field, func(v int64) { s.Nice = v })
		case 17:
			err = setIntField("num_threads", field, func(v int64) { s.NumThreads = v })
		case 18:
			err = setIntField("itrealvalue", field, func(v int64) { s.Itrealvalue = v })
		case 19:
			err = setUintField("starttime", field, func(v uint64) { s.Starttime = v })
		case 20:
			err = setUintField("vsize", field, func(v uint64) { s.Vsize = v })
		case 21:
			err = setIntField("rss", field, func(v int64) { s.Rss = v })
		case 22:
			err = setUintField("rsslim", field, func(v uint64) { s.Rsslim = v })
		case 23:
			err = setUintField("startcode", field, func(v uint64) { s.Startcode = v })
		case 24:
			err = setUintField("endcode", field, func(v uint64) { s.Endcode = v })
		case 25:
			err = setUintField("startstack", field, func(v uint64) { s.Startstack = v })
		case 26:
			err = setUintField("kstkesp", field, func(v uint64) { s.Kstkesp = v })
		case 27:
			err = setUintField("kstkeip", field, func(v uint64) { s.Kstkeip = v })
		case 28:
			err = setUintField("signal", field, func(v uint64) { s.Signal = v })
		case 29:
			err = setUintField("blocked", field, func(v uint64) { s.Blocked = v })
		case 30:
			err = setUintField("sigignore", field, func(v uint64) { s.Sigignore = v })
		case 31:
			err = setUintField("sigcatch", field, func(v uint64) { s.Sigcatch = v })
		case 32:
			err = setUintField("wchan", field, func(v uint64) { s.Wchan = v })
		case 33:
			err = setUintField("nswap", field, func(v uint64) { s.Nswap = v })
		case 34:
			err = setUintField("cnswap", field, func(v uint64) { s.Cnswap = v })
		case 35:
			err = setIntField("exit_signal", field, func(v int64) { s.ExitSignal = int(v) })
		case 36:
			err = setIntField("processor", field, func(v int64) { s.Processor = int(v) })
		case 37:
			err = setUintField("rt_priority", field, func(v uint64) { s.RtPriority = uint(v) })
		case 38:
			err = setUintField("policy", field, func(v uint64) { s.Policy = uint(v) })
		case 39:
			err = setUintField("delayacct_blkio_ticks", field, func(v uint64) { s.DelayacctBlkioTicks = v })
		case 40:
			err = setUintField("guest_time", field, func(v uint64) { s.GuestTime = v })
		case 41:
			err = setIntField("cguest_time", field, func(v int64) { s.CguestTime = v })
		case 42:
			err = setUintField("start_data", field, func(v uint64) { s.StartData = v })
		case 43:
			err = setUintField("end_data", field, func(v uint64) { s.EndData = v })
		case 44:
			err = setUintField("start_brk", field, func(v uint64) { s.StartBrk = v })
		case 45:
			err = setUintField("arg_start", field, func(v uint64) { s.ArgStart = v })
		case 46:
			err = setUintField("arg_end", field, func(v uint64) { s.ArgEnd = v })
		case 47:
			err = setUintField("env_start", field, func(v uint64) { s.EnvStart = v })
		case 48:
			err = setUintField("env_end", field, func(v uint64) { s.EnvEnd = v })
		case 49:
			err = setIntField("exit_code", field, func(v int64) { s.ExitCode = int(v) })

		default:
			break
		}

		if err != nil {
			return err
		}
	}
	return nil
}
