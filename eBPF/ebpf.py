#!/usr/bin/python3
# -*- coding: utf-8 -*-

# 加载 BCC 库
from bcc import BPF

# 加载 eBPF 内核态程序
b = BPF(src_file="ebpf.c")

# 将 eBPF 程序挂载到 kprobe
b.attach_kprobe(event="do_sys_openat2", fn_name="ebpf")

# 读取并且打印 eBPF 内核态程序输出的数据
b.trace_print()
